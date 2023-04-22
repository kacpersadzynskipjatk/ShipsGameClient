package game

import (
	"encoding/json"
	"main/game_client"

	"github.com/fatih/color"
	"github.com/grupawp/warships-lightgui/v2"
)


type Game struct {
	GameClient game_client.GameClient
	Token      string
	Status     string
	Board      board.Board
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

func (g *Game) Fire(coord string) {
	g.CheckGameStatus()
	
	//make a object to unmashall all responses???

	resp := g.GameClient.PostFire(coord, g.Token)
	//shotResult, _ := unmarshalToSlice(resp)
	print(resp)
	//g.Board.Set("Right",)
	//g.Board.Display()
}

func (g *Game) CheckGameStatus() {
	g.Status = g.GameClient.GetGameStatus(g.Token)
}

func (g *Game) DisplayBoard() {
	cfg := board.NewConfig()
	cfg.HitChar = '#'
	cfg.HitColor = color.FgRed
	cfg.BorderColor = color.BgRed
	cfg.RulerTextColor = color.BgBlue
	board := board.New(cfg)

	jsonShips := g.GameClient.GetGameBoards(g.Token)
	coords, _ := unmarshalToSlice(jsonShips)

	board.Import(coords)
	board.Display()
	g.Board = *board
}

func unmarshalToSlice(jsonStr string) ([]string, error) {
	var data map[string][]string
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil, err
	}

	board := data["board"]

	return board, nil
}