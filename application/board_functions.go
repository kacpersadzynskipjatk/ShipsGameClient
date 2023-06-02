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
		go SetShips(b, myBoardStates)
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
func SetShips(b *gui.Board, boardStates *[10][10]gui.State) {
	setBoardToEmpty(boardStates)
	counter := 20
	for counter > 0 {
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

	b1 := gui.NewBoard(1, 4, cfg)
	b2 := gui.NewBoard(50, 4, cfg)
	g := gui.NewGUI(true)
	g.Draw(b1)
	g.Draw(b2)
	return g, b1, b2
}