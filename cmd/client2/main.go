package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/NelsonJs/chat/internal/model"
	"github.com/gorilla/websocket"
)

func main() {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		println("err")
		return
	}

	go func() {
		_, msgBte, err := conn.ReadMessage()
		if err != nil {
			println(err.Error())
			return
		}
		item := model.Message{}
		json.Unmarshal(msgBte, &item)
		fmt.Printf("%d receive msg: %+v", 2, item)
	}()

	time.Sleep(8 * time.Second)
	msg := model.Message{
		Content:     "world",
		ContentType: model.Text,
		Id:          2,
		PeerId:      1,
	}
	bteData, _ := json.Marshal(msg)
	err = conn.WriteMessage(websocket.TextMessage, bteData)
	if err != nil {
		println(err.Error())
		return
	}

	time.Sleep(time.Minute * 10)
}
