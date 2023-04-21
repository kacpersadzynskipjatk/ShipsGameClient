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
	//newGame := game.NewGame(gc)
	coords := []string{
		"F1", "F2", "F3", "F4",
		"A1", "B1", "C1",
		"A3", "B3", "C3",
		"A5", "B5",
		"A7", "B7",
		"A9", "B9",
		"J9",
		"J7",
		"J5",
		"J3",
	}
	newGame := game.NewGameParams(gc, coords, "desc", "Kacper", "", true)
	print(newGame.Token)
	newGame.DisplayBoard()
}
