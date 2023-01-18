// Copyright (c) 2022 Haute école d'ingerie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

package tictactoe

import (
	"GoTicTacToe/pkg/api"
	"encoding/json"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

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

func getNowPlaying(g *Game) Symbol {
	if g.CurrentTurn%2 == 0 {
		return Cross
	}
	return Circle
}

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

func checkLanguageButtons(g *Game, input InputEvent) {
	if input.mouseY <= WINDOW_H && input.mouseY > WINDOW_H-15 {
		if input.mouseX >= 180 && input.mouseX < 217 {
			g.Lang = "fr-FR"
		} else if input.mouseX >= 217 && input.mouseX < 258 {
			g.Lang = "en-US"
		} else if input.mouseX >= 258 && input.mouseX < 290 {
			g.Lang = "de-DE"
		}
	}
}

func refreshLastGamesMenu(g *Game, input InputEvent) {
	if input.eventType == Mouse {
		checkLanguageButtons(g, input)
		if input.mouseY > 500 && input.mouseY < 550 {
			g.GameState = MainMenu
		}
		if input.mouseX > 350 {
			if input.mouseY > 160 && input.mouseY < 200 {
				g.LastGameEntriesViewId = 0
				g.GameState = OldBoardView
			} else if input.mouseY > 210 && input.mouseY < 250 {
				g.LastGameEntriesViewId = 1
				g.GameState = OldBoardView
			} else if input.mouseY > 260 && input.mouseY < 300 {
				g.LastGameEntriesViewId = 2
				g.GameState = OldBoardView
			} else if input.mouseY > 310 && input.mouseY < 350 {
				g.LastGameEntriesViewId = 3
				g.GameState = OldBoardView
			} else if input.mouseY > 360 && input.mouseY < 400 {
				g.LastGameEntriesViewId = 4
				g.GameState = OldBoardView
			}
		}
	} else if input.eventType == Quit {
		g.GameState = MainMenu
	}

}

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

func refreshOldBoardViewMenu(g *Game, input InputEvent) {
	if input.eventType == Mouse || input.eventType == Quit {
		checkLanguageButtons(g, input)
		g.GameState = LastGamesMenu
	}
}

func proceedEndGame(g *Game, input InputEvent) {
	g.CurrentTurn++
	if g.SqlUsable && !g.SqlProceed {
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
