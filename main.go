package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

type FileMapKey struct {
	FileName   string `json:"filename"`
	UploadDate int64  `json:"upload_date"`
}

type FileMap struct {
	Files map[string]FileMapKey `json:"files"`
}

type ApiKeys struct {
	Keys map[string]bool `json:"keys"`
}

func main() {
	wd, _ := os.Getwd()

	saveDir := wd + "/uploads/"
	dataDir := wd + "/data/"

	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		os.Mkdir(saveDir, os.ModePerm)
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, os.ModePerm)
	}

	fileMapData := readFile(dataDir+"filemap.json", []byte("{\"files\": {}}"))
	apiKeysData := readFile(dataDir+"apikeys.json", []byte("{\"keys\": {}}"))

	apiKeys := ApiKeys{}
	uploaded := FileMap{}

	json.Unmarshal(fileMapData, &uploaded)
	json.Unmarshal(apiKeysData, &apiKeys)

	log.Println(uploaded.Files)

	r := gin.Default()
	handleUploads(r, saveDir, dataDir, &uploaded, &apiKeys)
	r.Run()
}
