package game

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/grupawp/warships-lightgui/v2"
	"main/game_client"
	"os"
	"text/tabwriter"
	"time"
)

func boardConfig() *board.Board {
	cfg := board.NewConfig()
	cfg.HitChar = '#'
	cfg.HitColor = color.FgRed
	cfg.BorderColor = color.BgRed
	cfg.RulerTextColor = color.BgBlue
	return board.New(cfg)
}

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
	b := boardConfig()
	game := &Game{
		GameClient: *c,
		Token:      resp.Token,
		Board:      b,
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
	b := boardConfig()
	game := &Game{
		GameClient: *c,
		Token:      resp.Token,
		Board:      b,
	}
	return game
}
func (g *Game) endGameCheck() {
	if g.StatusResp.GameStatus == "ended" {
		fmt.Println("Gra zakończona")
		if g.StatusResp.LastGameStatus == "win" {
			fmt.Printf("Wygrałeś %s gratulacje!!! ", g.StatusResp.Nick)
		} else if g.StatusResp.LastGameStatus == "lose" {
			fmt.Printf("Przegrałeś :( wygrał: %s", g.StatusResp.Opponent)
		}
		os.Exit(0)
	}
}

func (g *Game) setOpponentShots() {
	oppShots := g.StatusResp.OppShots
	for _, shot := range oppShots {
		state, _ := g.Board.HitOrMiss(board.Left, shot)
		err := g.Board.Set(board.Left, shot, state)
		if err != nil {
			err = fmt.Errorf("board not set: %s", err)
			fmt.Println(err)
		}
	}
}
func (g *Game) makeShot() {
	for true {
		fmt.Println("Twoja Tura!\nWpisz koordynaty by strzelić:")
		var coordFromUser string
		fmt.Scanln(&coordFromUser)
		g.Fire(coordFromUser)
		state, _ := g.Board.HitOrMiss(board.Right, coordFromUser)
		err := g.Board.Set(board.Right, coordFromUser, state)
		if err != nil {
			err = fmt.Errorf("board not set: %s", err)
			fmt.Println(err)
		}
		if g.FireResp.Result != "Hit" {
			break
		}
	}
}
func (g *Game) StartGame() {
	for true {
		g.checkGameStatus()
		g.endGameCheck()
		if g.StatusResp.GameStatus != "game_in_progress" {
			time.Sleep(1 * time.Second)
			continue
		}
		if !g.StatusResp.ShouldFire {
			time.Sleep(1 * time.Second)
			continue
		}
		g.displayBoard()
		g.displayGameDescription()
		g.setOpponentShots()
		g.makeShot()
	}
}

func (g *Game) Fire(coord string) {
	r := game_client.FireRequest{
		Coord: coord,
	}
	resp := g.GameClient.PostFire(g.Token, &r)
	g.FireResp = resp
}

func (g *Game) checkGameStatus() {
	g.StatusResp = g.GameClient.GetGameStatus(g.Token)
}

func (g *Game) displayGameDescription() {
	g.StatusResp = g.GameClient.GetGameDescription(g.Token)
	writer := tabwriter.NewWriter(os.Stdout, 9, 8, 0, '\t', 0)
	fmt.Fprintf(writer, "%s  VS.  %s\n\n", g.StatusResp.Nick, g.StatusResp.Opponent)
	fmt.Fprintf(writer, "Gracz\tOpis\n")
	fmt.Fprintf(writer, "%s\t%s\n", "-----", "----")
	fmt.Fprintf(writer, "%s\t%s\n", g.StatusResp.Nick, g.StatusResp.Desc)
	fmt.Fprintf(writer, "%s\t%s\n\n", g.StatusResp.Opponent, g.StatusResp.OppDesc)
	writer.Flush()
}

func (g *Game) displayBoard() {
	g.BoardResp = g.GameClient.GetGameBoards(g.Token)
	err := g.Board.Import(g.BoardResp.Board)
	if err != nil {
		err = fmt.Errorf("import error: %s", err)
		fmt.Println(err)
	}
	g.Board.Display()
}
