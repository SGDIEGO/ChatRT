package client

type Room struct {
	id      int
	name    string
	clients []*Client
}

// Buffer for rooms
var Rooms []*Room

func NewRoom(name_ string) *Room {
	var newRoom = &Room{
		id:   len(Rooms) + 1,
		name: name_,
	}

	Rooms = append(Rooms, newRoom)

	return newRoom
}

func (r *Room) addClient(client *Client) {
	r.clients = append(r.clients, client)
}

func (r *Room) sendMessage(message string) {
	for _, client := range r.clients {
		client.Write(message)
	}
}
