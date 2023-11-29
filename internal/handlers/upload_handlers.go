package handlers

import (
	"log"
	"os"
)

const (
	uploadPath = "./uploads/"
	chunkSize  = 5 * 1024 * 1024 // 5 MB
)

func init() {
	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.Mkdir(uploadPath, os.ModePerm)
	}
}

func removeFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Printf("Ошибка при удалении файла: %v", err)
	}
}
