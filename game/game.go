package game

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/grupawp/warships-lightgui/v2"
	"log"
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
	GameDesc   game_client.StatusResponse
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
	start, err := c.PostStartGame(&r)
	if err != nil {
		log.Fatalln(err)
	}
	b := boardConfig()
	game := &Game{
		GameClient: *c,
		Token:      start.Token,
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
	start, err := c.PostStartGame(&r)
	if err != nil {
		log.Fatalln(err)
	}
	b := boardConfig()
	game := &Game{
		GameClient: *c,
		Token:      start.Token,
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
		state, err := g.Board.HitOrMiss(board.Left, shot)
		if err != nil {
			log.Fatalln(err)
		}
		err = g.Board.Set(board.Left, shot, state)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
func getValidCoords() string {
	var coordFromUser string
	for {
		fmt.Scanln(&coordFromUser)
		chars := []rune(coordFromUser)
		if !(string(chars[0]) >= "A" && string(chars[0]) <= "J") {
			fmt.Println("Niepoprawna litera, wpisz jeszcze raz:")
			continue
		} else if len(chars) == 2 {
			if !(string(chars[1]) >= "1" && string(chars[1]) <= "9") {
				fmt.Println("Niepoprawny numer, wpisz jeszcze raz:")
				continue
			}
			break
		} else if len(chars) == 3 {
			if !(string(chars[1]) == "1" && string(chars[2]) == "0") {
				fmt.Println("Niepoprawny numer, wpisz jeszcze raz:")
				continue
			}
			break
		} else {
			fmt.Println("Niepoprawna długość koordynatów")
			continue
		}
	}
	return coordFromUser
}

func (g *Game) makeShot() {
	for {
		fmt.Println("Twoja Tura!\nWpisz koordynaty by strzelić:")
		coord := getValidCoords()
		g.Fire(coord)
		state, err := g.Board.HitOrMiss(board.Right, coord)
		if err != nil {
			log.Fatalln(err)
		}
		err = g.Board.Set(board.Right, coord, state)
		if err != nil {
			log.Fatalln(err)
		}
		if g.FireResp.Result != "Hit" {
			break
		}
	}
}

func (g *Game) StartGame() {
	for {
		g.checkGameStatus()
		if g.StatusResp.GameStatus != "game_in_progress" {
			time.Sleep(1 * time.Second)
			continue
		}
		resp, err := g.GameClient.GetGameDescription(g.Token)
		if err != nil {
			log.Fatalln(err)
		}
		g.GameDesc = resp
		break
	}
	for {
		g.checkGameStatus()
		g.endGameCheck()
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
	resp, err := g.GameClient.PostFire(g.Token, &r)
	if err != nil {
		log.Fatalln(err)
	}
	g.FireResp = resp
}

func (g *Game) checkGameStatus() {
	resp, err := g.GameClient.GetGameStatus(g.Token)
	if err != nil {
		log.Fatalln(err)
	}
	g.StatusResp = resp
}

func (g *Game) displayGameDescription() {
	writer := tabwriter.NewWriter(os.Stdout, 9, 8, 0, '\t', 0)
	fmt.Fprintf(writer, "%s  VS.  %s\n\n", g.GameDesc.Nick, g.GameDesc.Opponent)
	fmt.Fprintf(writer, "Gracz\tOpis\n")
	fmt.Fprintf(writer, "%s\t%s\n", "-----", "----")
	fmt.Fprintf(writer, "%s\t%s\n", g.GameDesc.Nick, g.GameDesc.Desc)
	fmt.Fprintf(writer, "%s\t%s\n\n", g.GameDesc.Opponent, g.GameDesc.OppDesc)
	writer.Flush()
}

func (g *Game) displayBoard() {
	resp, err := g.GameClient.GetGameBoards(g.Token)
	if err != nil {
		log.Fatalln(err)
	}
	g.BoardResp = resp

	err = g.Board.Import(g.BoardResp.Board)
	if err != nil {
		log.Fatalln(err)
	}
	g.Board.Display()
}
