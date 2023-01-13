package tictactoe

import (
	"GoTicTacToe/pkg/api"
	"encoding/json"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func (g *Game) Update() error {

	inputAction(g)

	switch g.GameState {
	case MainMenu:
		refreshMainMenu(g)
	case Playing:
		refreshInGame(g)
	case Finished:
		proceedEndGame(g)
	case Pause:
		//refreshPauseMenu(g)
	case LastGamesMenu:
		refreshLastGamesMenu(g)
	}

	g.PlayerInput.eventType = Void
	return nil
}

func inputAction(g *Game) {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.PlayerInput.eventType = Mouse
		g.PlayerInput.mouseX, g.PlayerInput.mouseY = ebiten.CursorPosition()
	}

	if inpututil.KeyPressDuration(ebiten.KeyR) == KEY_PRESS_TIME {
		g.PlayerInput.eventType = Restart
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == KEY_PRESS_TIME {
		g.PlayerInput.eventType = Quit
	}
}

func getNowPlaying(g *Game) Symbol {
	if g.CurrentTurn%2 == 0 {
		return Cross
	}
	return Circle
}

func refreshMainMenu(g *Game) {
	if g.PlayerInput.eventType == Mouse {
		if g.PlayerInput.mouseY < 200 {
			g.GameState = LastGamesMenu
		} else {
			g.GameState = Playing
		}
	}
}

func refreshLastGamesMenu(g *Game) {
	if g.PlayerInput.eventType == Mouse {
		g.GameState = MainMenu
	}
}

func refreshInGame(g *Game) {

	if g.PlayerInput.eventType == Mouse {
		// check if on game area
		if g.PlayerInput.mouseX > 0 &&
			g.PlayerInput.mouseX < WINDOW_W &&
			g.PlayerInput.mouseY > 0 &&
			g.PlayerInput.mouseY < WINDOW_W {

			pX := g.PlayerInput.mouseX / (WINDOW_W / 3)
			pY := g.PlayerInput.mouseY / (WINDOW_W / 3)

			// place a symbol
			if g.GameBoard[pX][pY] == None {
				g.GameBoard[pX][pY] = getNowPlaying(g)
				incrementMarksCounter(g)
				if checkWinner(g, pX, pY, getNowPlaying(g)) {
					g.GameState = Finished
					return
				}
				g.CurrentTurn++

				if g.GameMode == IA {
					AIPlaceRandom(g)
				}

				if g.CurrentTurn > 8 {
					g.GameState = Finished
					return
				}
			}

		}

	}
}

func incrementMarksCounter(g *Game) {
	if getNowPlaying(g) == Cross {
		g.XMarks++
	} else {
		g.OMarks++
	}
}

func proceedEndGame(g *Game) {
	g.CurrentTurn++
	if sql && !sqlProceed {
		jsonAsBytes, _ := json.Marshal(g)
		jsonString := string(jsonAsBytes[:])
		api.UploadNewGame(jsonString)
		sqlProceed = true
	}
	if g.PlayerInput.eventType == Mouse {
		var newBoard [3][3]Symbol
		g.GameBoard = newBoard
		g.CurrentTurn = 0
		g.XMarks = 0
		g.OMarks = 0
		g.GameState = MainMenu
		sqlProceed = false
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

	for _, v := range sum {
		if int(v) == (int(sym) * 3) {
			return true
		}
	}

	return false
}

func (g *Game) InitGame() {
	g.GameState = MainMenu
	g.GenerateAssets()
	g.GenerateFonts()
	go g.processMainMenuAnimation()
	g.GameMode = IA

	setupWindow(g)
}
