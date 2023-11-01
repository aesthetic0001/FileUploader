package main

import (
	"github.com/gin-gonic/gin"
	"os"
)

func handleDownloads(r *gin.Engine, saveDir string, uploaded *FileMap) {
	fileCache := make(map[string][]byte)

	r.GET("/download/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(404, "File not found")
			return
		}
		c.FileAttachment(saveDir+hash, file.FileName)
	})

	r.GET("/cdn/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(404, "File not found")
			return
		}
		if fileCache[hash] == nil {
			fileCache[hash], _ = os.ReadFile(saveDir + hash)
		}
		// serve file from memory with filename
		c.Data(200, "application/octet-stream", fileCache[hash])
	})
}
