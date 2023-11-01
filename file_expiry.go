package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func handleFileExpiry(saveDir string, dataDir string, fileHash string, uploaded *FileMap) {
	file := uploaded.Files[fileHash]
	remainingTime := (time.Now().UnixMilli() - file.UploadDate) + int64(time.Hour*24*180)

	time.AfterFunc(time.Duration(remainingTime), func() {
		os.Remove(saveDir + fileHash)
		fmt.Printf("Deleting %s (%s)\n", file.FileName, fileHash)
		delete(uploaded.Files, fileHash)
		dat, _ := json.Marshal(uploaded)
		os.WriteFile(dataDir+"filemap.json", dat, os.ModePerm)
	})

}
