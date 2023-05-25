package application

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
	"main/client"
	"net/http"
	"os"
	"strconv"
	"time"
)

func guiConfig() (*gui.GUI, *gui.Board, *gui.Board) {
	cfg := gui.NewBoardConfig()
	//change cfg values to configure the board
	b1 := gui.NewBoard(1, 4, cfg)
	b2 := gui.NewBoard(50, 4, cfg)
	g := gui.NewGUI(true)
	g.Draw(b1)
	g.Draw(b2)
	return g, b1, b2
}

type Application struct {
	Client              client.Client
	Token               string
	StatusResp          *client.StatusResponse
	GameDesc            *client.StatusResponse
	BoardResp           *client.BoardResponse
	FireResp            *client.FireResponse
	OpponentsList       *[]client.OpponentResponse
	Gui                 *gui.GUI
	MyBoard             *gui.Board
	EnemyBoard          *gui.Board
	lastTurnShotsAmount int
	MyBoardStates       [10][10]gui.State
	EnemyBoardStates    [10][10]gui.State
	Ctx                 context.Context
}

type Coord struct {
	X int
	Y int
}

func (a *Coord) containedIn(sl []Coord) bool {
	for _, s := range sl {
		if s.X == a.X && s.Y == a.Y {
			return true
		}
	}
	return false
}

func NewGame(c *client.Client) *Application {
	g, b1, b2 := guiConfig()
	states := [10][10]gui.State{}
	for i := range states {
		for j := range states[i] {
			states[i][j] = gui.Empty
		}
	}
	app := &Application{
		Client:              *c,
		Gui:                 g,
		MyBoard:             b1,
		EnemyBoard:          b2,
		lastTurnShotsAmount: 0,
		MyBoardStates:       states,
		EnemyBoardStates:    states,
		Ctx:                 context.Background(),
	}
	return app
}

func (g *Application) GenerateGameToken(coords []string, desc, nick, targetNick string, wpbot bool) {
	r := client.StartGameRequest{
		Coords:     coords,
		Desc:       desc,
		Nick:       nick,
		TargetNick: targetNick,
		Wpbot:      wpbot,
	}
	_, err := g.Client.SendRequest(http.MethodPost, client.GameEndpoint, "", r, &[]client.StartGameResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	g.Token = g.Client.Token
}

func (g *Application) endGameCheck() {
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
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (g *Application) setOpponentShots() {
	oppShots := g.StatusResp.OppShots
	currentAmount := len(oppShots)

	var state gui.State
	for i := g.lastTurnShotsAmount; i < currentAmount; i++ {
		if contains(g.BoardResp.Board, oppShots[i]) {
			state = gui.Hit
		} else {
			state = gui.Miss
		}
		c := StringToIntCoord(oppShots[i])
		g.MyBoardStates[c.X][c.Y] = state
	}
	g.MyBoard.SetStates(g.MyBoardStates)
	g.lastTurnShotsAmount = currentAmount
}

func (g *Application) makeShot() {
	loop := true
	for loop {
		//fmt.Println("Type 'Abandon' to end the game")
		//fmt.Println("Twoja Tura!\nWpisz koordynaty by strzelić:")
		coord := g.EnemyBoard.Listen(g.Ctx)
		g.Fire(coord)
		c := StringToIntCoord(coord)
		switch g.FireResp.Result {
		case "hit":
			g.EnemyBoardStates[c.X][c.Y] = gui.Hit
			g.EnemyBoard.SetStates(g.EnemyBoardStates)
		case "miss":
			g.EnemyBoardStates[c.X][c.Y] = gui.Miss
			g.EnemyBoard.SetStates(g.EnemyBoardStates)
			loop = false
		case "sunk":
			g.EnemyBoardStates[c.X][c.Y] = gui.Hit
			visitedCoords := make([]Coord, 0)
			g.CreateShipBorder(c, &visitedCoords)
			g.EnemyBoard.SetStates(g.EnemyBoardStates)
			g.checkGameStatus()
			g.endGameCheck()
		}
	}
}

func (g *Application) CreateShipBorder(c Coord, visitedCoords *[]Coord) {
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			dx := x + c.X
			dy := y + c.Y
			if dx < 0 || dx > 9 {
				continue
			}
			if dy < 0 || dy > 9 {
				continue
			}
			nextCoord := Coord{dx, dy}
			if nextCoord.containedIn(*visitedCoords) {
				continue
			}
			*visitedCoords = append(*visitedCoords, c)
			if g.EnemyBoardStates[dx][dy] == gui.Empty {
				g.EnemyBoardStates[dx][dy] = gui.Miss
			} else if g.EnemyBoardStates[dx][dy] == gui.Hit {
				g.CreateShipBorder(nextCoord, visitedCoords)
			}
		}
	}
}

