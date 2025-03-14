package main

import (
	"net/http"

	"github.com/NelsonJs/chat/internal/chat"
)

func main() {

	hub := chat.NewHub()
	go hub.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		println("add client")
		chat.ServeWs(hub, w, r)
	})
	println("ws server start")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
