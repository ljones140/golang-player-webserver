package poker_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	poker "github.com/ljones140/golang-player-webserver"
)

func TestGETPlayers(t *testing.T) {
	store := poker.StubPlayerStore{
		Scores: map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}

	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseBody(t, response.Body.String(), "20")
		assertStatus(t, response, http.StatusOK)
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := poker.NewGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		poker.AssertResponseBody(t, response.Body.String(), "10")
		assertStatus(t, response, http.StatusOK)
	})
	t.Run("returns 404 when player does not exist", func(t *testing.T) {
		request := poker.NewGetScoreRequest("NonExistantPlayer")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusNotFound)
	})
}

func TestLeague(t *testing.T) {

	t.Run("it returns the league table as JSON", func(t *testing.T) {
		wantedLeague := []poker.Player{
			{"Cleo", 32},
			{"Cleo", 20},
			{"Cleo", 14},
		}

		store := poker.StubPlayerStore{nil, nil, wantedLeague}
		server := mustMakePlayerServer(t, &store, dummyGame)

		request := poker.NewGetLeagueRequest()
		response := httptest.NewRecorder()
		server.ServeHTTP(response, request)
		got := poker.GetLeagueFromResponse(t, response.Body)

		poker.AssertLeague(t, got, wantedLeague)
		assertStatus(t, response, http.StatusOK)
		poker.AssertContentType(t, response, "application/json")
	})
}

func TestStoreWins(t *testing.T) {
	store := poker.StubPlayerStore{
		Scores: map[string]int{},
	}

	server := mustMakePlayerServer(t, &store, dummyGame)

	t.Run("it records win when POST", func(t *testing.T) {
		player := "Pepper"
		request := poker.NewPostWinRequest(player)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusAccepted)

		//TODO: use helper here
		if len(store.WinCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.WinCalls), 1)
		}

		if store.WinCalls[0] != player {
			t.Errorf("did not store correct winner go %q want %q", store.WinCalls[0], player)
		}
	})
}

func TestGame(t *testing.T) {
	t.Run("GET /game returns a 200", func(t *testing.T) {
		server := mustMakePlayerServer(t, &poker.StubPlayerStore{}, dummyGame)
		request := poker.NewGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)
		assertStatus(t, response, http.StatusOK)

	})

	t.Run("when we aget a message over a websocket it is a winner of a game", func(t *testing.T) {
		store := &StubPlayerStore{}
		winner := "Ruth"

		server := httptest.NewServer(mustMakePlayerServer(t, store))
		defer server.Close()

		wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

		conn := mustDialWS(t, wsURL)

		defer conn.Close()

		writeWSMessage(t, conn, winner)

		time.Sleep(10 * time.Millisecond)
		AssertPlayerWin(t, store, winner)

	})
}

func assertStatus(t testing.TB, response *httptest.ResponseRecorder, want int) {
	t.Helper()
	got := response.Code
	if got != want {
		t.Errorf("did not get correct status, got %d want %d", got, want)
	}
}

func mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}

	return server
}

func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}

func writeWSMessage(t *testing.T, conn *websocket.Conn, message string) {
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
