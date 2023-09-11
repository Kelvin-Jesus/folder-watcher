package utils

import (
	"fmt"
	"os"
)

func MoveToFolder(destinatiorFolder string, file *ReceivedFile) {
	folderSeparator := string(os.PathSeparator)

	err := os.Rename(file.Path, destinatiorFolder+folderSeparator+file.Name)
	if err != nil {
		fmt.Println("Error moving file:", err)
		return
	}

	file.Close()
}
