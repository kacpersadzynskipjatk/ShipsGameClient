package application

import (
	"fmt"
	"log"
	"main/client"
	"net/http"
	"os"
	"sort"
)

// displayMainMenu displays the main menu options to the user and handles their choices.
// It enters a loop to continue displaying the menu until an option is selected.
// The function performs the corresponding actions based on the selected menu choice.
func (g *Application) displayMainMenu() {
	loop := true
	for loop {
		fmt.Println("Menu, enter the number to choose an option:")
		fmt.Println("1. Play game")
		fmt.Println("2. Display Waiting Players")
		fmt.Println("3. Display statistics")
		fmt.Println("4. Display player ranking")
		fmt.Println("5. How to play")
		fmt.Println("Press 'Ctrl+C' to exit")
		menuChoice, err := getUserInputFromOptions([]string{"1", "2", "3", "4", "5"})
		if err != nil {
			log.Fatalln(err)
		}
		clearTerminal()
		switch menuChoice {
		case "1":
			g.OpponentChoice()
			loop = false
		case "2":
			g.DisplayOpponentsList()
		case "3":
			displayStatistics()
		case "4":
			g.displayPlayerRanking()
		case "5":
			DisplayInstructions()
		}
	}
}

// OpponentChoice allows the user to choose the opponent for the game.
// It prompts the user to select between playing with a bot or with another player.
// It modifies the currentShipsCoords if the ShipPlacementChoice function returns valid ship coordinates.
// The function displays the available options and validates the user's input.
// After the user makes a choice, it generates a game token based on the selected opponent
// and the user's own information (nick, desc, and targetNick if applicable).
func (g *Application) OpponentChoice() {
	coords := ShipPlacementChoice(&g.MyBoardStates)
	if coords != nil {
		g.currentShipsCoords = coords
	}
	fmt.Println("1. Play with bot")
	fmt.Println("2. Play with player")

	// Prompt the user to choose an opponent
	enemyChoice, err := getUserInputFromOptions([]string{"1", "2"})
	if err != nil {
		log.Fatalln(err)
	}
	clearTerminal()

	// Handle the user's choice
	switch enemyChoice {
	case "1":
		// Generate a game token for playing with a bot
		g.GenerateGameToken(g.currentShipsCoords, g.desc, g.nick, "", true)
	case "2":
		g.DisplayOpponentsList()
		fmt.Println("Type enemy nick or press ENTER to wait")

		// Prompt the user to enter the target player's nick
		var targetNick string
		fmt.Scanln(&targetNick)
		clearTerminal()

		// Generate a game token for playing with another player
		g.GenerateGameToken(g.currentShipsCoords, g.desc, g.nick, targetNick, false)
	}
}

// DisplayOpponentsList retrieves and displays the list of available opponents.
// It sends a request to the lobby API endpoint and stores the response in the LobbiesList field.
// The function then prints the list of opponents' nicknames and game statuses.
func (g *Application) DisplayOpponentsList() {
	resp, err := g.Client.SendRequest(http.MethodGet, client.LobbyEndpoint, g.Token, nil, &[]client.LobbyResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	g.LobbiesList = resp.(*[]client.LobbyResponse)
	fmt.Println("Waiting Opponents:")
	for i, lobby := range *g.LobbiesList {
		fmt.Printf("%d. nick: %s, status: %s\n", i+1, lobby.Nick, lobby.GameStatus)
	}
	fmt.Println()
}

// checkGameStatus retrieves and updates the game status in the Application struct.
// It sends a request to retrieve the current game status from the server and updates the application's StatusResp field accordingly.
func (g *Application) checkGameStatus() {
	resp, err := g.Client.SendRequest(http.MethodGet, client.GameEndpoint, g.Token, nil, &client.StatusResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	g.StatusResp = resp.(*client.StatusResponse)
}

// displayStatistics displays the game statistics based on the contents of the "stats.txt" file.
// If the file does not exist, it creates a new file and initializes the statistics with 0 values.
// The function reads the statistics from the file and calculates the hit ratio.
func displayStatistics() {
	if _, err := os.Stat("stats.txt"); os.IsNotExist(err) {
		saveGameStats(0, 0)
	}

	file, err := os.Open("stats.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Print("My epic stats:\n\n")

	var allShots, hitShots int
	_, err = fmt.Fscanf(file, "%d\n%d", &allShots, &hitShots)
	if err != nil {
		fmt.Println(err)
		return
	}

	var ratio float64
	if allShots == 0 {
		ratio = 0
	} else {
		ratio = float64(hitShots) / float64(allShots)
	}

	fmt.Printf("All my shots: %d\n", allShots)
	fmt.Printf("On target: %d\n", hitShots)
	fmt.Printf("Hit Ratio: %f\n", ratio)

	if ratio > 0.39 {
		fmt.Printf("You are a true Capitan!\n\n")
	} else if ratio > 0.19 {
		fmt.Printf("You are an average player, get some sailing courses.\n\n")
	} else {
		fmt.Printf("You are a weak player, train more!\n\n")
	}
}

// displayPlayerRanking retrieves and displays the player ranking and statistics.
// It sends requests to the server to get the player's individual stats and the top players' stats.
// The player's stats are added to the top players' stats and then sorted by points in descending order.
// Finally, it prints the ranking information, including the player's rank, nickname, games played, points, and wins.
func (g *Application) displayPlayerRanking() {
	resp, err := g.Client.SendRequest(http.MethodGet, client.PlayerStatsEndpoint+"/"+g.nick, g.Token, nil, &client.PlayerStatsResponse{})
	if err != nil {
		log.Println(err)
	}
	playerStats := resp.(*client.PlayerStatsResponse)

	resp, err = g.Client.SendRequest(http.MethodGet, client.PlayerStatsEndpoint, g.Token, nil, &client.TopPlayersStatsResponse{})
	if err != nil {
		log.Fatalln(err)
	}
	topPlayers := resp.(*client.TopPlayersStatsResponse)
	topPlayers.Stats = append(topPlayers.Stats, playerStats.Stats)

	// Sort the stats array by points in descending order
	sort.Slice(topPlayers.Stats, func(i, j int) bool {
		return topPlayers.Stats[i].Points > topPlayers.Stats[j].Points
	})

	for _, player := range topPlayers.Stats {
		fmt.Printf("Rank: %d\n", player.Rank)
		fmt.Printf("Nick: %s\n", player.Nick)
		fmt.Printf("Games played: %d\n", player.Games)
		fmt.Printf("Points: %d\n", player.Points)
		fmt.Printf("Wins: %d\n\n", player.Wins)
	}
}
