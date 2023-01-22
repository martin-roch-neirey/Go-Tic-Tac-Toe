// Copyright (c) 2022 Haute école d'ingénierie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

package main

import (
	"GoTicTacToe/pkg/tictactoe"
)

func main() {
	// instantiation of a new game and initialisation of it
	game := &tictactoe.Game{}
	game.InitGame()
}
