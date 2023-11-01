package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
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

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(CORSMiddleware())
	handleUploads(r, saveDir, dataDir, &uploaded, &apiKeys)
	handleDownloads(r, saveDir, &uploaded)
	handleDeletions(r, saveDir, dataDir, &uploaded, &apiKeys)
	for fileHash := range uploaded.Files {
		handleFileExpiry(saveDir, dataDir, fileHash, &uploaded)
	}
	r.Run(":8080")
}
