package checkers

import (
	"fmt"

	"github.com/daburch/golang-game-server/pkg/websocket"
	gorilla "github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

// Game represents a checkers game with 2 players
type Game struct {
	GameID  int
	Pool    *websocket.Pool
	Player1 *websocket.Client
	Player2 *websocket.Client
}

var gameCounter = 0

// AvailableGames is the slice of games currently waiting for a second player
var AvailableGames []*Game

// ActiveGames is the slice of games currently being played
var ActiveGames []*Game

// GetGame returns the next available game
func GetGame() *Game {
	var game *Game
	if len(AvailableGames) != 0 {
		game = AvailableGames[0]
		AvailableGames = AvailableGames[1:]
	} else {
		game = NewGame()
		AvailableGames = append(AvailableGames, game)
	}

	return game
}

// NewGame initalizes a new game with the next available gameID
func NewGame() *Game {
	pool := websocket.NewPool()
	go pool.Start()

	return &Game{
		GameID: getNextID(),
		Pool:   pool,
	}
}

// getNextID returns the next available gameID
func getNextID() int {
	gameCounter++
	return gameCounter
}

// Join the game
func (game *Game) Join(ws *gorilla.Conn) {
	var client *websocket.Client
	if game.Player1 == nil {
		// create the client for player1
		client = &websocket.Client{
			ID:   "player1",
			Conn: ws,
			Pool: game.Pool,
		}

		// 1 / 2 players in the game. wait for another person to join.
		log.WithFields(log.Fields{
			"gameID": game.GameID,
		}).Debug("Player1 joined.")
		game.Player1 = client
	} else if game.Player2 == nil {
		// create the client for player2
		client = &websocket.Client{
			ID:   "player2",
			Conn: ws,
			Pool: game.Pool,
		}

		// game is now full, start the game
		log.WithFields(log.Fields{
			"gameID": game.GameID,
		}).Debug("Player2 joined.")
		game.Player2 = client

		game.Start()
	}

	// register the client to the game pool
	game.Pool.Register <- client

	// listen for messages
	client.Read()
}

// IsFull checks if the game is full
func (game *Game) IsFull() bool {
	return game.Player1 != nil && game.Player2 != nil
}

// Start starts the game
func (game *Game) Start() {
	log.WithFields(log.Fields{
		"gameID": game.GameID,
	}).Info("Starting game.")

	assignColor(game.Player1, "white")
	assignColor(game.Player2, "black")

	// add the game to the active game queue
	ActiveGames = append(ActiveGames, game)
}

func assignColor(player *websocket.Client, color string) {
	p := fmt.Sprintf(`{ "action": "assignColor", "color": "%s" }`, color)
	message := websocket.Message{Type: 1, Body: string(p)}
	player.Conn.WriteJSON(message)
}
