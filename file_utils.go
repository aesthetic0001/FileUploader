package main

import "os"

func readFile(path string, placeholder []byte) []byte {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, _ := os.Create(path)
		file.Write(placeholder)
		file.Close()
	}

	data, _ := os.ReadFile(path)

	return data
}
