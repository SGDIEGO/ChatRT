package client

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

var (
	END = string([]byte{13, 10})
)

type Client struct {
	Id         int
	Name       string
	Connection net.Conn
	Actual     *Room
	Rooms      []*Room
}

func NewClient(connection_ net.Conn) *Client {
	// Create client
	id_ := len(Clients) + 1
	var newClient = &Client{
		Id:         id_,
		Name:       fmt.Sprintf("Random %d", id_),
		Connection: connection_,
		Actual:     nil,
		Rooms:      make([]*Room, 0),
	}

	// Add client
	Clients = append(Clients, newClient)

	// Return client created
	return newClient
}

func (c *Client) Listen() {
	defer c.Connection.Close()

	for {
		// Read message from connection
		message, err := bufio.NewReader(c.Connection).ReadString('\n')
		if err != nil {
			log.Fatalf("%s", err.Error())
		}

		// Delete last chars
		message = strings.TrimSuffix(message, "\r\n")

		// Verify commands
		switch message {
		// Command exit
		case "/exit":
			c.Connection.Close()
			return

		// Command create
		case "/create":
			c.Connection.Write([]byte("Name for room: "))
			roomName, err := bufio.NewReader(c.Connection).ReadString('\n')
			if err != nil {
				c.Connection.Write([]byte("Error creating room: " + err.Error()))
				break
			}

			// Room created
			newRoom := NewRoom(roomName)
			newRoom.addClient(c)
			c.addRoom(newRoom)
			c.Write("Room " + roomName + " created!")

		// Command for list rooms
		case "/list":
			// If dont have any room registered
			if len(c.Rooms) == 0 {
				c.Write("You dont have any room registered")
				break
			}

			// Else show rooms
			for _, room := range c.Rooms {
				c.Write(room.name)
			}

		// Command for enter to room
		case "/join":
			c.Write("Name of room: ")
			roomName, err := bufio.NewReader(c.Connection).ReadString('\n')
			if err != nil {
				c.Write("Error creating room: " + err.Error())
				break
			}

			// Find room and set as actual to client
			room := c.findRoom(roomName)
			if room == nil {
				c.Write("You dont have registered room with name " + roomName)
				break
			}
			c.Actual = room
			c.Write("You are in room " + roomName + " right now!")

		// Command for show users
		case "/users":
			c.showClients()

		// Command for invite user to room
		case "/invite":
			if c.Actual == nil {
				c.Write("Join any room first")
				break
			}
			c.Write("Id user: ")
			input, err := bufio.NewReader(c.Connection).ReadString('\n')
			if err != nil {
				c.Write("Error reading input, " + err.Error())
				break
			}
			input = strings.TrimSuffix(input, "\r\n")

			// Convert input string to int
			id, err := strconv.Atoi(input)
			if err != nil {
				c.Write("Error with id, " + err.Error())
				break
			}

			// Find user from list
			friend := Clients.findClient(id)

			// Add new friend to room
			c.Actual.addClient(friend)
			friend.Rooms = append(friend.Rooms, c.Actual)
			c.Write(fmt.Sprintf("%s has been adding to room %s", friend.Name, c.Actual.name))

		// By default send message to all users
		default:
			// Send message to all clients
			c.Send(message)
		}
	}

}

func (c *Client) Write(message string) error {
	_, err := c.Connection.Write([]byte(message + "\n"))
	return err
}

func (c *Client) findRoom(name string) *Room {

	for _, room := range c.Rooms {
		if strings.Compare(room.name, name) == 0 {
			return room
		}
	}

	return nil
}

func (c *Client) addRoom(r *Room) {
	c.Rooms = append(c.Rooms, r)
}

func (c *Client) Send(message string) {
	message = fmt.Sprintf("(%s): %s", c.Name, message)
	// Send message to all
	if c.Actual == nil {
		// Iterate all clients
		for _, client := range Clients {
			// Verify if client id is equal
			if client.Id != c.Id && client.Actual == nil {
				// Verify if client is in same room
				client.Write(message)
			}
		}
		return
	}

	// Else send to room users
	c.Actual.sendMessage(message)
}

func (c *Client) showClients() {
	for _, client := range Clients {
		if c.Id != client.Id {
			c.Write(fmt.Sprintf("User %s with id: %d", client.Name, client.Id))
		}
	}
}
