package main

import (
	"github.com/tv2169145/go_websocket_use_chan/internal/roomHandler"
	"log"
	"net/http"
)

func main() {
	mux := routes()
	//roomHandler.StartChatRoom()

	go roomHandler.ChatRoom.ProcessTask()

	log.Println("start server on port 8081...")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Println(err)
		return
	}
}
