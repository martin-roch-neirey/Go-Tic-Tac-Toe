// Copyright (c) 2022 Haute école d'ingerie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

// Package tictactoe implements game fonctions, data and variables
package tictactoe

import (
	"math/rand"
)

// A Case is a 2D value of map coordinate
type Case struct {
	X, Y int
}

// AI move type
const (
	MAX_MOVE  = +10
	MIN_MOVE  = -10
	NULL_MOVE = 0
	MAX_BEST  = -100
	MIN_BEST  = +100
)

// AIGetEmptyCases returns a array with all empty Cases of gameBoard
func (g *Game) AIGetEmptyCases() []Case {
	array := make([]Case, 0, 9)

	for x, boardY := range g.GameBoard {
		for y, _ := range boardY {
			if g.GameBoard[x][y] == None {
				array = append(array, Case{x, y})
			}
		}
	}

	return array
}

// AIPlaceRandom places a Circle randomly (AI move) on gameboard and then checks the victory condition.
func (g *Game) AIPlaceRandom() {

	emptyCase := g.AIGetEmptyCases()
	randomIndex := rand.Intn(len(emptyCase))

	g.GameBoard[emptyCase[randomIndex].X][emptyCase[randomIndex].Y] = getNowPlaying(g)
	if checkWinner(g, emptyCase[randomIndex].X, emptyCase[randomIndex].Y, getNowPlaying(g)) {
		g.GameState = Finished
		return
	}

	g.CurrentTurn++
}

// AIPlace places a Circle with a minimax Algo on gameboard and then checks the victory condition.
func (g *Game) AIPlace() {

	bestVal := MAX_BEST
	bestMove := Case{-1, -1}

	for _, v := range g.AIGetEmptyCases() {
		// Place Circle and run minimax algo
		g.GameBoard[v.X][v.Y] = Circle
		moveVal := minimax(g, false)
		// back move
		g.GameBoard[v.X][v.Y] = None

		// check best move
		if moveVal > bestVal {
			bestVal = moveVal
			bestMove = v
		}
	}

	g.GameBoard[bestMove.X][bestMove.Y] = getNowPlaying(g)
	if checkWinner(g, bestMove.X, bestMove.Y, getNowPlaying(g)) {
		g.GameState = Finished
		return
	}
	g.CurrentTurn++
}

// minimax function evaluate actions and return best choice for AI
func minimax(g *Game, isMax bool) int {

	var best int
	// get score of current game state
	score := evaluate(g)

	// if max or min score return
	if score == MAX_MOVE || score == MIN_MOVE {
		return score
	}

	// if empty case return
	if len(g.AIGetEmptyCases()) == 0 {
		return NULL_MOVE
	}

	// maximize Circle to win
	if isMax {
		best = MAX_BEST
		for _, v := range g.AIGetEmptyCases() {
			// place Circle
			g.GameBoard[v.X][v.Y] = Circle
			// recursive minimax
			best = max(minimax(g, !isMax), best)
			// remove Circle
			g.GameBoard[v.X][v.Y] = None
		}
	} else {
		best = MIN_BEST

		for _, v := range g.AIGetEmptyCases() {
			// place Cross
			g.GameBoard[v.X][v.Y] = Cross
			// recursive minimax
			best = min(minimax(g, !isMax), best)
			// remove Cross
			g.GameBoard[v.X][v.Y] = None
		}
	}

	return best
}

// evaluate current state of game and give next best AI move
func evaluate(g *Game) int {
	for i := range g.GameBoard {
		// check all ligne
		if g.GameBoard[i][0] == g.GameBoard[i][1] && g.GameBoard[i][1] == g.GameBoard[i][2] {
			if g.GameBoard[i][0] == Circle {
				return MAX_MOVE
			} else if g.GameBoard[i][0] == Cross {
				return MIN_MOVE
			}
		}
		// check all colonne
		if g.GameBoard[0][i] == g.GameBoard[1][i] && g.GameBoard[1][i] == g.GameBoard[2][i] {
			if g.GameBoard[0][i] == Circle {
				return MAX_MOVE
			} else if g.GameBoard[0][i] == Cross {
				return MIN_MOVE
			}
		}
	}

	// check diagonal 1
	if g.GameBoard[0][0] == g.GameBoard[1][1] && g.GameBoard[1][1] == g.GameBoard[2][2] {
		if g.GameBoard[0][0] == Circle {
			return MAX_MOVE
		} else if g.GameBoard[0][0] == Cross {
			return MIN_MOVE
		}
	}
	// check diagonal 2
	if g.GameBoard[0][2] == g.GameBoard[1][1] && g.GameBoard[1][1] == g.GameBoard[2][0] {
		if g.GameBoard[0][2] == Circle {
			return MAX_MOVE
		} else if g.GameBoard[0][2] == Cross {
			return MIN_MOVE
		}
	}

	return NULL_MOVE
}

// max return the largest number
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// min return the smallest number
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
