package main

import "github.com/SGDIEGO/ChatRT/internal/server"

const (
	HOST = "localhost"
	PORT = "3000"
)

var SERVER = server.NewServer(HOST, PORT)

func main() {
	SERVER.Load()
}
