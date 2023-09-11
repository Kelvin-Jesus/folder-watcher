package utils

import (
	"log"

	"github.com/go-ini/ini"
)

type IniFile struct {
	FolderToWatch   string
	OnSuccessFolder string
	OnErrorFolder   string
	TimeInMinutes   string
	AuthToken       string
	URLToSendFiles  string
}

func (iniFile *IniFile) Load() *IniFile {
	configurationVariables, err := ini.Load("send.ini")
	if err != nil {
		log.Fatal("Failed to load .ini file", err)
	}

	return &IniFile{
		FolderToWatch:   configurationVariables.Section("VARS").Key("FOLDER_TO_WATCH").String(),
		OnSuccessFolder: configurationVariables.Section("VARS").Key("ON_SUCCESS_FOLDER").String(),
		OnErrorFolder:   configurationVariables.Section("VARS").Key("ON_ERROR_FOLDER").String(),
		TimeInMinutes:   configurationVariables.Section("VARS").Key("TIME_IN_MINUTES").String(),
		AuthToken:       configurationVariables.Section("VARS").Key("AUTH_TOKEN").String(),
		URLToSendFiles:  configurationVariables.Section("VARS").Key("URL_TO_SEND").String(),
	}
}
