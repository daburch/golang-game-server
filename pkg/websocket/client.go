package websocket

import (
	"encoding/json"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

// Client represents a game connection with identification and pool details
type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

// Message is a message representing a game action
type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	quitRecieved := false
	for !quitRecieved {
		// wait for a message from the connection
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			if strings.HasPrefix(err.Error(), "websocket: close") {
				log.WithFields(log.Fields{
					"id":  c.ID,
					"err": err,
				}).Debug("Disconnect.")
			} else {
				log.WithFields(log.Fields{
					"err": err,
				}).Error("Error reading message.")
			}
			return
		}

		message := Message{Type: messageType, Body: string(p)}
		log.WithFields(log.Fields{
			"MessageBody": message.Body,
		}).Debug("Message received.")

		msg := ParseMessage(&message)
		if msg["action"] == "quit" {
			quitRecieved = true
		} else {
			// broadcast the message
			c.Pool.Broadcast <- message
		}
	}
}

// ParseMessage unmarshals a message body into a map
func ParseMessage(message *Message) map[string]interface{} {
	var result map[string]interface{}

	json.Unmarshal([]byte(message.Body), &result)

	return result
}
