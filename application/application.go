package application

import (
	"context"
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
	"main/client"
	"net/http"
	"time"
)

// Application represents the game application and its state.
type Application struct {
	Client              client.Client           // Client represents the client for making API requests.
	Token               string                  // Token stores the game token obtained from the API.
	StatusResp          *client.StatusResponse  // StatusResp stores the response from the status API endpoint.
	GameDesc            *client.StatusResponse  // GameDesc stores the response from the game description API endpoint.
	BoardResp           *client.BoardResponse   // BoardResp stores the response from the board API endpoint.
	FireResp            *client.FireResponse    // FireResp stores the response from the fire API endpoint.
	LobbiesList         *[]client.LobbyResponse // LobbiesList stores the list of available lobbies.
	Gui                 *gui.GUI                // Gui represents the graphical user interface.
	MyBoard             *gui.Board              // MyBoard represents the player's own game board.
	EnemyBoard          *gui.Board              // EnemyBoard represents the enemy's game board.
	lastTurnShotsAmount int                     // lastTurnShotsAmount stores the number of shots made in the last turn.
	MyBoardStates       [10][10]gui.State       // MyBoardStates stores the state of each cell in the player's own game board.
	EnemyBoardStates    [10][10]gui.State       // EnemyBoardStates stores the state of each cell in the enemy's game board.
	Ctx                 context.Context         // Ctx represents the application's context.
	currentShipsCoords  []string                // currentShipsCoords stores the current ship coordinates.
	allShots            int                     // allShots stores the total number of shots fired.
	hitShots            int                     // hitShots stores the number of shots that hit the enemy's ships.
	nick                string                  // nick stores the player's nickname.
	desc                string                  // desc stores the game description.
}

// NewGame creates a new game application.
// It takes a client instance as a parameter and initializes the application's fields.
// It sets up the GUI, game boards, and other initial values.
// The function returns the created Application instance.
func NewGame(c *client.Client) *Application {
	g, b1, b2 := guiConfig()

	states := [10][10]gui.State{}
	setBoardToEmpty(&states)

	shipsCoords := []string{
		"F1", "F2", "F3", "F4", "A1", "B1", "C1", "A3", "B3", "C3", "A5", "B5", "A7", "B7", "A9", "B9", "J9", "J7", "J5", "J3"}

	app := &Application{
		Client:              *c,
		Gui:                 g,
		MyBoard:             b1,
		EnemyBoard:          b2,
		lastTurnShotsAmount: 0,
		MyBoardStates:       states,
		EnemyBoardStates:    states,
		Ctx:                 context.Background(),
		currentShipsCoords:  shipsCoords,
		allShots:            0,
		hitShots:            0,
	}
	return app
}

// GenerateGameToken generates a game token by sending a request to the game endpoint.
// It takes the user's ship coordinates (coords), game description (desc), user's nick (nick),
// target player's nick (targetNick), and a flag indicating whether it's a game with a bot (wpbot) as parameters.
// The function constructs a StartGameRequest with the provided information and sends a POST request to the game endpoint.
// It updates the Application's Token field with the generated token if the request is successful.
// If there is an error during the request, the function logs the error and terminates the program.
func (g *Application) GenerateGameToken(coords []string, desc, nick, targetNick string, wpbot bool) {
	r := client.StartGameRequest{
		Coords:     coords,
		Desc:       desc,
		Nick:       nick,
		TargetNick: targetNick,
		Wpbot:      wpbot,
	}
	_, err := g.Client.SendRequest(http.MethodPost, client.GameEndpoint, "", &r, nil)
	if err != nil {
		log.Fatalln(err)
	}
	g.Token = g.Client.Token
}

// StartGame is the entry point for starting and managing the game.
// It performs the necessary steps to start and continue playing the game.
// The function does not return any value and runs in an infinite loop until the game is exited.
func (g *Application) StartGame() {
	deleteOldStatsFile()
	clearTerminal()
	g.nick, g.desc = NickAndDescChoice()
	g.displayMainMenu()
	turnChannel := make(chan bool)

	for true {
		g.initializeGame()

		go g.manageTurn(turnChannel)

		g.displayGameDescription()

		g.Gui.Start(g.Ctx, nil)
		turnChannel <- true

		clearTerminal()

		saveGameStats(g.allShots, g.hitShots)

		g.abandonGame()

		g.playAgainChoice()
	}
}

