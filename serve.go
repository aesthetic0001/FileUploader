package main

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func handleDownloads(r *gin.Engine, saveDir string, uploaded *FileMap) {
	r.GET("/download/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		uploadedFile, _ := os.Open(saveDir + hash)
		r := bufio.NewReader(uploadedFile)
		content, _ := r.ReadBytes('\n')
		uploadedFile.Close()

		c.Data(http.StatusOK, "application/octet-stream", content)
	})

	r.GET("/cdn/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		uploadedFile, _ := os.Open(saveDir + hash)
		r := bufio.NewReader(uploadedFile)
		content, _ := r.ReadBytes('\n')
		uploadedFile.Close()
		contentType := http.DetectContentType(content)

		c.Data(http.StatusOK, contentType, content)
	})

	r.GET("/total_files", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"files": uploaded.Files,
		})
	})
}
