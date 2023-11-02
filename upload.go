package main

import (
	"bufio"
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
			hash := sha256.New()
			writer := bufio.NewWriterSize(outFile, 8*1024*1024)
			// copy part to writer
			io.Copy(writer, io.TeeReader(part, hash))
			writer.Flush()
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
