package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"time"
)

func handleUploads(r *gin.Engine, saveDir string, dataDir string, uploaded *FileMap, apiKeys *ApiKeys) {
	r.MaxMultipartMemory = 8 << 31 // 2GiB
	r.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		h := sha256.New()
		fileContents, _ := file.Open()
		io.Copy(h, fileContents)
		hash := fmt.Sprintf("%x", h.Sum(nil))
		if uploaded.Files[hash] != (FileMapKey{}) {
			c.JSON(http.StatusOK, gin.H{
				"hash":   hash,
				"status": "ok",
			})
			return
		}
		fmt.Printf("Storing %s (%s)\n", file.Filename, hash)
		if err := c.SaveUploadedFile(file, saveDir+hash); err != nil {
			return
		}
		uploaded.Files[hash] = FileMapKey{
			FileName:   file.Filename,
			UploadDate: time.Now().UnixMilli(),
		}
		handleFileExpiry(saveDir, dataDir, hash, uploaded)
		dat, _ := json.Marshal(uploaded)
		os.WriteFile(dataDir+"filemap.json", dat, os.ModePerm)
		c.JSON(http.StatusOK, gin.H{
			"hash": fmt.Sprintf("%x", h.Sum(nil)),
		})
	})
}
