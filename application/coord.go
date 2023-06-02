package application

import (
	"log"
	"strconv"
)

//Coord representation of the x, y coordinates
type Coord struct {
	X int
	Y int
}

// containedIn checks if a Coord is contained in a slice of Coord.
// It takes a pointer to a Coord and a slice of Coord as input.
// The function iterates over the slice and compares the X and Y values of each Coord with the given Coord.
// If a matching Coord is found, it returns true.
// If no matching Coord is found, it returns false.
func (a *Coord) containedIn(sl []Coord) bool {
	for _, s := range sl {
		if s.X == a.X && s.Y == a.Y {
			return true
		}
	}
	return false
}

// StringToIntCoord converts a string coordinate to a Coord struct.
// It takes a string coordinate as input and returns a Coord struct.
// The function extracts the numeric part from the string coordinate and converts it to an integer.
// The letter part of the string coordinate is converted to the corresponding integer representation.
// The resulting values are used to create a Coord struct, which is then returned.
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

// IntCoordToString converts a Coord struct to a string coordinate.
// It takes a Coord struct as input and returns a string coordinate.
// The X and Y values of the Coord struct are converted to their corresponding character representations.
// The character X value is incremented by 'A' to convert it to the corresponding letter.
// The character Y value is incremented by '1' to convert it to the corresponding numeric representation.
// The resulting characters are combined to form the string coordinate, which is then returned.
func IntCoordToString(coord Coord) string {
	var result string
	charX := rune(coord.X)
	charY := rune(coord.Y)
	charX += 'A'
	charY += '1'
	result = string(charX)
	if coord.Y == 9 {
		result += "10"
	} else {
		result += string(charY)
	}
	return result
}
