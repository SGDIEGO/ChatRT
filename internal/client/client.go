package client

import (
	"bufio"
	"log"
	"net"
)

var (
	END = string([]byte{13, 10})
)

type Client struct {
	Id         int
	Connection net.Conn
}

// Buffer for clients
var Clients []*Client

func NewClient(id_ int, connection_ net.Conn) *Client {
	// Create client
	var newClient *Client
	newClient.Id = id_
	newClient.Connection = connection_

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

		// Verify commands
		switch message {
		// Command exit
		case "/exit" + END:
			c.Connection.Close()
			return
		// Command create
		case "/create new room" + END:
			c.Connection.Write([]byte("Name for room: "))
			roomName, err := bufio.NewReader(c.Connection).ReadString('\n')
			if err != nil {
				c.Connection.Write([]byte("Error creating room: " + err.Error()))
				break
			}

			// Room created
			newRoom := NewRoom(roomName)
			newRoom.addClient(c)

		default:
			// Send message to all clients
			err = c.Send(message)
			if err != nil {
				log.Printf("Error writing message %s from %d", message, c.Id)
			}
		}
	}

}

func (c *Client) Send(message string) error {
	var ERROR error
	// Iterate all clients
	for _, client := range Clients {
		// Verify if client id is equal
		if client.Id != c.Id {
			// Write message to client
			_, ERROR = client.Connection.Write([]byte(message))
		}
	}

	// Return any error
	return ERROR
}
