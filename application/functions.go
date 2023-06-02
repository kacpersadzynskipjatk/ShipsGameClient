package application

import (
	"fmt"
	gui "github.com/grupawp/warships-gui/v2"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// splitString splits a string into chunks of a specified size.
// It takes a string and an integer as input, representing the input string and the desired chunk size, respectively.
// The function splits the input string into words and starts building chunks.
// It iterates over the words, adding them to the current chunk until the length of the chunk exceeds the chunk size.
// When the chunk size is exceeded, the current chunk is appended to the list of chunks, and a new chunk is started.
// The function trims any leading or trailing spaces from each chunk.
// Finally, it appends the last remaining chunk (if any) to the list of chunks and returns the resulting slice of strings.
func splitString(input string, chunkSize int) []string {
	var chunks []string
	words := strings.Fields(input)
	currentChunk := ""
	for _, word := range words {
		if len(currentChunk)+len(word)+1 > chunkSize {
			chunks = append(chunks, strings.TrimSpace(currentChunk))
			currentChunk = ""
		}
		currentChunk += " " + word
	}
	if currentChunk != "" {
		chunks = append(chunks, strings.TrimSpace(currentChunk))
	}
	return chunks
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

// contains checks if a string is present in a slice of strings.
// It takes a slice of strings and a string as input.
// The function iterates over the slice and compares each element with the given string.
// If a match is found, it returns true.
// If no match is found, it returns false.
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// clearTerminal clears the terminal screen by executing platform-specific commands.
func clearTerminal() {
	switch runtime.GOOS {
	case "windows":
		// For Windows, use "cmd" with the "/c" flag and "cls" command to clear the screen
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		// For Unix-like systems, use the "clear" command to clear the screen
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

// GetNickAndDesc prompts the user to enter their nickname and description,
// reads the input from the user, clears the terminal, and returns the entered
// nickname and description.
func GetNickAndDesc() (string, string) {
	fmt.Println("Enter nick:")
	nick, err := getUserInputWithLengthLimit(20)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Enter description:")
	desc, err := getUserInputWithLengthLimit(200)
	if err != nil {
		log.Fatalln(err)
	}

	clearTerminal()
	return nick, desc
}

// NickAndDescChoice prompts the user to choose whether they want to provide
// their nickname and description. If the user chooses 'y', it calls GetNickAndDesc
// to retrieve the nickname and description. Otherwise, it returns empty strings.
func NickAndDescChoice() (string, string) {
	fmt.Println("Do you want to enter a nick and description? Enter (y or n): ")
	choice, err := getUserInputFromOptions([]string{"y", "n"})
	if err != nil {
		log.Fatalln(err)
	}
	clearTerminal()

	if choice == "y" {
		return GetNickAndDesc()
	}
	return "", ""
}

// saveGameStats saves the game statistics (allShots and hitShots) to a file named "stats.txt".
// It creates the file (overwriting if it exists) and writes the statistics to it.
func saveGameStats(allShots, hitShots int) {
	// Create the file (overwrite if it exists)
	file, err := os.Create("stats.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Write stats to the file
	_, err = fmt.Fprintf(file, "%d\n%d", allShots, hitShots)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// deleteOldStatsFile deletes the file "stats.txt" if it exists.
// It handles different scenarios such as file not found or error accessing the file.
func deleteOldStatsFile() {
	filename := "stats.txt"

	// Check if the file exists
	_, err := os.Stat(filename)
	if err != nil {
		// If the file doesn't exist, no need to delete it
		if os.IsNotExist(err) {
			fmt.Printf("File %s does not exist.\n", filename)
			return
		}

		// If there was an error accessing the file, handle it accordingly
		fmt.Printf("Error accessing file %s: %s\n", filename, err.Error())
		return
	}

	// The file exists, so proceed with deletion
	err = os.Remove(filename)
	if err != nil {
		fmt.Printf("Error deleting file %s: %s\n", filename, err.Error())
		return
	}
}
