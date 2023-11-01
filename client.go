package main

import (
	"log"
	"net/http"
)

func upload() {
	// multipart upload to localhost:8080/upload
	// Create a new file upload request
	_, err := http.NewRequest("POST", "http://localhost:8080/upload", nil)
	if err != nil {
		log.Fatal(err)
	}
}
