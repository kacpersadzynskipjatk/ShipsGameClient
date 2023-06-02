package main

import (
	"main/application"
	"main/client"
	"net/http"
	"time"
)

func main() {
	httpClient := http.Client{Timeout: 15 * time.Second}
	gameClient := client.NewClient(&httpClient)
	newGame := application.NewGame(gameClient)
	newGame.StartGame()
}
