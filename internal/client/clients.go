package client

// Buffer for clients
type ClientsSTR []*Client

var Clients = make(ClientsSTR, 0)

// Find client
func (C *ClientsSTR) findClient(id int) *Client {
	for _, client := range *C {
		if client.Id == id {
			return client
		}
	}
	return nil
}
