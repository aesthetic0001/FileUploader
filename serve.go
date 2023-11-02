package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type cachedFile struct {
	ContentType string
	Content     []byte
}

func handleDownloads(r *gin.Engine, saveDir string, uploaded *FileMap) {
	fileCache := make(map[string]cachedFile)

	// todo: maybe use buffer instead of reading the whole file into memory
	r.GET("/download/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		if _, exists := fileCache[hash]; !exists {
			content, _ := os.ReadFile(saveDir + hash)
			contentType := http.DetectContentType(content)
			fileCache[hash] = cachedFile{
				ContentType: contentType,
				Content:     content,
			}
		}

		c.Data(http.StatusOK, "application/octet-stream", fileCache[hash].Content)
	})

	r.GET("/cdn/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		if _, exists := fileCache[hash]; !exists {
			content, _ := os.ReadFile(saveDir + hash)
			contentType := http.DetectContentType(content)
			fileCache[hash] = cachedFile{
				ContentType: contentType,
				Content:     content,
			}
		}

		c.Data(http.StatusOK, fileCache[hash].ContentType, fileCache[hash].Content)
	})

	r.GET("/total_files", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"files": uploaded.Files,
		})
	})
}
