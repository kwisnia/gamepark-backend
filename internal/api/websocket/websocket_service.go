package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// MesssageType enum
type MesssageType string

const (
	// ChatMessage type
	ChatMessage MesssageType = "chatMessage"
)

type SocketMessage struct {
	messageType MesssageType
	data        map[string]any
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

func VerifyIncomingMessage(message []byte) error {
	parsedMessage := SocketMessage{}
	err := json.Unmarshal(message, &parsedMessage)
	if err != nil {
		return err
	}
	if parsedMessage.messageType != ChatMessage {
		return fmt.Errorf("invalid message type")
	}
	return nil
}
