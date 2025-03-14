package model

import "encoding/json"

type MessageType int

const (
	_                = iota
	Text MessageType = iota
	Image
	File
)

type Message struct {
	Id          int         `json:"id"`      // 消息id
	UserId      int         `json:"userId"`  // 发送者id
	PeerId      int         `json:"peerId"`  // 接收者id
	GroupId     int         `json:"groupId"` // 群id
	Content     string      `json:"content"`
	ContentType MessageType `json:"contentType"`
}

func (m *Message) Unmarshal(bytes []byte) error {
	return json.Unmarshal(bytes, m)
}
