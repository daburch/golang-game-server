package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/daburch/golang-game-server/pkg/checkers"
	"github.com/daburch/golang-game-server/pkg/websocket"
)

func main() {
	log.SetLevel(log.InfoLevel)

	log.Info("Game Server starting.")

	handleRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRoutes() {
	log.Debug("creating route: /ws")
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(w, r)
	})
}

func serveWs(w http.ResponseWriter, r *http.Request) {
	log.Debug("WebSocket Endpoint Hit.")

	ws, err := websocket.Upgrade(w, r)
	if err != nil {
		log.WithFields(log.Fields{
			"w":   w,
			"r":   r,
			"err": err,
		}).Fatal("Error upgrading connection.")
	}

	log.Debug("Finding Game.")
	game := checkers.GetGame()

	log.WithFields(log.Fields{
		"gameID": game.GameID,
	}).Info("Joining game.")
	game.Join(ws)
}
