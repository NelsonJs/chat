package chat

import "github.com/gorilla/websocket"

type Hub struct {
	onlines     map[int]*Client // 在线客户端
	room        map[int]*Room
	broadcaster chan []byte
}

func NewHub() *Hub {
	return &Hub{
		onlines:     make(map[int]*Client),
		broadcaster: make(chan []byte),
		room:        make(map[int]*Room),
	}
}

func (h *Hub) Run() {
	for message := range h.broadcaster {
		for _, client := range h.onlines {
			client.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
	close(h.broadcaster)
}

func (h *Hub) GetClient(id int) (*Client, bool) {
	client, ok := h.onlines[id]
	return client, ok
}

func (h *Hub) RemoveClient(id int) {
	delete(h.onlines, id)
}

func (h *Hub) AddClient(client *Client) {
	max := 0
	for k := range h.onlines {
		if k > max {
			max = k
		}
	}
	id := max + 1
	client.id = id
	h.onlines[id] = client
	println("current online: ", max)
}

func (h *Hub) AddRoom(room *Room) {
	h.room[room.id] = room
}

func (h *Hub) GetRoom(id int) (*Room, bool) {
	room, ok := h.room[id]
	return room, ok
}
