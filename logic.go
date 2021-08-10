package main

import (
	"log"
	"math/rand"
)

// This function is called when you register your Battlesnake on play.battlesnake.com
// See https://docs.battlesnake.com/guides/getting-started#step-4-register-your-battlesnake
// It controls your Battlesnake appearance and author permissions.
// For customization options, see https://docs.battlesnake.com/references/personalization
func info() BattlesnakeInfoResponse {
	log.Println("INFO")
	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "belimawr",
		Color:      "#1a0a74",
		Head:       "scarf",
		Tail:       "ice-skate",
	}
}

// This function is called everytime your Battlesnake is entered into a game.
// The provided GameState contains information about the game that's about to be played.
// It's purely for informational purposes, you don't have to make any decisions here.
func start(state GameState) {
	log.Printf("%s START\n", state.Game.ID)
}

// This function is called when a game your Battlesnake was in has ended.
// It's purely for informational purposes, you don't have to make any decisions here.
func end(state GameState) {
	log.Printf("%s END\n\n", state.Game.ID)
}

// This function is called on every turn of a game. Use the provided GameState to decide
// where to move -- valid moves are "up", "down", "left", or "right".
// We've provided some code and comments to get you started.
func move(state GameState) BattlesnakeMoveResponse {
	possibleMoves := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	myHead := state.You.Body[0]

	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	if myHead.X == boardWidth-1 {
		possibleMoves["right"] = false
	}
	if myHead.Y == boardHeight-1 {
		possibleMoves["up"] = false
	}
	if myHead.X == 0 {
		possibleMoves["left"] = false
	}
	if myHead.Y == 0 {
		possibleMoves["down"] = false
	}

	// deathPoints are positions we can't go
	// Populate it with our own body
	deathPoints := append([]Coord{}, state.You.Body...)

	// Thne add all other snakes
	for _, s := range state.Board.Snakes {
		deathPoints = append(deathPoints, s.Body...)
	}

	for m, safe := range possibleMoves {
		if safe {
			switch m {
			case "up":
				nextHead := myHead
				nextHead.Y++
				for _, c := range deathPoints {
					if nextHead == c {
						possibleMoves["up"] = false
					}

				}

			case "down":
				nextHead := myHead
				nextHead.Y--
				for _, c := range deathPoints {
					if nextHead == c {
						possibleMoves["down"] = false
					}

				}

			case "left":
				nextHead := myHead
				nextHead.X--
				for _, c := range deathPoints {
					if nextHead == c {
						possibleMoves["left"] = false
					}

				}

			case "right":
				nextHead := myHead
				nextHead.X++
				for _, c := range deathPoints {
					if nextHead == c {
						possibleMoves["right"] = false
					}

				}
			}
		}
	}

	// TODO: Step 4 - Find food.
	// Use information in GameState to seek out and find food.

	// Finally, choose a move from the available safe moves.
	// TODO: Step 5 - Select a move to make based on strategy, rather than random.
	var nextMove string

	safeMoves := []string{}
	for move, isSafe := range possibleMoves {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		nextMove = "down"
		log.Printf("%s %s MOVE %d: No safe moves detected! Moving %s\n",
			state.Game.ID, state.You.ID, state.Turn, nextMove)
	} else {
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
		log.Printf("%s %s MOVE %d: %s\n", state.Game.ID, state.You.ID, state.Turn, nextMove)
	}
	return BattlesnakeMoveResponse{
		Move: nextMove,
	}
}
