package application

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
)

// ShipPlacementChoice allows the player to choose whether he wants to place his ships himself or not.
// If no default placement is selected.
// The user clicks on the board to place the ships.
func ShipPlacementChoice(myBoardStates *[10][10]gui.State) []string {
	fmt.Println("Do you want to place your ships? Enter (y or n): ")
	choice, err := getUserInputFromOptions([]string{"y", "n"})
	if err != nil {
		log.Fatalln(err)
	}

	clearTerminal()
	if choice == "y" {
		g := gui.NewGUI(true)
		b := gui.NewBoard(1, 3, nil)
		text1 := gui.NewText(1, 1, "Click on board to set ships:", nil)
		g.Draw(text1)
		g.Draw(b)
		go SetShips(g, b, myBoardStates)
		g.Start(context.Background(), nil)

		var coords []string
		for i := range myBoardStates {
			for j := range myBoardStates[i] {
				c := Coord{i, j}
				if myBoardStates[i][j] == gui.Ship {
					res := IntCoordToString(c)
					coords = append(coords, res)
				}
			}
		}
		return coords
	}
	return nil
}

// SetShips allows the user to set ships on the board by listening for coordinate inputs from a given board.
// It takes a pointer to a gui.Board and a pointer to a 10x10 array of gui.State as input.
// The function initializes the board states to Empty and starts a loop until all ships are placed.
// In each iteration, it listens for a coordinate input from the board.
// If the corresponding board state is Empty, it sets it to Ship and decreases the ship counter.
// If the corresponding board state is already a Ship, it resets it to Empty and increases the ship counter.
// The function updates the board states and continues until all ships are placed.
func SetShips(g *gui.GUI, b *gui.Board, boardStates *[10][10]gui.State) {
	setBoardToEmpty(boardStates)
	counter := 20
	loop := true
	for loop {
		coord := b.Listen(context.Background())
		c := StringToIntCoord(coord)
		switch boardStates[c.X][c.Y] {
		case gui.Empty:
			boardStates[c.X][c.Y] = gui.Ship
			counter--
		case gui.Ship:
			boardStates[c.X][c.Y] = gui.Empty
			counter++
		}
		if counter == 0 {
			visitedCoords := make([]Coord, 0)
			leftShips := map[int]int{
				1: 4,
				2: 3,
				3: 2,
				4: 1,
			}
			for i := 0; i < 10; i++ {
				for j := 0; j < 10; j++ {
					x := Coord{i, j}
					if boardStates[i][j] == gui.Ship && !x.containedIn(visitedCoords) {
						shipSize := 0
						DetectShip(boardStates, x, &visitedCoords, &shipSize, gui.Ship, false)
						if shipSize >= 1 && shipSize <= 4 {
							leftShips[shipSize]--
						}
					}
				}
			}
			var flag = true
			for i := 1; i <= 4; i++ {
				if leftShips[i] != 0 {
					flag = false
					text1 := gui.NewText(1, 28, "Wrong ships placement, correct it!", nil)
					g.Draw(text1)
					break
				}
			}
			if flag == true {
				text1 := gui.NewText(1, 28, "Your ships are correct, press Ctrl + C", nil)
				g.Draw(text1)
				loop = false
			}
		}
		b.SetStates(*boardStates)
	}
}

// setBoardToEmpty sets all states of a given board to Empty.
// It takes a pointer to a 10x10 array of gui.State as input.
// The function iterates over each element of the board and assigns it the value of gui.Empty.
func setBoardToEmpty(states *[10][10]gui.State) {
	for i := range states {
		for j := range states[i] {
			states[i][j] = gui.Empty
		}
	}
}

// shipsOnlyBoard sets all states of a given board to Empty, except for the cells that contain a Ship.
// It takes a pointer to a 10x10 array of gui.State as input.
// The function iterates over each element of the board and checks if it is not a Ship.
// If the state is not a Ship, it is assigned the value of gui.Empty.
func shipsOnlyBoard(states *[10][10]gui.State) {
	for i := range states {
		for j := range states[i] {
			if states[i][j] != gui.Ship {
				states[i][j] = gui.Empty
			}
		}
	}
}

// guiConfig configures and returns a GUI, along with two boards.
// It creates a new GUI with the "true" parameter indicating ANSI color support.
// It also creates two boards using the provided configuration.
// The configuration can be modified by changing the values of the `cfg` variable.
// The function then draws the boards on the GUI and returns the GUI, the first board, and the second board.
func guiConfig() (*gui.GUI, *gui.Board, *gui.Board) {

	// Change cfg values to configure the board
	cfg := gui.NewBoardConfig()

	b1 := gui.NewBoard(1, 6, cfg)
	b2 := gui.NewBoard(50, 6, cfg)
	g := gui.NewGUI(true)
	g.Draw(b1)
	g.Draw(b2)
	return g, b1, b2
}
