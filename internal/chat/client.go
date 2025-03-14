package chat

import (
	"log"
	"net/http"
	"time"

	"github.com/NelsonJs/chat/internal/plugins/objpool"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	id   int
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
	room map[int]struct{}
}

func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println(err.Error())
		return
	}
	client := &Client{
		conn: conn,
		hub:  hub,
		send: make(chan []byte),
		room: make(map[int]struct{}),
	}
	hub.AddClient(client)
	go client.Read()
	go client.Write()
}

func (c *Client) Read() {
	defer func() {
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		messageType, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			c.hub.RemoveClient(c.id)
			break
		}
		if messageType != websocket.TextMessage {
			continue
		}
		msgObj := objpool.Pool().GetMessage()
		if err = msgObj.Unmarshal(message); err != nil {
			log.Printf("UnmarshalJSON error: %v", err)
			continue
		}
		switch {
		case msgObj.GroupId > 0:
			if _, ok := c.room[msgObj.GroupId]; !ok {
				continue
			}
			room, ok := c.hub.GetRoom(msgObj.GroupId)
			if !ok {
				continue
			}
			room.Send(message)
		case msgObj.PeerId > 0:
			println("read")
			if peer, ok := c.hub.GetClient(msgObj.PeerId); ok {
				println("finded client")
				peer.send <- message
				println("send msg")
			}
		}
	}
}

func (c *Client) Write() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		println("write for")
		select {
		case message, ok := <-c.send:
			println("write msg")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				continue
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
