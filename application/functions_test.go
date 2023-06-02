package application

import (
	"testing"
)

func TestStringToIntCoord(t *testing.T) {
	t.Run("Middle coord", func(t *testing.T) {
		coord := "E5"
		expected := Coord{4, 4}
		if result := StringToIntCoord(coord); result.X != expected.X || result.Y != expected.Y {
			t.Errorf("Expected %d %d but got %d %d", expected.X, expected.Y, result.X, result.Y)
		}
	})
	t.Run("Left edge coord", func(t *testing.T) {
		coord := "A2"
		expected := Coord{0, 1}
		if result := StringToIntCoord(coord); result.X != expected.X || result.Y != expected.Y {
			t.Errorf("Expected %d %d but got %d %d", expected.X, expected.Y, result.X, result.Y)
		}
	})

	t.Run("Right edge coord", func(t *testing.T) {
		coord := "J8"
		expected := Coord{9, 7}
		if result := StringToIntCoord(coord); result.X != expected.X || result.Y != expected.Y {
			t.Errorf("Expected %d %d but got %d %d", expected.X, expected.Y, result.X, result.Y)
		}
	})
	t.Run("Corner coord", func(t *testing.T) {
		coord := "A10"
		expected := Coord{0, 9}
		if result := StringToIntCoord(coord); result.X != expected.X || result.Y != expected.Y {
			t.Errorf("Expected %d %d but got %d %d", expected.X, expected.Y, result.X, result.Y)
		}
	})
}

func TestIntCoordToString(t *testing.T) {
	t.Run("Middle coord", func(t *testing.T) {
		coord := Coord{4, 4}
		expected := "E5"
		if result := IntCoordToString(coord); result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})
	t.Run("Left edge coord", func(t *testing.T) {
		coord := Coord{0, 1}
		expected := "A2"
		if result := IntCoordToString(coord); result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})

	t.Run("Right edge coord", func(t *testing.T) {
		coord := Coord{9, 7}
		expected := "J8"
		if result := IntCoordToString(coord); result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})
	t.Run("Corner coord", func(t *testing.T) {
		coord := Coord{0, 9}
		expected := "A10"
		if result := IntCoordToString(coord); result != expected {
			t.Errorf("Expected %s but got %s", expected, result)
		}
	})
}
