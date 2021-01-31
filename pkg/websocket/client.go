package websocket

import (
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		c.Conn.Close()
	}()

	for {
		// wait for a message from the connection
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Info("Error reading message.")
			return
		}

		message := Message{Type: messageType, Body: string(p)}
		log.WithFields(log.Fields{
			"MessageBody": message.Body,
		}).Debug("Message received.")

		// broadcast the message
		c.Pool.Broadcast <- message
	}
}
