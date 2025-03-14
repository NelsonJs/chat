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
		fmt.Printf("%d receive msg: %+v", 1, item)
	}()

	time.Sleep(5 * time.Second)
	msg := model.Message{
		Content:     "hello",
		ContentType: model.Text,
		Id:          1,
		PeerId:      2,
	}
	bteData, _ := json.Marshal(msg)
	err = conn.WriteMessage(websocket.TextMessage, bteData)
	if err != nil {
		println(err.Error())
		return
	}

	time.Sleep(time.Minute * 10)
}
