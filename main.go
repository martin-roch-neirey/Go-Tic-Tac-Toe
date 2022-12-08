package main

import (
	"embed"
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	WINDOW_W       = 480
	WINDOW_H       = 600
	LINE_THICKNESS = 6
	SYMBOL_SIZE    = 50
	FONT_SIZE      = 15
	KEY_PRESS_TIME = 60
)

type Symbol uint

// Calcule via carré magique (1 - 4 - 16)
const (
	None   Symbol = 0
	Cross  Symbol = 1
	Circle Symbol = 4
)

type Mode uint

const (
	MultiPlayer Mode = iota
	IA
)

type State uint

const (
	MainMenu State = iota
	Playing
	Finished
	Pause
)

type Event uint

const (
	Quit Event = iota
	Restart
	Mouse
)

type InputEvent struct {
	eventType Event
	mouseX    int
	mouseY    int
}

// change to generate image
//
//go:embed images/*
var imageFS embed.FS

type Game struct {
	assets      map[string]*ebiten.Image
	fonts       map[string]font.Face
	gameBoard   [3][3]Symbol
	gameState   State
	gameMode    Mode
	playerInput InputEvent
	currentTurn uint
}

func (g *Game) Update() error {

	inputAction(g)
	refreshInGame(g)

	switch g.gameState {
	case MainMenu:
		//refreshMainMenu(g)
	case Playing:
		refreshInGame(g)
	case Finished:
		//proceedEndGame(g)
	case Pause:
		//refreshPauseMenu(g)
	}

	return nil
}

func inputAction(g *Game) {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.playerInput.eventType = Mouse
		g.playerInput.mouseX, g.playerInput.mouseY = ebiten.CursorPosition()
	}

	if inpututil.KeyPressDuration(ebiten.KeyR) == KEY_PRESS_TIME {
		g.playerInput.eventType = Restart
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == KEY_PRESS_TIME {
		g.playerInput.eventType = Quit
	}
}

func getNowPlaying(g *Game) Symbol {
	if g.currentTurn%2 == 0 {
		return Cross
	}
	return Circle
}

func refreshInGame(g *Game) {

	if getWinner(g) != None {
		return
	}

	if g.playerInput.eventType == Mouse {
		// check if on game area
		if g.playerInput.mouseX > 0 &&
			g.playerInput.mouseX < WINDOW_W &&
			g.playerInput.mouseY > 0 &&
			g.playerInput.mouseY < WINDOW_W {

			pX := g.playerInput.mouseX / (WINDOW_W / 3)
			pY := g.playerInput.mouseY / (WINDOW_W / 3)

			// place a symbol
			if g.gameBoard[pX][pY] == None {
				g.gameBoard[pX][pY] = getNowPlaying(g)
				g.currentTurn++
			}

		}
	}
}

func getWinner(g *Game) Symbol {
	var sum int
	var sym int

	for i := 0; i < 3; i++ {
		sym = int(g.gameBoard[i][0])
		sum = int(g.gameBoard[i][0]) + int(g.gameBoard[i][1]) + int(g.gameBoard[i][2])
		if sum == (3 * sym) {
			return g.gameBoard[i][0]
		}
		sym = int(g.gameBoard[0][i])
		sum = int(g.gameBoard[0][i]) + int(g.gameBoard[1][i]) + int(g.gameBoard[2][i])
		if sum == (3 * sym) {
			return g.gameBoard[i][0]
		}
	}

	sym = int(g.gameBoard[1][1])
	sum = int(g.gameBoard[0][0]) + int(g.gameBoard[1][1]) + int(g.gameBoard[2][2])
	if sum == (3 * sym) {
		return g.gameBoard[1][1]
	}

	sum = int(g.gameBoard[0][2]) + int(g.gameBoard[1][1]) + int(g.gameBoard[2][0])
	if sum == (3 * sym) {
		return g.gameBoard[1][1]
	}

	return None
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.DrawImage(g.assets["map"], nil)

	for x, array := range g.gameBoard {
		for y, sym := range array {
			g.DrawSymbol(x, y, sym, screen)
		}
	}

	msgFPS := fmt.Sprintf("TPS: %0.2f\nFPS: %0.2f", ebiten.CurrentTPS(), ebiten.CurrentFPS())
	text.Draw(screen, msgFPS, g.fonts["normal"], 0, WINDOW_H-LINE_THICKNESS, color.White)

}

func (g *Game) DrawSymbol(x int, y int, sym Symbol, screen *ebiten.Image) {
	opSymbol := &ebiten.DrawImageOptions{}
	opSymbol.GeoM.Translate(float64(WINDOW_W/3)*float64(x), float64(WINDOW_W/3)*float64(y))

	switch sym {
	case Circle:
		screen.DrawImage(g.assets["circle"], opSymbol)
	case Cross:
		screen.DrawImage(g.assets["cross"], opSymbol)
	}
}

func (g *Game) GenerateAssets() {
	g.assets = make(map[string]*ebiten.Image)

	// Generate MAP
	img := gg.NewContext(WINDOW_W, WINDOW_W)
	img.SetRGB(1, 1, 1)

	img.DrawLine((WINDOW_W / 3), 0, (WINDOW_W / 3), WINDOW_W)
	img.DrawLine((WINDOW_W/3)*2, 0, (WINDOW_W/3)*2, WINDOW_W)
	img.DrawLine(0, (WINDOW_W / 3), WINDOW_W, (WINDOW_W / 3))
	img.DrawLine(0, (WINDOW_W/3)*2, WINDOW_W, (WINDOW_W/3)*2)
	img.SetLineWidth(float64(LINE_THICKNESS))
	img.Stroke()

	g.assets["map"] = ebiten.NewImageFromImage(img.Image())

	// Generate Cross
	symbolPos := float64((WINDOW_W / 3) / 2)
	img = gg.NewContext(WINDOW_W/3, WINDOW_W/3)
	img.SetRGB(1, 1, 1)

	img.DrawLine(symbolPos-SYMBOL_SIZE, symbolPos-SYMBOL_SIZE, symbolPos+SYMBOL_SIZE, symbolPos+SYMBOL_SIZE)
	img.DrawLine(symbolPos+SYMBOL_SIZE, symbolPos-SYMBOL_SIZE, symbolPos-SYMBOL_SIZE, symbolPos+SYMBOL_SIZE)
	img.SetLineWidth(float64(LINE_THICKNESS))
	img.Stroke()

	g.assets["cross"] = ebiten.NewImageFromImage(img.Image())

	// Generate Circle
	img = gg.NewContext(WINDOW_W/3, WINDOW_W/3)
	img.SetRGB(1, 1, 1)

	img.DrawCircle(symbolPos, symbolPos, SYMBOL_SIZE)
	img.SetLineWidth(float64(LINE_THICKNESS))
	img.Stroke()

	g.assets["circle"] = ebiten.NewImageFromImage(img.Image())
}

func (g *Game) GenerateFonts() {
	g.fonts = make(map[string]font.Face)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	g.fonts["normal"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    15,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {

	return WINDOW_W, WINDOW_H
}

func main() {
	game := &Game{}
	game.GenerateAssets()
	game.GenerateFonts()

	ebiten.SetWindowSize(WINDOW_W, WINDOW_H)
	ebiten.SetWindowTitle("TicTacToe")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
