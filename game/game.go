package game

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/grupawp/warships-lightgui/v2"
	"main/game_client"
	"time"
)

type Game struct {
	GameClient game_client.GameClient
	Token      string
	StatusResp game_client.StatusResponse
	BoardResp  game_client.BoardResponse
	FireResp   game_client.FireResponse
	Board      *board.Board
}

func NewGame(c *game_client.GameClient) *Game {
	r := game_client.StartGameRequest{
		Coords:     nil,
		Desc:       "",
		Nick:       "",
		TargetNick: "",
		Wpbot:      true,
	}
	resp := c.PostStartGame(&r)
	game := &Game{
		GameClient: *c,
		Token:      resp.Token,
	}
	return game
}

func NewGameParams(c *game_client.GameClient, coords []string, desc, nick, targetNick string, wpbot bool) *Game {
	r := game_client.StartGameRequest{
		Coords:     coords,
		Desc:       desc,
		Nick:       nick,
		TargetNick: targetNick,
		Wpbot:      wpbot,
	}
	resp := c.PostStartGame(&r)
	game := &Game{
		GameClient: *c,
		Token:      resp.Token,
	}
	return game
}

func (g *Game) StartGame() {
	for true {
		g.CheckGameStatus()
		if g.StatusResp.GameStatus == "end" {
			fmt.Println("Gra zakończona")
			//informacje o zwycięzcy
			time.Sleep(10 * time.Second)
			break
		}
		if g.StatusResp.GameStatus != "game_in_progress" {
			time.Sleep(1 * time.Second)
			continue
		}
		g.DisplayBoard()
		g.DisplayGameDescription()
		if !g.StatusResp.ShouldFire {
			time.Sleep(1 * time.Second)
			continue
		}

		//pobierz wsp przeciwnika
		oppShots := g.StatusResp.OppShots
		for _, shot := range oppShots {
			state, _ := g.Board.HitOrMiss(0, shot)
			g.Board.Set(0, shot, state)
		}
		for true {
			fmt.Println("Wykonaj ruch wpisz koordynaty:")
			var coordFromUser string
			fmt.Scanln(&coordFromUser)
			g.Fire(coordFromUser)
			state, _ := g.Board.HitOrMiss(1, coordFromUser)
			g.Board.Set(1, coordFromUser, state)
			if g.FireResp.Result != "Hit" {
				break
			}
		}

	}
}

func (g *Game) Fire(coord string) {
	r := game_client.FireRequest{
		Coord: coord,
	}
	resp := g.GameClient.PostFire(g.Token, &r)
	g.FireResp = resp
}

func (g *Game) CheckGameStatus() {
	g.StatusResp = g.GameClient.GetGameStatus(g.Token)
}

func (g *Game) DisplayGameDescription() {
	g.StatusResp = g.GameClient.GetGameDescription(g.Token)
	fmt.Printf("Gracz 1: %s\nOpis: %s\n", g.StatusResp.Nick, g.StatusResp.Desc)
	fmt.Printf("Gracz 2: %s\nOpis: %s\n", g.StatusResp.Opponent, g.StatusResp.OppDesc)
}

func (g *Game) DisplayBoard() {
	cfg := board.NewConfig()
	cfg.HitChar = '#'
	cfg.HitColor = color.FgRed
	cfg.BorderColor = color.BgRed
	cfg.RulerTextColor = color.BgBlue
	b := board.New(cfg)
	g.Board = b
	g.BoardResp = g.GameClient.GetGameBoards(g.Token)
	b.Import(g.BoardResp.Board)
	b.Display()
}
