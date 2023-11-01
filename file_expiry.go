package main

import (
	"fmt"
	"os"
	"time"
)

func handleFileExpiry(directory string, fileHash string, uploaded *FileMap) {
	file := uploaded.Files[fileHash]
	remainingTime := (time.Now().UnixMilli() - file.UploadDate) + 24*60*60*1000

	time.AfterFunc(time.Duration(remainingTime), func() {
		fmt.Printf("Deleting %s (%s)\n", file.FileName, fileHash)
		os.Remove(directory + file.FileName)
		delete(uploaded.Files, fileHash)
	})
}
