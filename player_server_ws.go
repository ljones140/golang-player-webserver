package poker

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type playerServerWS struct {
	*websocket.Conn
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading to connection to WebSockets %v\n", err)
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, message, err := w.ReadMessage()

	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(message)
}