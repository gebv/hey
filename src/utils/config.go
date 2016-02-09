package utils

import (
	"encoding/json"
	"models"
	"os"
	"path/filepath"
)

var Cfg *models.Config = &models.Config{}

var BuildDate string = ""
var Version string = ""
var GitHash string = ""
var IsTesting bool

const (
	MODE_PROD = "prod"
	MODE_DEV  = "dev"
)

func FindConfigFile(fileName string) string {
	if len(fileName) == 0 {
		panic("empty file name")
	}

	if _, err := os.Stat("./config/" + fileName); err == nil {
		fileName, _ = filepath.Abs("./config/" + fileName)
	} else if _, err := os.Stat("../config/" + fileName); err == nil {
		fileName, _ = filepath.Abs("../config/" + fileName)
	} else if _, err := os.Stat(fileName); err == nil {
		fileName, _ = filepath.Abs(fileName)
	} else {
		panic("not found " + fileName)
	}

	return fileName
}

func FindDir(dir string) string {

	if len(dir) == 0 {
		panic("empty dir name")
	}

	fileName := "."
	if _, err := os.Stat("./" + dir + "/"); err == nil {
		fileName, _ = filepath.Abs("./" + dir + "/")
	} else if _, err := os.Stat("../" + dir + "/"); err == nil {
		fileName, _ = filepath.Abs("../" + dir + "/")
	} else {
		panic("not found " + dir)
	}

	return fileName + "/"
}

func LoadConfig(fileName string) {
	if (len(BuildDate) == 0 || len(GitHash) == 0 || len(Version) == 0) && !IsTesting {
		panic("empty BuildDate OR GitHash OR Version")
	}

	fileName = FindConfigFile(fileName)

	file, err := os.Open(fileName)

	if err != nil {
		panic("error opening config file=" + fileName + ", err=" + err.Error())
	}

	decoder := json.NewDecoder(file)
	config := models.Config{}
	err = decoder.Decode(&config)
	if err != nil {
		panic("error decoding config file=" + fileName + ", err=" + err.Error())
	}

	Cfg = &config
}
