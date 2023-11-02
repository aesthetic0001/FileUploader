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
	r.POST("/upload", func(c *gin.Context) {
		mpReader, _ := c.Request.MultipartReader()
		for {
			part, err := mpReader.NextPart()
			if err == io.EOF {
				break
			} else if err != nil {
				panic(err)
			}
			if part.FileName() == "" {
				continue
			}
			fileName := part.FileName()
			tempFileName := fmt.Sprintf("%stmp%s-%d", saveDir, fileName, time.Now().UnixMilli())
			outFile, _ := os.Create(tempFileName)
			buf := make([]byte, 1024)
			hash := sha256.New()
			for {
				n, err := part.Read(buf)
				if err == io.EOF {
					break
				}
				if _, err := hash.Write(buf[:n]); err != nil {
					panic(err)
				}
				if _, err := outFile.Write(buf[:n]); err != nil {
					panic(err)
				}
			}
			outFile.Close()
			hashString := fmt.Sprintf("%x", hash.Sum(nil))
			if uploaded.Files[hashString] != (FileMapKey{}) {
				fmt.Println("File already exists! Deleting temp file...")
				os.Remove(tempFileName)
				c.JSON(http.StatusOK, gin.H{
					"hash": hashString,
				})
				return
			}
			fmt.Printf("File uploaded: %s (%s)\n", fileName, hashString)
			os.Rename(tempFileName, saveDir+hashString)
			uploaded.Files[hashString] = FileMapKey{
				FileName:   fileName,
				UploadDate: time.Now().UnixMilli(),
			}
			marshalled, _ := json.Marshal(uploaded)
			os.WriteFile(dataDir+"filemap.json", marshalled, os.ModePerm)
			c.JSON(http.StatusOK, gin.H{
				"hash": hashString,
			})
			return
		}
	})
}
