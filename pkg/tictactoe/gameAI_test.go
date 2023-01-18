// Copyright (c) 2022 Haute école d'ingerie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

// Package tictactoe implements game fonctions, data and variables
package tictactoe

import "testing"

// TestAIGetEmptyCases validate function
func TestAIGetEmptyCases(t *testing.T) {
	gameTest := &Game{}

	board := gameTest.AIGetEmptyCases()
	// make a loop and place Cross for check all empty element
	for i := 9; i > 0; i-- {
		if len(board) != i {
			t.Errorf("Error with number of EmptyCases")
		}
		gameTest.GameBoard[board[0].X][board[0].Y] = Cross
		board = gameTest.AIGetEmptyCases()
	}

	if len(board) != 0 {
		t.Errorf("Error with number of EmptyCases")
	}
}

// TestAIPlaceRandom place random symbol
func TestAIPlaceRandom(t *testing.T) {
	gameTest := &Game{}

	// check all places
	for i := 9; i > 0; i-- {
		if len(gameTest.AIGetEmptyCases()) != i {
			t.Errorf("Error: with number of EmptyCases")
		}
		gameTest.AIPlaceRandom()
	}

	if len(gameTest.AIGetEmptyCases()) != 0 {
		t.Errorf("Error: with number of EmptyCases")
	}

	countCross := 0
	countCircle := 0

	for _, array := range gameTest.GameBoard {
		for _, value := range array {
			switch value {
			case Cross:
				countCross++
			case Circle:
				countCircle++
			}
		}
	}

	// circle < cross (4 and 5)
	if countCircle < countCross {
		t.Errorf("Error: number of circle is not valid")
		t.Logf("%v %v", countCircle, countCross)
	}
}
