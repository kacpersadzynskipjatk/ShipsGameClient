package application

import (
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
	"main/client"
	"net/http"
	"strconv"
	"time"
)

// manageTurn manages a player's turn in the game.
// It displays text on the GUI to indicate whether it's the player's turn or the opponent's turn.
// It continuously checks the game status and waits until the player is allowed to make a shot.
// Once the player is allowed to make a shot, it starts a timer and displays it on the GUI.
// It then sets the opponent's shots on the player's board and updates the GUI.
// After that, it displays text to indicate that it's the player's turn and calls the makeShotAndCheckEndGame function.
// The function repeats the above steps until the game ends.
func (g *Application) manageTurn(channel chan bool) {
	enemyTurnText := gui.NewText(5, 0, "Enemy turn wait! 😒", nil)
	myTurnText := gui.NewText(5, 0, "Your turn shoot! 😎", nil)
	loop := true
	ch := make(chan bool, 1)
	for loop {
		select {
		case <-channel:
			ch <- true
			return
		default:
			g.Gui.Draw(enemyTurnText)
			g.checkGameStatus()
			if !g.StatusResp.ShouldFire {
				time.Sleep(1 * time.Second)
				continue
			}
			go g.displayTimer(ch)
			g.setOpponentShots()
			time.Sleep(200 * time.Millisecond)
			g.Gui.Draw(myTurnText)
			loop = !g.makeShotAndCheckEndGame()
			ch <- true
		}
	}
}

// displayTimer displays and updates the timer on the GUI during the player's turn.
// It receives a channel as a parameter to listen for signals indicating the end of the turn.
// The function starts with an initial timeLeft value of 60 seconds.
// It creates a text element on the GUI to display the remaining time.
// The function continuously updates the timeLeft value by checking the game status every 5 seconds.
// It removes the previous text element from the GUI and creates a new one with the updated timeLeft.
// It sleeps for 1 second and decrements the timeLeft value.
// The function repeats these steps until it receives a signal on the provided channel or the timeLeft value reaches 0.
func (g *Application) displayTimer(ch chan bool) {
	timeLeft := 60
	var text = gui.NewText(0, 0, "", nil)
	i := 0
	for {
		if i%5 == 0 {
			g.checkGameStatus()
			timeLeft = g.StatusResp.Timer
		}
		g.Gui.Remove(text)
		select {
		case <-ch:
			return
		default:
			text = gui.NewText(28, 0, "Time left:"+strconv.Itoa(timeLeft), nil)
			g.Gui.Draw(text)
			time.Sleep(1 * time.Second)
			timeLeft--
		}
		i++
		if timeLeft == 0 {
			fmt.Println(" Time ran out, press Ctrl + C")
			break
		}
	}
}

// setOpponentShots updates the player's game board with the shots made by the opponent.
// It retrieves the opponent's shots from the StatusResp field of the Application instance.
// For each shot, it checks whether it is a hit or a miss by comparing it with the player's board configuration.
// It updates the corresponding state in the MyBoardStates field of the Application instance.
// Finally, it updates the GUI to reflect the changes by setting the states of the player's game board.
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

// makeShotAndCheckEndGame handles the player's turn to make a shot and checks if the game has ended.
// It listens for the coordinates of the shot from the enemy board.
// It performs the shot and updates the game board states accordingly.
// It determines the result of the shot (hit, miss, or sunk) and updates the GUI to reflect the changes.
// It also checks the game status and performs additional actions if the game has ended.
// The function returns a boolean indicating whether the game has ended.
func (g *Application) makeShotAndCheckEndGame() bool {
	loop := true
	endGame := false
	for loop {
		coord := g.EnemyBoard.Listen(g.Ctx)
		g.Fire(coord)
		g.allShots++
		c := StringToIntCoord(coord)

		// Shot the same ship more than one time
		if g.EnemyBoardStates[c.X][c.Y] == gui.Hit {
			return false
		}

		switch g.FireResp.Result {
		case "hit":
			g.hitShots++
			g.EnemyBoardStates[c.X][c.Y] = gui.Hit
			g.EnemyBoard.SetStates(g.EnemyBoardStates)
		case "miss":
			g.EnemyBoardStates[c.X][c.Y] = gui.Miss
			g.EnemyBoard.SetStates(g.EnemyBoardStates)
			loop = false
		case "sunk":
			g.hitShots++
			g.EnemyBoardStates[c.X][c.Y] = gui.Hit
			visitedCoords := make([]Coord, 0)
			shipSize := 0
			DetectShip(&g.EnemyBoardStates, c, &visitedCoords, &shipSize, gui.Hit, true)
			g.enemyLeftShips[shipSize]--
			g.displayEnemyLeftShips()
			g.EnemyBoard.SetStates(g.EnemyBoardStates)
			g.checkGameStatus()
			endGame = g.endGameCheck()
			loop = !endGame
		}
	}
	return endGame
}

func (g *Application) displayEnemyLeftShips() {
	x := 96
	ships := gui.NewText(x, 15, "Enemy ships left:\n", nil)
	g.Gui.Draw(ships)
	text := fmt.Sprintf("Size 1 - %d\n", g.enemyLeftShips[1])
	ships = gui.NewText(x, 16, text, nil)
	g.Gui.Draw(ships)
	text = fmt.Sprintf("Size 2 - %d\n", g.enemyLeftShips[2])
	ships = gui.NewText(x, 17, text, nil)
	g.Gui.Draw(ships)
	text = fmt.Sprintf("Size 3 - %d\n", g.enemyLeftShips[3])
	ships = gui.NewText(x, 18, text, nil)
	g.Gui.Draw(ships)
	text = fmt.Sprintf("Size 4 - %d\n", g.enemyLeftShips[4])
	ships = gui.NewText(x, 19, text, nil)
	g.Gui.Draw(ships)
}

// Fire sends a fire request to the server to execute a shot at the specified coordinate.
// It takes the coordinate as a parameter and sends the request to the server.
// It stores the response in the `FireResp` field of the application.
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

func DetectShip(boardToUpdate *[10][10]gui.State, c Coord, visitedCoords *[]Coord, shipSize *int, countedField gui.State, makeBorder bool) {
	*shipSize++
	*visitedCoords = append(*visitedCoords, c)
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
			if makeBorder && boardToUpdate[dx][dy] == gui.Empty {
				boardToUpdate[dx][dy] = gui.Miss
			} else if boardToUpdate[dx][dy] == countedField {
				DetectShip(boardToUpdate, nextCoord, visitedCoords, shipSize, countedField, makeBorder)
			}
		}
	}
}

// endGameCheck checks if the game has ended and displays the appropriate message on the GUI.
// It examines the `GameStatus` and `LastGameStatus` fields of the `StatusResponse` to determine the game outcome.
// It draws the corresponding text on the GUI using the `Gui` field of the application.
// The function returns true if the game has ended, and false otherwise.
func (g *Application) endGameCheck() bool {
	if g.StatusResp.GameStatus == "ended" {
		text1 := gui.NewText(5, 0, "The game has ended", nil)
		var text2 *gui.Text
		g.Gui.Draw(text1)
		if g.StatusResp.LastGameStatus == "win" {
			s := fmt.Sprintf("You won %s congratulations!!! ", g.StatusResp.Nick)
			text2 = gui.NewText(25, 0, s, nil)
		} else if g.StatusResp.LastGameStatus == "lose" {
			s := fmt.Sprintf("You lost, player %s won", g.StatusResp.Opponent)
			text2 = gui.NewText(25, 0, s, nil)
		}
		g.Gui.Draw(text2)
		return true
	}
	return false
}
