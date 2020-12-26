package main

import (
	"fmt"
	"log"
	"os"

	poker "github.com/ljones140/golang-player-webserver"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {name} wins to record a win")

	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	game := poker.NewCLI(store, os.Stdin, poker.BlindAlerterFunc(poker.StdOutAlerter))

	game.PlayPoker()

}
