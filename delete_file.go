package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"os"
)

func handleDeletions(r *gin.Engine, saveDir string, dataDir string, uploaded *FileMap, apiKeys *ApiKeys) {
	r.DELETE("/delete/:hash", func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !apiKeys.Keys[authHeader] {
			c.String(401, "Unauthorized")
			return
		}

		hash := c.Param("hash")
		file := uploaded.Files[hash]
		if file == (FileMapKey{}) {
			c.String(404, "File not found")
			return
		}
		c.String(200, "Deleted %s (%s)\n", file.FileName, hash)
		os.Remove(saveDir + hash)
		delete(uploaded.Files, hash)
		dat, _ := json.Marshal(uploaded)
		os.WriteFile(dataDir+"filemap.json", dat, os.ModePerm)
		// todo: timer is not removed but it should be in the future
	})
}
