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
	MAX_MOVE  = 10
	MIN_MOVE  = 10
	NULL_MOVE = 0
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
	bestScore := -1000
	bestMove := Case{-1, -1}
	emptyCase := g.AIGetEmptyCases()

	for _, v := range emptyCase {

		g.GameBoard[v.X][v.Y] = Circle
		score := minimax(g, 0, false)
		g.GameBoard[v.X][v.Y] = None

		if score > bestScore {
			bestScore = score
			bestMove = Case{v.X, v.Y}
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
func minimax(g *Game, depth int, isMax bool) int {

	var best int
	val := evaluate(g)
	if val == MAX_MOVE {
		return 10 - depth
	}

	if val == MIN_MOVE {
		return -10 + depth
	}

	emptyCases := g.AIGetEmptyCases()
	if len(emptyCases) == 0 {
		return 0
	}

	if isMax {
		best = -1000
	} else {
		best = 1000
	}

	for x, array := range g.GameBoard {
		for y, _ := range array {
			if g.GameBoard[x][y] == None {
				if isMax {
					g.GameBoard[x][y] = Circle
					best = max(best, minimax(g, depth+1, false))
				} else {
					g.GameBoard[x][y] = Cross
					best = min(best, minimax(g, depth+1, true))
				}

				g.GameBoard[x][y] = None
			}
		}
	}

	return best
}

// evaluate current state of game and give next best AI move
func evaluate(g *Game) int {

	for x, array := range g.GameBoard {

		if g.GameBoard[x][0] == g.GameBoard[x][1] && g.GameBoard[x][1] == g.GameBoard[x][2] {
			if g.GameBoard[x][0] == Circle {
				return MAX_MOVE
			} else if g.GameBoard[x][0] == Cross {
				return MIN_MOVE
			}
		}

		for y, _ := range array {
			if g.GameBoard[0][y] == g.GameBoard[1][y] && g.GameBoard[1][y] == g.GameBoard[2][y] {
				if g.GameBoard[0][y] == Circle {
					return MAX_MOVE
				} else if g.GameBoard[0][y] == Cross {
					return MIN_MOVE
				}
			}
		}
	}

	if g.GameBoard[0][0] == g.GameBoard[1][1] && g.GameBoard[1][1] == g.GameBoard[2][2] {
		if g.GameBoard[0][0] == Circle {
			return MAX_MOVE
		} else if g.GameBoard[0][0] == Cross {
			return MIN_MOVE
		}
	}

	if g.GameBoard[0][2] == g.GameBoard[1][1] && g.GameBoard[1][1] == g.GameBoard[2][2] {
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
