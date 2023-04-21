package game

import "main/game_client"

type Game struct {
	GameClient game_client.GameClient
}

func NewGame(c *game_client.GameClient) *Game {
	game := &Game{
		GameClient: *c,
	}
	return game
}
