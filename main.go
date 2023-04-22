package main

import (
	"fmt"
	"main/game"
	"main/game_client"
	"net/http"
	"time"
)

func main() {
	c := http.Client{Timeout: 15 * time.Second}
	gc := game_client.NewGameClient(&c)
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
	//newGame := game.NewGame(gc)
	newGame := game.NewGameParams(gc, coords, "desc", "Kacper", "", true)
	newGame.DisplayBoard()
	fmt.Print(newGame.Token + "\n")
	newGame.Fire("A2")
}
