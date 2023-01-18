// Copyright (c) 2022 Haute école d'ingénierie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

// Package tictactoe implements game fonctions, data and variables
package tictactoe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// The function checks if the current player has won the game
func TestCheckWinner(t *testing.T) {
	gameTest := &Game{}
	currentSym := [2]Symbol{Cross, Circle}
	var emptyBoard [3][3]Symbol

	for _, v := range currentSym {
		// check all win condition
		for y := range gameTest.GameBoard {
			for x := range gameTest.GameBoard {
				gameTest.GameBoard[x][y] = v
			}

			for x := range gameTest.GameBoard {
				assert.Equal(t, true, checkWinner(gameTest, x, y, v), "should be a win condition")
			}

			gameTest.GameBoard = emptyBoard
		}

		// Test win Diag1
		gameTest.GameBoard[0][0] = v
		gameTest.GameBoard[1][1] = v
		gameTest.GameBoard[2][2] = v

		assert.Equal(t, true, checkWinner(gameTest, 0, 0, v))
		assert.Equal(t, true, checkWinner(gameTest, 1, 1, v))
		assert.Equal(t, true, checkWinner(gameTest, 2, 2, v))

		gameTest.GameBoard = emptyBoard
		// Test win Diag2
		gameTest.GameBoard[0][2] = v
		gameTest.GameBoard[1][1] = v
		gameTest.GameBoard[2][0] = v

		assert.Equal(t, true, checkWinner(gameTest, 0, 2, v))
		assert.Equal(t, true, checkWinner(gameTest, 1, 1, v))
		assert.Equal(t, true, checkWinner(gameTest, 2, 0, v))

		gameTest.GameBoard = emptyBoard
	}
}

// It tests that the function getNowPlaying returns the correct value based on the current turn
func TestGetNowPlaying(t *testing.T) {
	gameTest := &Game{}

	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			assert.Equal(t, Cross, getNowPlaying(gameTest))
		} else {
			assert.Equal(t, Circle, getNowPlaying(gameTest))
		}
		gameTest.CurrentTurn++
	}
}

// Test mouse event in MainMenu
func TestRefreshMainMenu(t *testing.T) {
	gameTest := &Game{}
	testcases := []struct {
		in    InputEvent
		state State
		mode  Mode
	}{
		// default
		{InputEvent{Event(None), 0, 0}, MainMenu, IA},
		// chose Multiplayer
		{InputEvent{Event(Mouse), 180, 370}, MainMenu, MultiPlayer},
		// chose IA
		{InputEvent{Event(Mouse), 270, 370}, MainMenu, IA},
		// chose IARandom
		{InputEvent{Event(Mouse), 330, 370}, MainMenu, IARandom},
		// empty click
		{InputEvent{Event(Mouse), 330, 200}, MainMenu, IARandom},
	}

	gameTest.GameState = MainMenu
	gameTest.GameMode = IA

	for _, tc := range testcases {

		refreshMainMenu(gameTest, tc.in)
		assert.Equal(t, tc.state, gameTest.GameState)
		assert.Equal(t, tc.mode, gameTest.GameMode)
	}
}
