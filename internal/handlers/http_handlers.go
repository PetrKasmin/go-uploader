package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Upload(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileName := handler.Filename
	savePath := filepath.Join(uploadPath, fileName)

	out, err := os.OpenFile(savePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		http.Error(w, "Error open file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Error copy file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "File %s uploaded!", fileName)
}
