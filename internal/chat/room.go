package chat

type RoomStatus uint8

const (
	Normal  RoomStatus = iota
	Silence            // 禁言
)

type Room struct {
	id     int
	users  []int
	hub    *Hub
	status RoomStatus
}

func NewRoom(id int, users []int, status RoomStatus) *Room {
	return &Room{
		id:     id,
		users:  users,
		status: Normal,
	}
}

func (r *Room) Send(msg []byte) {
	if r.status == Silence {
		return
	}
	for _, v := range r.users {
		if client, ok := r.hub.GetClient(v); ok {
			client.send <- msg
		}
	}
}
