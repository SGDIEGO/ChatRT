package server

import (
	"fmt"
	"log"
	"net"

	"github.com/SGDIEGO/ChatRT/internal/client"
)

type Server struct {
	host, port string
}

func NewServer(host_, port_ string) *Server {
	return &Server{
		host: host_,
		port: port_,
	}
}

func (s *Server) Load() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.host, s.port))
	if err != nil {
		return err
	}

	defer listener.Close()

	log.Printf("Running server on %s:%s", s.host, s.port)
	for {
		// Verify client connection
		newConn, err := listener.Accept()
		if err != nil {
			return err
		}

		// New client
		newClient := client.NewClient(newConn)

		log.Printf("New connection: %d, Total connection: %d", newClient.Id, len(client.Clients))
		// Listen client activity
		go newClient.Listen()
	}
}