func (g *Application) DisplayOpponentsList() {
	resp, err := g.Client.SendRequest(http.MethodGet, client.ListEndpoint, g.Token, nil, &[]client.OpponentResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	g.OpponentsList = resp.(*[]client.OpponentResponse)
	fmt.Println("Waiting Opponents:")
	for i, opponent := range *g.OpponentsList {
		fmt.Printf("%d. nick: %s, status: %s\n", i+1, opponent.Nick, opponent.GameStatus)
	}
	fmt.Println()
}

func (g *Application) StartGame() {
	g.initializeMyBoard()
	for {
		g.checkGameStatus()
		if g.StatusResp.GameStatus == "waiting_wpbot" {
			fmt.Println(g.StatusResp.GameStatus)
			time.Sleep(1 * time.Second)
			continue
		} else if g.StatusResp.GameStatus == "waiting" {
			fmt.Println(g.StatusResp.GameStatus)
			_, err := g.Client.SendRequest(http.MethodGet, client.RefreshEndpoint, g.Token, nil, &client.AbandonResponse{})
			time.Sleep(2 * time.Second)
			if err != nil {
				log.Fatalln(err)
			}
			continue
		}
		resp, err := g.Client.SendRequest(http.MethodGet, client.DescEndpoint, g.Token, nil, &client.StatusResponse{})
		if err != nil {
			log.Fatalln(err)
		}
		g.GameDesc = resp.(*client.StatusResponse)
		break
	}
	go g.manageTurn()
	g.displayGameDescription()

	g.Gui.Start(g.Ctx, nil)
}

func (g *Application) Fire(coord string) {
	r := client.FireRequest{
		Coord: coord,
	}
	resp, err := g.Client.SendRequest(http.MethodPost, client.FireEndpoint, g.Token, &r, &client.FireResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	g.FireResp = resp.(*client.FireResponse)
}

func (g *Application) checkGameStatus() {
	resp, err := g.Client.SendRequest(http.MethodGet, client.GameEndpoint, g.Token, nil, &client.StatusResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	g.StatusResp = resp.(*client.StatusResponse)
}

func (g *Application) displayGameDescription() {
	playersNicks := g.GameDesc.Nick + " vs. " + g.GameDesc.Opponent
	t := gui.NewText(1, 2, playersNicks, nil)
	g.Gui.Draw(t)
	//writer := tabwriter.NewWriter(os.Stdout, 9, 8, 0, '\t', 0)
	//fmt.Fprintf(writer, "%s  VS.  %s\n\n", g.GameDesc.Nick, g.GameDesc.Opponent)
	//fmt.Fprintf(writer, "Gracz\tOpis\n")
	//fmt.Fprintf(writer, "%s\t%s\n", "-----", "----")
	//fmt.Fprintf(writer, "%s\t%s\n", g.GameDesc.Nick, g.GameDesc.Desc)
	//fmt.Fprintf(writer, "%s\t%s\n\n", g.GameDesc.Opponent, g.GameDesc.OppDesc)
	//writer.Flush()
}

func (g *Application) initializeMyBoard() {
	resp, err := g.Client.SendRequest(http.MethodGet, client.BoardEndpoint, g.Token, nil, &client.BoardResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	g.BoardResp = resp.(*client.BoardResponse)
	for _, coord := range g.BoardResp.Board {
		c := StringToIntCoord(coord)
		g.MyBoardStates[c.X][c.Y] = gui.Ship
		g.MyBoard.SetStates(g.MyBoardStates)
	}
}

func StringToIntCoord(coord string) Coord {
	chars := []rune(coord)
	var sb = ""
	for i, char := range chars {
		if i != 0 {
			sb += string(char)
		}
	}
	y, err := strconv.Atoi(sb)
	if err != nil {
		log.Fatalln(err)
	}
	y -= 1
	x := int(chars[0] - 'A')
	c := Coord{x, y}
	return c
}

func (g *Application) manageTurn() {
	for {
		g.checkGameStatus()
		g.endGameCheck()
		if !g.StatusResp.ShouldFire {
			time.Sleep(1 * time.Second)
			continue
		}
		g.setOpponentShots()
		g.makeShot()
	}
}
