package main

import (
	"github.com/bmizerany/pat"
	"go-upload/internal/handlers"
	"net/http"
)

func routes() http.Handler {
	mux := pat.New()

	mux.Get("/", http.HandlerFunc(handlers.Main))
	mux.Post("/upload", http.HandlerFunc(handlers.Upload))
	mux.Post("/upload-http2", http.HandlerFunc(handlers.UploadHttp2))
	mux.Get("/upload-progress", http.HandlerFunc(handlers.UploadProgress))
	mux.Get("/upload-websocket", http.HandlerFunc(handlers.UploadWebsocket))

	return mux
}
