package game

import (
	"encoding/json"
	"fmt"
	"main/game_client"

	"github.com/fatih/color"
	"github.com/grupawp/warships-lightgui/v2"
)


type Game struct {
	GameClient game_client.GameClient
	Token      string
	Status     string
}

func NewGame(c *game_client.GameClient) *Game {
	game := &Game{
		GameClient: *c,
		Token:      c.PostStartGame(nil),
	}
	return game
}

func NewGameParams(c *game_client.GameClient, coords []string, desc, nick, targetNick string, wpbot bool) *Game {
	params := map[string]any{
		"coords":      coords,
		"desc":        desc,
		"nick":        nick,
		"target_nick": targetNick,
		"wpbot":       wpbot,
	}
	game := &Game{
		GameClient: *c,
		Token:      c.PostStartGame(params),
	}
	return game
}

func (g *Game) StartGame() {
	//main game loop
}

func (g *Game) CheckGameStatus() {
	g.Status = g.GameClient.GetGameStatus(g.Token)
}

func (g *Game) DisplayBoard() {
	cfg := board.NewConfig()
	cfg.HitChar = '*'
	cfg.HitColor = color.FgRed
	cfg.BorderColor = color.BgRed
	cfg.RulerTextColor = color.BgBlue
	board := board.New(cfg)
	board.Display()
	
	jsonShips := g.GameClient.GetGameBoards(g.Token)
	slice, _ := getBoardSlice(jsonShips)
	fmt.Println(slice)
	board.Import(slice)
}

func getBoardSlice(jsonStr string) ([]string, error) {
	var data map[string][]string
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}

	board := data["board"]

	// Create a new slice of strings in the desired format
	var boardSlice []string
	for _, s := range board {
		boardSlice = append(boardSlice, fmt.Sprintf(`"%s"`, s))
	}

	return boardSlice, nil
}