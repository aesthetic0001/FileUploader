package main

import "github.com/gin-gonic/gin"

func handleDownloads(r *gin.Engine, saveDir string, uploaded *FileMap) {
	r.GET("/download/:hash", func(c *gin.Context) {
		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(404, "File not found")
			return
		}
		c.FileAttachment(saveDir+hash, file.FileName)
	})
}
