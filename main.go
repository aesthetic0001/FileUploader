package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
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
	r.MaxMultipartMemory = 8 << 20 // 8 MiB
	r.POST("/upload", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !apiKeys.Keys[authHeader] {
			c.String(http.StatusUnauthorized, "Unauthorized")
			return
		}

		log.Println(authHeader)

		form, _ := c.MultipartForm()
		files := form.File["upload[]"]
		var hashes = make([]string, len(files))

		for i, file := range files {
			h := sha256.New()
			fileContents, _ := file.Open()
			io.Copy(h, fileContents)
			hashes[i] = fmt.Sprintf("%x", h.Sum(nil))
			if uploaded.Files[hashes[i]] != (FileMapKey{}) {
				log.Println("File already uploaded")
				continue
			}
			if err := c.SaveUploadedFile(file, saveDir+hashes[i]); err != nil {
				log.Println(err)
				return
			}
			uploaded.Files[hashes[i]] = FileMapKey{
				FileName:   file.Filename,
				UploadDate: time.Now().UnixMilli(),
			}
			dat, _ := json.Marshal(uploaded)
			os.WriteFile(dataDir+"filemap.json", dat, os.ModePerm)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded! | Hashes: %s", len(files), hashes))
	})
	r.Run()
}
