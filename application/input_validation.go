package application

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// getUserInputFromOptions prompts the user to enter their choice from a list of options
// and returns the valid user input.
// It continuously prompts the user until a valid option is provided.
// The function expects a string slice of valid options.
func getUserInputFromOptions(options []string) (string, error) {

	// Create a map of valid options for quick lookup
	validOptions := make(map[string]bool)
	for _, option := range options {
		validOptions[option] = true
	}

	for {
		fmt.Printf("Enter your choice (%s): ", strings.Join(options, "/"))

		// Read user input
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading input: %v", err)
		}

		// Remove the trailing newline character
		input = strings.TrimSpace(input)

		// Check if the input is a valid option
		if _, ok := validOptions[input]; !ok {
			// If the input is not valid, prompt the user again
			fmt.Printf("Invalid choice. Please enter one of %s.\n", strings.Join(options, "/"))
			continue
		}

		// Return the valid input
		return input, nil
	}
}

// getUserInputWithLengthLimit prompts the user to enter input and returns the valid user input
// with a maximum length limit and restricted to the specified valid characters.
// It continuously prompts the user until a valid input within the length and character limit is provided.
// The function expects the maximum length of the user input.
// It returns the user input as a string and any error that occurred during input reading.
func getUserInputWithLengthLimit(maxLength int) (string, error) {
	validChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789._-"

	for {
		// Prompt the user to enter their input
		fmt.Printf("Enter your input (maximum %d characters): ", maxLength)

		// Read user input from standard input
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("error reading input: %v", err)
		}

		// Remove the trailing newline character and trim whitespace
		input = strings.TrimSpace(input)

		// Validate the input length
		if len(input) > maxLength {
			fmt.Printf("Invalid input. Maximum input length is %d characters.\n", maxLength)
			continue
		}

		// Validate the input characters
		for _, ch := range input {
			if !strings.ContainsRune(validChars, ch) {
				input = ""
				fmt.Println("Invalid input. Input contains invalid characters.")
				break
			}
		}

		// If the input is valid, return it
		if input != "" {
			return input, nil
		}
	}
}
