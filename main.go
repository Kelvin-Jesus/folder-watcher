package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/kelvin-jesus/folder-watcher/utils"

	"github.com/fsnotify/fsnotify"
)

func main() {
	ini := utils.IniFile{}
	iniFile := ini.Load()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Started watching folder: ", iniFile.FolderToWatch)
	defer watcher.Close()

	err = watcher.Add(iniFile.FolderToWatch)
	if err != nil {
		log.Fatal(err)
	}

	startEventLoop(watcher, sendReceivedFileOverHTTP)
}

func sendReceivedFileOverHTTP(filePath string) {
	log.Println("Enviando requisição do arquivo: ", filePath)
	receivedFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	fileName, _ := receivedFile.Stat()
	fileStruct := utils.ReceivedFile{
		Pointer: receivedFile,
		Name:    fileName.Name(),
		Path:    filePath,
	}

	defer receivedFile.Close()

	var requestBody bytes.Buffer

	multipartWriter := multipart.NewWriter(&requestBody)

	err = multipartWriter.WriteField("tempo", iniFile.TimeInMinutes)
	if err != nil {
		log.Println("Error writing form field:", err)
		return
	}

	fileWriter, err := multipartWriter.CreateFormFile("file", filePath)
	if err != nil {
		log.Println("Error creating form file:", err)
		return
	}

	_, err = io.Copy(fileWriter, receivedFile)
	if err != nil {
		log.Println("Error copying file content:", err)
		return
	}

	multipartWriter.Close()

	req, err := http.NewRequest("POST", URL, &requestBody)
	if err != nil {
		log.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+authToken)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("Error sending HTTP request:", err)
		return
	}

	response.Body.Close()

	log.Println("Status Code:", response.StatusCode)
	if response.StatusCode != http.StatusOK {
		moveFile.MoveToFolder(onErrorFolder, &fileStruct)
		log.Printf("HTTP request failed with status code %d\n", response.StatusCode)
		return
	}

	moveFile.ToFolder(onSuccessFolder, &fileStruct)
}

func startEventLoop(
	fsWatcher *fsnotify.Watcher,
	execOnCreatedFileFN func(filePath string),
) {
	for {
		select {
		case event, ok := <-fsWatcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Create == fsnotify.Create {
				if !strings.Contains(event.String(), ".part") {
					log.Println(event.Name)
					go execOnCreatedFileFN(event.Name)
				}
			}
		case err, ok := <-fsWatcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}
