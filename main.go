package main

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"main/application"
	"main/client"
	"net/http"
	"time"
)

func getNickAndDesc() (string, string) {
	fmt.Println("Podaj nick:")
	var nick string
	fmt.Scanln(&nick)
	fmt.Println("Podaj opis:")
	var desc string
	fmt.Scanln(&desc)
	return nick, desc
}
func nickAndDescChoice() (string, string) {
	fmt.Println("Czy chcesz wprowadzić nick i opis gracza? wpisz: y/n")
	var choice string
	fmt.Scanln(&choice)
	if choice == "y" {
		n, d := getNickAndDesc()
		return n, d
	}
	return "", ""
}

func shipPlacementChoice() []string {
	fmt.Println("Czy chcesz rozmieścić swoje statki? wpisz: y/n")
	var choice string
	fmt.Scanln(&choice)
	if choice == "y" {
		myBoardStates := [10][10]gui.State{}
		g := gui.NewGUI(true)
		b := gui.NewBoard(1, 1, nil)
		g.Draw(b)
		go setShips(b, &myBoardStates)
		g.Start(context.Background(), nil)
		var coords []string
		for i := range myBoardStates {
			for j := range myBoardStates[i] {
				if myBoardStates[i][j] == gui.Ship {
					var res string
					charX := rune(i)
					charY := rune(j)
					charX += 'A'
					charY += '1'
					res = string(charX)
					if j == 9 {
						res += "10"
					} else {
						res += string(charY)
					}
					coords = append(coords, res)
				}
			}
		}
		return coords
	}
	return nil
}

func setShips(b *gui.Board, boardStates *[10][10]gui.State) {
	for i := range boardStates {
		for j := range boardStates[i] {
			boardStates[i][j] = gui.Empty
		}
	}
	counter := 20
	for counter > 0 {
		coord := b.Listen(context.Background())
		c := application.StringToIntCoord(coord)
		switch boardStates[c.X][c.Y] {
		case gui.Empty:
			boardStates[c.X][c.Y] = gui.Ship
			counter--
		case gui.Ship:
			boardStates[c.X][c.Y] = gui.Empty
			counter++
		}
		b.SetStates(*boardStates)
	}
}

func main() {
	c := http.Client{Timeout: 15 * time.Second}
	gc := client.NewGameClient(&c)
	newGame := application.NewGame(gc)
	displayMainMenu(newGame)
}

func displayMainMenu(app *application.Application) {
	defaultCoords := []string{
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
	for true {
		fmt.Println("Menu, enter number to choose option:")
		fmt.Println("1. Play game")
		fmt.Println("2. Display Waiting Players")
		fmt.Println("3. Display statistics")
		var menuChoice string
		fmt.Scanln(&menuChoice)
		switch menuChoice {
		case "1":
			nick, desc := nickAndDescChoice()
			coords := shipPlacementChoice()
			if coords != nil {
				defaultCoords = coords
			}
			fmt.Println("1. Play with bot")
			fmt.Println("2. Play with player")
			var enemyChoice string
			fmt.Scanln(&enemyChoice)
			switch enemyChoice {
			case "1":
				app.GenerateGameToken(defaultCoords, desc, nick, "", true)
			case "2":
				app.DisplayOpponentsList()
				targetNick := ""
				fmt.Println("Wprowadź Nick oponenta lub pomiń klikając enter:")
				fmt.Scanln(&targetNick)
				app.GenerateGameToken(defaultCoords, desc, nick, targetNick, false)
			}
			app.StartGame()
			break
		case "2":
			app.DisplayOpponentsList()
			continue
		case "3":
			fmt.Println("Not done yet")
			continue
		}
	}
}
