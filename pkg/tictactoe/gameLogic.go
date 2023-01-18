// Copyright (c) 2022 Haute école d'ingénierie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

package tictactoe

import (
	"GoTicTacToe/pkg/api"
	"encoding/json"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"os"
)

// Update is called approximately 60 times per second (from ebiten.CurrentTPS())
func (g *Game) Update() error {

	playerInput := getInputs()

	switch g.GameState {
	case MainMenu:
		refreshMainMenu(g, playerInput)
	case Playing:
		refreshInGame(g, playerInput)
	case Finished:
		proceedEndGame(g, playerInput)
	case LastGamesMenu:
		refreshLastGamesMenu(g, playerInput)
	case OldBoardView:
		refreshOldBoardViewMenu(g, playerInput)
	}

	return nil
}

// InitGame is called by main package after game instantiation
func (g *Game) InitGame() {
	g.GameState = MainMenu
	g.Lang = "fr-FR"
	g.GenerateAssets()
	g.GenerateFonts()
	go g.processMainMenuAnimation()
	g.GameMode = IA
	g.SqlUsable = api.IsSqlApiUsable()
	setupWindow(g)
}

// getInputs() returns present inputs made by player actions
func getInputs() InputEvent {

	input := InputEvent{Event(None), 0, 0}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		input.eventType = Mouse
		input.mouseX, input.mouseY = ebiten.CursorPosition()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		input.eventType = Quit
	}

	return input
}

// Returns next Symbol to be placed
func getNowPlaying(g *Game) Symbol {
	if g.CurrentTurn%2 == 0 {
		return Cross
	}
	return Circle
}

// refreshMainMenu is a function that checks every player action
// when main menu is active
func refreshMainMenu(g *Game, input InputEvent) {
	if input.eventType == Mouse {
		checkLanguageButtons(g, input)
		if input.mouseY >= 230 && input.mouseY < 270 { // check start button
			g.GameState = Playing
		} else if input.mouseY >= 470 && input.mouseY < 510 && g.SqlUsable { // check game history button
			g.GameState = LastGamesMenu
		} else if input.mouseX >= 350 && input.mouseY > 570 { // check game exiting button
			os.Exit(0)
		} else if input.mouseY >= 350 && input.mouseY < 390 { // check gamemode section
			if input.mouseX < 190 { // check multiplayer button
				g.GameMode = MultiPlayer
			} else if input.mouseX > 300 { // check IARandom button
				g.GameMode = IARandom
			} else { // if not multiplayer and not IARandom, then it is IA
				g.GameMode = IA
			}
		}
	}
}

// checkLanguageButtons is a function called in every view refresh menu
// function to check if player is clicking on a language switch button
func checkLanguageButtons(g *Game, input InputEvent) {
	if input.mouseY <= WINDOW_H && input.mouseY > WINDOW_H-15 {
		if input.mouseX >= 180 && input.mouseX < 217 { // check FR button
			g.Lang = "fr-FR"
		} else if input.mouseX >= 217 && input.mouseX < 258 { // check EN button
			g.Lang = "en-US"
		} else if input.mouseX >= 258 && input.mouseX < 290 { // check DE button
			g.Lang = "de-DE"
		}
	}
}

// refreshLastGamesMenu is a function that checks every player action
// when last games menu is active
func refreshLastGamesMenu(g *Game, input InputEvent) {
	if input.eventType == Mouse {
		checkLanguageButtons(g, input)
		if input.mouseY > 500 && input.mouseY < 550 { // back to main menu check
			g.GameState = MainMenu
		}
		if input.mouseX > 350 {
			if input.mouseY > 160 && input.mouseY < 200 { // click on most recent played game
				g.LastGameEntriesViewId = 0
				g.GameState = OldBoardView
			} else if input.mouseY > 210 && input.mouseY < 250 { // 2nd played game
				g.LastGameEntriesViewId = 1
				g.GameState = OldBoardView
			} else if input.mouseY > 260 && input.mouseY < 300 { // 3rd played game
				g.LastGameEntriesViewId = 2
				g.GameState = OldBoardView
			} else if input.mouseY > 310 && input.mouseY < 350 { // 4th played game
				g.LastGameEntriesViewId = 3
				g.GameState = OldBoardView
			} else if input.mouseY > 360 && input.mouseY < 400 { // 5th played game
				g.LastGameEntriesViewId = 4
				g.GameState = OldBoardView
			}
		}
	} else if input.eventType == Quit { // ESC button check
		g.GameState = MainMenu
	}

}

