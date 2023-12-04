package client

type Room struct {
	id      int
	name    string
	clients []*Client
}

// Buffer for rooms
var Rooms []*Room

func NewRoom(name_ string) *Room {
	var newRoom *Room
	newRoom.id = len(Rooms) + 1
	newRoom.name = name_

	Rooms = append(Rooms, newRoom)

	return newRoom
}

func (r *Room) addClient(client *Client) {
	r.clients = append(r.clients, client)
}
