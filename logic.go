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

	// Avoid walls
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

	// snakesBodies are positions we can't go
	// Populate it with our own body
	snakesBodies := []Coord{}
	snakesHeads := []Coord{}

	// The add all snakes (including ourselves)
	for _, s := range state.Board.Snakes {
		snakesBodies = append(snakesBodies, s.Body...)

		if s.ID != state.You.ID {
			snakesHeads = append(snakesHeads, s.Head)
		}
	}

	for m, safe := range possibleMoves {
		if safe {
			switch m {
			case "up":
				nextHead := myHead
				nextHead.Y++
				for _, c := range snakesBodies {
					if nextHead == c {
						possibleMoves["up"] = false
					}

				}

				// for _, c := snakesHeads{
				// }

			case "down":
				nextHead := myHead
				nextHead.Y--
				for _, c := range snakesBodies {
					if nextHead == c {
						possibleMoves["down"] = false
					}

				}

			case "left":
				nextHead := myHead
				nextHead.X--
				for _, c := range snakesBodies {
					if nextHead == c {
						possibleMoves["left"] = false
					}

				}

			case "right":
				nextHead := myHead
				nextHead.X++
				for _, c := range snakesBodies {
					if nextHead == c {
						possibleMoves["right"] = false
					}

				}
			}
		}
	}

	// Avoid hitting snakes on their next move
	for _, m := range avoidHeadColision(myHead, snakesHeads) {
		possibleMoves[m] = false
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

// avoidHeadColision returns dangerous moves, they might hit other snake's head
func avoidHeadColision(myHead Coord, snakesHeads []Coord) []string {
	moves := []string{"up", "down", "left", "right"}
	dangerousMoves := map[string]bool{
		"up":    false,
		"down":  false,
		"left":  false,
		"right": false,
	}

	for _, sh := range snakesHeads {
		for _, m := range moves {
			switch m {
			case "up":
				nextHead := Coord{X: myHead.X, Y: myHead.Y + 1}
				if willHitOtherHead(nextHead, sh) {
					dangerousMoves["up"] = true
				}
			case "down":
				nextHead := Coord{X: myHead.X, Y: myHead.Y - 1}
				if willHitOtherHead(nextHead, sh) {
					dangerousMoves["up"] = true
				}
			case "left":
				nextHead := Coord{X: myHead.X - 1, Y: myHead.Y}
				if willHitOtherHead(nextHead, sh) {
					dangerousMoves["up"] = true
				}
			case "right":
				nextHead := Coord{X: myHead.X + 1, Y: myHead.Y}
				if willHitOtherHead(nextHead, sh) {
					dangerousMoves["up"] = true
				}
			}
		}
	}

	movesList := []string{}
	for m, notSafe := range dangerousMoves {
		if notSafe {
			movesList = append(movesList, m)
		}
	}

	return movesList
}

func willHitOtherHead(me, other Coord) bool {
	others := []Coord{
		{X: other.X, Y: other.Y + 1}, // up
		{X: other.X, Y: other.Y - 1}, // down
		{X: other.X + 1, Y: other.Y}, // right
		{X: other.X - 1, Y: other.Y}, // left
	}

	for _, o := range others {
		if me == o {
			return true
		}
	}

	return false
}
