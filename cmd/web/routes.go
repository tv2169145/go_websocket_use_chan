package main

import (
	"github.com/bmizerany/pat"
	"github.com/tv2169145/go_websocket_use_chan/internal/homeHandler"
	"github.com/tv2169145/go_websocket_use_chan/internal/roomHandler"
	"net/http"
)

func routes() http.Handler {
	mux := pat.New()
	roomHandler.NewRoom()
	mux.Get("/", http.HandlerFunc(homeHandler.Home))
	mux.Get("/broadcast", http.HandlerFunc(roomHandler.Broadcast))
	mux.Get("/ws", http.HandlerFunc(roomHandler.ChatRoomHandle))

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return mux
}
