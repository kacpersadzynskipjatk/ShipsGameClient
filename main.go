package main

import (
	"main/game"
	"main/game_client"
	"net/http"
	"time"
)

func main() {
	c := http.Client{Timeout: 5 * time.Second}
	gc := game_client.NewGameClient(&c)
	game := game.NewGame(gc)
	game.StartGame()
}
