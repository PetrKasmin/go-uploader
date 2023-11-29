package handlers

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var (
	upgradeConnection = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	clients      = make(map[WebSocketConnection]string)
	clientsMutex = sync.Mutex{}
	wsChan       = make(chan float64)
)

type WebSocketConnection struct {
	*websocket.Conn
}

func UploadWebsocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgradeConnection.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	conn := WebSocketConnection{Conn: ws}
	clients[conn] = ""

	go handleClient(&conn)
}

func ListenToWsChannel() {
	for {
		e := <-wsChan
		broadcastProgress(e)
	}
}

func broadcastProgress(progress float64) {
	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%f.2", progress)))
		if err != nil {
			log.Println(err)
		}
	}
}

func handleClient(conn *WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
		}
	}()

	var totalBytesRead int
	buf := make([]byte, chunkSize)

	var payload struct {
		Filename string `json:"filename"`
		Filesize int64  `json:"filesize"`
	}

	err := conn.ReadJSON(&payload)
	if err != nil {
		return
	}

	fileName := filepath.Join(uploadPath, payload.Filename)
	out, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer out.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if messageType == websocket.BinaryMessage {
			n := copy(buf, p)

			_, err := out.Write(buf[:n])
			if err != nil {
				return
			}

			totalBytesRead += len(p)

			progress := float64(totalBytesRead) / float64(payload.Filesize) * 100
			wsChan <- progress
		}
	}
}