// refreshInGame is a function that checks every player action
// when ingame menu is active
func refreshInGame(g *Game, input InputEvent) {

	if input.eventType == Mouse {
		checkLanguageButtons(g, input)
		// check if on game area
		if input.mouseX > 0 &&
			input.mouseX < WINDOW_W &&
			input.mouseY > 0 &&
			input.mouseY < WINDOW_W {

			pX := input.mouseX / (WINDOW_W / 3)
			pY := input.mouseY / (WINDOW_W / 3)

			// place a symbol
			if g.GameBoard[pX][pY] == None {
				g.GameBoard[pX][pY] = getNowPlaying(g)
				switch getNowPlaying(g) {
				case Cross:
					g.XMarks++
				case Circle:
					g.OMarks++
				}
				if checkWinner(g, pX, pY, getNowPlaying(g)) {
					g.GameState = Finished
					return
				}
				g.CurrentTurn++

				if g.CurrentTurn > 8 {
					g.GameState = Finished
					return
				}

				// IA MOVE
				if g.GameMode == IA {
					g.AIPlace()
					g.OMarks++
				} else if g.GameMode == IARandom {
					g.AIPlaceRandom()
					g.OMarks++
				}

			}

		}

	}
}

// refreshOldBoardViewMenu is a function that checks every player action
// when player watches an old gameboard
func refreshOldBoardViewMenu(g *Game, input InputEvent) {
	if input.eventType == Mouse || input.eventType == Quit {
		checkLanguageButtons(g, input)
		g.GameState = LastGamesMenu
	}
}

// proceedEndGame is active when game is over and winrod has been drawn.
// It sends to DB (if usable) all information about game one time.
func proceedEndGame(g *Game, input InputEvent) {
	g.CurrentTurn++
	if g.SqlUsable && !g.SqlProceed { // prevent sending to DB multiple times the same game
		jsonAsBytes, _ := json.Marshal(g)
		jsonString := string(jsonAsBytes[:])
		api.UploadNewGame(jsonString)
		g.SqlProceed = true
	}
	if input.eventType == Mouse || input.eventType == Quit {
		checkLanguageButtons(g, input)
		var newBoard [3][3]Symbol
		g.GameBoard = newBoard
		g.CurrentTurn = 0
		g.XMarks = 0
		g.OMarks = 0
		g.Winner = "/"
		g.WinRod.rodType = NORod
		g.GameState = MainMenu
		g.SqlProceed = false
	}
}

// checkWinner is a magical square function to check if
// a given Symbol is an actual winner of the game. Returns true if so, false otherwise.
func checkWinner(g *Game, x int, y int, sym Symbol) bool {
	sum := make([]int, 4)

	for i := range g.GameBoard {
		sum[0] += int(g.GameBoard[i][y])
	}
	for i := range g.GameBoard[x] {
		sum[1] += int(g.GameBoard[x][i])
	}

	sum[2] = int(g.GameBoard[0][0]) + int(g.GameBoard[1][1]) + int(g.GameBoard[2][2])
	sum[3] = int(g.GameBoard[0][2]) + int(g.GameBoard[1][1]) + int(g.GameBoard[2][0])

	for i, v := range sum {
		if int(v) == (int(sym) * 3) {
			switch i {
			case 0:
				g.WinRod.rodType = HRod
				g.WinRod.location = uint(y)
			case 1:
				g.WinRod.rodType = VRod
				g.WinRod.location = uint(x)
			case 2:
				g.WinRod.rodType = D1Rod
			case 3:
				g.WinRod.rodType = D2Rod
			}

			g.Winner = string(sym.toString())
			return true
		}
	}
	g.Winner = "/"
	return false
}
