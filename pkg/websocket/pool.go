package websocket

import (
	log "github.com/sirupsen/logrus"
)

// Pool is used to broadcast a message to a pool of connected clients
type Pool struct {
	Clients    map[*Client]bool
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan Message
}

// NewPool creates a new pool
func NewPool() *Pool {
	return &Pool{
		Clients:    make(map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan Message),
	}
}

// Start activates the pool
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.WithFields(log.Fields{
				"clients": len(pool.Clients),
			}).Debug("Client successfully registered.")
			break

		case client := <-pool.Unregister:
			delete(pool.Clients, client)
			log.WithFields(log.Fields{
				"clients": len(pool.Clients),
			}).Debug("Client successfully unregistered.")
			break

		case message := <-pool.Broadcast:
			log.WithFields(log.Fields{
				"messageBody": message.Body,
			}).Debug("Broadcasting message.")
			for client := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.WithFields(log.Fields{
						"clientID":    client.ID,
						"messageBody": message.Body,
						"err":         err,
					}).Error("Error broadcasting message.")
					return
				}
			}
		}
	}
}
