package main

import (
	"log"
	"net/http"
	"os"

	poker "github.com/ljones140/golang-player-webserver"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, nil := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("Problem creating file system player store, %v ", err)

	}
	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listn on port 5000 %v", err)
	}
}