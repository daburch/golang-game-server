package checkers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/daburch/golang-game-server/pkg/common"
	"github.com/daburch/golang-game-server/pkg/websocket"
	gorilla "github.com/gorilla/websocket"

	log "github.com/sirupsen/logrus"
)

var upgrader = gorilla.Upgrader{}
var game *Game
var messages []websocket.Message

func assertEquals(t *testing.T, a interface{}, b interface{}, message string) {
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}

	assert(t, a == b, message)
}

func assertNotEquals(t *testing.T, a interface{}, b interface{}, message string) {
	if len(message) == 0 {
		message = fmt.Sprintf("%v == %v", a, b)
	}

	assert(t, a != b, message)
}

func assert(t *testing.T, cond bool, message string) {
	if cond {
		return
	}

	t.Fatal(message)
}

// Test_game get a game and has 2 players join
// 		- the game should be created and added to the queue when player 1 joins.
func Test_game(t *testing.T) {
	common.InitLogger()
	game = GetGame()

	// Create test server with with ws func to join the game
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()

		game.Join(c)
	}))
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.1
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	assert(t, game.Player1 == nil && game.Player2 == nil, "player1 and player2 should be nil before anyone joins")
	assert(t, !game.IsFull(), "game should not be full now - no players present")

	// Connect player1 to the server
	player1 := connectToGame(t, u)
	defer player1.Close()
	time.Sleep(1 * time.Second)

	assert(t, game.Player1 != nil, "player1 didn't successfully register to the game")
	assert(t, !game.IsFull(), "game should not be full now - 1 player present")

	// Connect player2 to the server
	player2 := connectToGame(t, u)
	defer player2.Close()
	time.Sleep(1 * time.Second)

	assert(t, game.Player2 != nil, "player2 didn't successfully register to the game")
	assert(t, game.IsFull(), "game should be full now")

	// check assignColor
	assert(t, len(messages) == 2, "assignColor messages were not recieved.")

	msg0 := websocket.ParseMessage(&messages[0])
	msg1 := websocket.ParseMessage(&messages[1])

	// either player could recieve the action first so don't worry about checking color here
	assert(t, msg0["action"] == "assignColor", "player1 didn't recieve assignColor message on start.")
	assert(t, msg1["action"] == "assignColor", "player2 didn't recieve assignColor message on start.")

	quitAction := "{ \"action\": \"quit\" }"
	// message := websocket.Message{Type: 1, Body: quitAction}

	// player1 quits
	player1.WriteMessage(1, []byte(quitAction))
	time.Sleep(1 * time.Second)

	assert(t, game.Player1 == nil, "player1 wasn't able to quit from the game")
	assert(t, !game.IsFull(), "game should not be full now")

	// player 2 disconnect
	player2.Close()
	time.Sleep(1 * time.Second)

	assert(t, game.Player1 == nil, "player2 wasn't able to disconnect from the game")
	assert(t, !game.IsFull(), "game should not be full now")

	// finished!
	log.Info("Done")
}

func connectToGame(t *testing.T, url string) *gorilla.Conn {
	// Connect to the server
	ws, _, err := gorilla.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}

	go read(t, ws)

	return ws
}

func read(t *testing.T, ws *gorilla.Conn) {
	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			return
		}

		message := websocket.Message{Type: messageType, Body: string(p)}
		messages = append(messages, message)
	}
}
