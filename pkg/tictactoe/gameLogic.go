package tictactoe

import (
	"GoTicTacToe/pkg/api"
	"encoding/json"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {

	playerInput := GetInputs()

	switch g.GameState {
	case MainMenu:
		refreshMainMenu(g, playerInput)
	case Playing:
		refreshInGame(g, playerInput)
	case Finished:
		proceedEndGame(g, playerInput)
	case Pause:
		//refreshPauseMenu(g)
	}

	return nil
}

func (g *Game) InitGame() {
	g.GameState = MainMenu
	g.GenerateAssets()
	g.GenerateFonts()
	go g.processMainMenuAnimation()
	g.GameMode = IA

	setupWindow(g)
}

func GetInputs() InputEvent {

	input := InputEvent{Event(None), 0, 0}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		input.eventType = Mouse
		input.mouseX, input.mouseY = ebiten.CursorPosition()
	}

	if inpututil.KeyPressDuration(ebiten.KeyR) == KEY_PRESS_TIME {
		input.eventType = Restart
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == KEY_PRESS_TIME {
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
		//g.GameState = Playing

		if input.mouseX > 0 &&
			input.mouseX < (WINDOW_W/3) &&
			input.mouseY > WINDOW_W &&
			input.mouseY < WINDOW_H {
			fmt.Printf("ok")
			g.GameState = Playing
		}

	}
}

func refreshInGame(g *Game, input InputEvent) {

	if input.eventType == Mouse {
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
				incrementMarksCounter(g)
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
				} else if g.GameMode == IARandom {
					g.AIPlaceRandom()
				}

			}

		}
	}
}

func proceedEndGame(g *Game, input InputEvent) {
	g.CurrentTurn++
	if sql && !sqlProceed {
		jsonAsBytes, _ := json.Marshal(g)
		jsonString := string(jsonAsBytes[:])
		api.UploadNewGame(jsonString)
		sqlProceed = true
	}
	if input.eventType == Mouse {
		var newBoard [3][3]Symbol
		g.GameBoard = newBoard
		g.CurrentTurn = 0
		g.XMarks = 0
		g.OMarks = 0
		g.WinRod.rodType = NORod
		g.GameState = MainMenu
		sqlProceed = false
	}
}

func incrementMarksCounter(g *Game) {
	if getNowPlaying(g) == Cross {
		g.XMarks++
	} else {
		g.OMarks++
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

			return true
		}
	}

	return false
}
