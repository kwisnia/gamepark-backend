package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// MessageType enum
type MessageType string

const (
	// ChatMessage type
	ChatMessage MessageType = "chatMessage"
)

type SocketMessage struct {
	MessageType MessageType            `json:"messageType"`
	Data        map[string]interface{} `json:"data"`
}

var ClientHub = newHub()

func CreateConnection(w http.ResponseWriter, r *http.Request, userID uint) error {
	fmt.Println(len(ClientHub.clients))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return err
	}
	client := &Client{hub: ClientHub, conn: conn, send: make(chan []byte, 256), userID: userID}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	return nil
}

func VerifyIncomingMessage(message []byte) (*SocketMessage, error) {
	parsedMessage := SocketMessage{}
	err := json.Unmarshal(message, &parsedMessage)
	if err != nil {
		return nil, err
	}
	fmt.Println(parsedMessage.MessageType)
	if parsedMessage.MessageType != ChatMessage {
		return nil, fmt.Errorf("invalid message type")
	}
	return &parsedMessage, nil
}
