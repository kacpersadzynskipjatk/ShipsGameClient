# ShipsGameClient

The Ships Game is a turn-based strategy game where two players take turns firing shots at each other's ships on a grid-based board. The goal is to sink all the opponent's ships before they sink yours.

## Game Overview

The game follows these basic rules:

1. Each player has their own board with ships placed on it.
2. Players take turns firing shots at the opponent's board.
3. Hits and misses are recorded on the boards.
4. Ships are sunk when all their coordinates have been hit.
5. The game continues until all of one player's ships have been sunk.
6. The player who sinks all the opponent's ships first wins the game.

## Features

- GUI-based interface for a visual representation of the game boards.
- Interactive gameplay with options for firing shots and making strategic decisions.
- Timer functionality to limit the time available for each turn.
- Game statistics tracking to record the number of shots and hits.
- Ability to play multiple games and view player rankings.

## Getting Started

To play the Ships Game, follow these steps:

1. Clone the repository: `git clone https://github.com/kacpersadzynskipjatk/ShipsGameClient`
2. Install the necessary dependencies.
3. Build the application: `go build`
4. Run the application: `go run .`

## Game Controls

- Use the GUI interface to interact with the game boards.
- Use the keyboard to input menu choices and options during gameplay.