// abandonGame sends a request to abandon the game.
// It uses the `Client` field of the application to send a DELETE request to the abandonment endpoint.
// The function does not return any value but logs a fatal error if an error occurs during the request.
func (g *Application) abandonGame() {
	resp, err := g.Client.SendRequest(http.MethodDelete, client.AbandonEndpoint, g.Token, nil, &client.AbandonResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	_ = resp.(*client.AbandonResponse)
}

// playAgainChoice asks the user if they want to play again and handles the user's choice.
// It prompts the user to enter "y" or "n" and validates the input.
// If the user chooses to play again, it resets the board states and proceeds with a new game by calling `OpponentChoice`.
// If the user chooses not to play again, it goes back to the main menu by calling `displayMainMenu`.
// The function does not return any value but clears the terminal, updates the board states, and logs a fatal error if an error occurs during input validation.
func (g *Application) playAgainChoice() {
	fmt.Println("Do you want to play again? Enter (y or n): ")
	choice, err := getUserInputFromOptions([]string{"y", "n"})
	if err != nil {
		log.Fatalln()
	}

	clearTerminal()
	setBoardToEmpty(&g.EnemyBoardStates)
	g.EnemyBoard.SetStates(g.EnemyBoardStates)
	shipsOnlyBoard(&g.MyBoardStates)
	if choice == "y" {
		g.OpponentChoice()
	} else {
		g.displayMainMenu()
	}
}

// initializeGame initializes the game by setting up the player's board, checking the game status, and retrieving the game description from the server.
// It starts by calling initializeMyBoard to set up the player's board based on the server's response.
// Then, it enters a loop to continuously check the game status until the status changes from "waiting_wpbot" or "waiting".
// Once the game status changes to a different value, it sends a request to retrieve the game description and stores it in the GameDesc field.
func (g *Application) initializeGame() {
	g.initializeMyBoard()
	for {
		g.checkGameStatus()
		if g.StatusResp.GameStatus == "waiting_wpbot" {
			fmt.Println(g.StatusResp.GameStatus)
			time.Sleep(1 * time.Second)
			continue
		} else if g.StatusResp.GameStatus == "waiting" {
			fmt.Println(g.StatusResp.GameStatus)
			_, err := g.Client.SendRequest(http.MethodGet, client.RefreshEndpoint, g.Token, nil, nil)
			time.Sleep(2 * time.Second)
			if err != nil {
				log.Fatalln(err)
			}
			continue
		}
		resp, err := g.Client.SendRequest(http.MethodGet, client.DescEndpoint, g.Token, nil, &client.StatusResponse{})
		if err != nil {
			log.Fatalln(err)
			break
		}
		g.GameDesc = resp.(*client.StatusResponse)
		break
	}
}

// displayGameDescription displays the game description on the GUI.
// It creates a text element to show the players' nicks and draws it on the GUI.
// It then splits the game description and opponent's description into chunks of maximum length 45 characters.
// Each chunk is displayed as a separate text element on the GUI, starting from specific coordinates.
func (g *Application) displayGameDescription() {
	playersNicks := g.GameDesc.Nick + " vs. " + g.GameDesc.Opponent
	t := gui.NewText(5, 2, playersNicks, nil)
	g.Gui.Draw(t)

	chunks := splitString(g.GameDesc.Desc, 45)
	for i, chunk := range chunks {
		d := gui.NewText(1, 26+i, chunk, nil)
		g.Gui.Draw(d)
	}

	chunks2 := splitString(g.GameDesc.OppDesc, 45)
	for i, chunk := range chunks2 {
		d := gui.NewText(50, 26+i, chunk, nil)
		g.Gui.Draw(d)
	}
}

// initializeMyBoard initializes the user's game board based on the response received from the server.
// It sends a request to retrieve the user's board configuration and updates the application's MyBoardStates field accordingly.
func (g *Application) initializeMyBoard() {
	resp, err := g.Client.SendRequest(http.MethodGet, client.BoardEndpoint, g.Token, nil, &client.BoardResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	g.BoardResp = resp.(*client.BoardResponse)

	setBoardToEmpty(&g.MyBoardStates)

	for _, coord := range g.BoardResp.Board {
		c := StringToIntCoord(coord)
		g.MyBoardStates[c.X][c.Y] = gui.Ship
		g.MyBoard.SetStates(g.MyBoardStates)
	}
}
