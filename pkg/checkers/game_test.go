package checkers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

func Test_game(t *testing.T) {
	log.SetLevel(log.DebugLevel)

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

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect player1 to the server
	player1 := connectToGame(t, u)
	defer player1.Close()
	time.Sleep(1 * time.Second)
	assert(t, game.Player1 != nil, "player1 didn't successfully register to the game")

	// Connect player2 to the server
	player2 := connectToGame(t, u)
	defer player2.Close()
	time.Sleep(1 * time.Second)
	assert(t, game.Player2 != nil, "player2 didn't successfully register to the game")

	time.Sleep(5 * time.Second)

	// check assignColor
	assert(t, len(messages) == 2, "assignColor messages were not recieved.")

	msg0 := parseMessage(&messages[0])
	msg1 := parseMessage(&messages[1])

	assert(t, msg0["action"] != nil, "player1 didn't recieve assignColor message on start.")
	assert(t, msg1["action"] != nil, "player2 didn't recieve assignColor message on start.")

	time.Sleep(1 * time.Second)

	// msg = read(t, game.Player2.Conn)
	// assert(t, msg["action"] == "assignColor", "player2 didn't recieve assignColor message on start.")
	// time.Sleep(1 * time.Second)

	time.Sleep(10 * time.Second)

	log.Debug("Done")
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
			t.Fatalf("%v", err)
		}

		message := websocket.Message{Type: messageType, Body: string(p)}
		log.Debug(message)
		messages = append(messages, message)
	}
}

func parseMessage(message *websocket.Message) map[string]interface{} {
	var result map[string]interface{}
	var body map[string]interface{}

	json.Unmarshal([]byte(message.Body), &result)
	json.Unmarshal([]byte(fmt.Sprintf("%v", result["body"])), &body)

	return body
}
