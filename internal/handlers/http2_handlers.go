package handlers

import (
	"fmt"
	"golang.org/x/net/context"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type FileInfo struct {
	totalBytesRead int64
	totalSize      int64
	percent        float64
}

var fileInfo = &FileInfo{}

func UploadProgress(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Stream not support", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	tick := time.Tick(time.Second / 10)

	for {
		select {
		case <-ctx.Done():
			fileInfo = &FileInfo{}
			return
		case <-tick:
			fmt.Fprintf(w, "data: %f\n\n", fileInfo.percent)
			flusher.Flush()
		}
	}
}

func UploadHttp2(w http.ResponseWriter, r *http.Request) {
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	totalSize, err := strconv.ParseFloat(r.FormValue("totalSize"), 64)
	if err != nil {
		http.Error(w, "Error get size file", http.StatusBadRequest)
		return
	}

	fileName := handler.Filename
	savePath := filepath.Join(uploadPath, fileName)

	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "Error create file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error read chunk file", http.StatusInternalServerError)
		return
	}

	_, err = out.Write(fileBytes)
	if err != nil {
		http.Error(w, "Error copy file", http.StatusInternalServerError)
		return
	}

	fileInfo.totalBytesRead += handler.Size
	fileInfo.percent = math.Round(float64(fileInfo.totalBytesRead) / totalSize * 100)
}
