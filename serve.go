package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleDownloads(r *gin.Engine, saveDir string, uploaded *FileMap) {
	// todo: maybe use buffer instead of reading the whole file into memory
	r.GET("/download/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		c.FileAttachment(saveDir+hash, file.FileName)
	})

	r.GET("/cdn/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(http.StatusNotFound, "File not found")
			return
		}

		c.File(saveDir + hash)
	})

	r.GET("/total_files", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"files": uploaded.Files,
		})
	})
}
