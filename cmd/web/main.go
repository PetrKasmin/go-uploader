package main

import (
	"go-upload/internal/handlers"
	"golang.org/x/net/http2"
	"log"
	"net/http"
)

func main() {
	mux := routes()

	go handlers.ListenToWsChannel()

	server := &http.Server{
		Addr:    ":3000",
		Handler: mux,
	}

	err := http2.ConfigureServer(server, &http2.Server{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Server listening on :3000...")
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
