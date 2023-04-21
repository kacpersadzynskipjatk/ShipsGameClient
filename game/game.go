package game

import "main/game_client"

type Game struct {
	GameClient game_client.GameClient
	Token string
}

func NewGame(c *game_client.GameClient) *Game {
	game := &Game{
		GameClient: *c,
		Token: c.PostStartGame(nil),
	}
	return game
}

func NewGameParams(c *game_client.GameClient, coords, desc, nick, targetNick string, wpbot bool) *Game {
	params := map[string]any {
			"coords": coords,
			"desc": desc,
			"nick": nick,
			"target_nick": targetNick,
			"wpbot": wpbot,
	}
	game := &Game{
		GameClient: *c,
		Token: c.PostStartGame(params),
	}
	return game
}

func (g *Game)StartGame() {
	//main game loop
}

func (g *Game)DisplayBoard() {
	
}
