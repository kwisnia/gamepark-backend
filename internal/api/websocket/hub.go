package websocket

import "fmt"

type Hub struct {
	// Registered clients.
	clients map[uint]*Client

	// Inbound messages from the clients.
	Send chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

type Message struct {
	ID   uint
	Data []byte
}

func newHub() *Hub {
	return &Hub{
		Send:       make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[uint]*Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			fmt.Println("register")
			h.clients[client.userID] = client
		case client := <-h.unregister:
			fmt.Println("unregister")
			if _, ok := h.clients[client.userID]; ok {
				delete(h.clients, client.userID)
				close(client.send)
			}
		case message := <-h.Send:
			if client, ok := h.clients[message.ID]; ok {

				select {
				case client.send <- message.Data:

				default:
					close(client.send)
					delete(h.clients, client.userID)
				}
			}
		}
	}
}
