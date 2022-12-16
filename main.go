package main

import (
	"GoTicTacToe/utils"
	"embed"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"strings"
	"time"

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

var lang = "fr-FR"
var animatedSize = 1
var animatedFont font.Face
var animatedFontList []font.Face
var listPointer = 0

type Symbol uint

// Calcule via carré magique ()
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
	Void Event = iota
	Quit
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

	switch g.gameState {
	case MainMenu:
		refreshMainMenu(g)
	case Playing:
		refreshInGame(g)
	case Finished:
		//proceedEndGame(g)
	case Pause:
		//refreshPauseMenu(g)
	}

	g.playerInput.eventType = Void
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

func refreshMainMenu(g *Game) {
	if g.playerInput.eventType == Mouse {
		g.gameState = Playing
	}
}

func refreshInGame(g *Game) {

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
				if checkWinner(g, pX, pY, getNowPlaying(g)) {

					var newBoard [3][3]Symbol
					g.gameBoard = newBoard

					g.gameState = MainMenu
					return
				}
				g.currentTurn++
			}

		}
	}
}

func checkLine(g *Game, x int, y int, dx int, dy int, len int) bool {

	var sum int

	for i := 0; i < len; i++ {
		sum += int(g.gameBoard[x+(dx*i)][y+(dy*i)])
	}

	if sum == 3 || sum == 12 {
		return true
	}

	return false
}

func checkWinner(g *Game, x int, y int, sym Symbol) bool {
	return false
}

func (g *Game) Draw(screen *ebiten.Image) {

	switch g.gameState {
	case MainMenu:
		DrawCenteredText(screen, utils.GetTranslation("game_name", lang), g.fonts["title"], WINDOW_H/2.5, color.White)
		DrawCenteredText(screen, utils.GetTranslation("click_to_play", lang), animatedFont, WINDOW_H/2, color.White)
		g.DrawSymbol(0, 0, Cross, screen)
		g.DrawSymbol(2, 2, Circle, screen)
	case Playing:
		screen.DrawImage(g.assets["map"], nil)

		for x, array := range g.gameBoard {
			for y, sym := range array {
				g.DrawSymbol(x, y, sym, screen)
			}
		}
	}

	msgFPS := strings.Replace(utils.GetTranslation("tps_fps", lang), "{tps}",
		fmt.Sprintf("%0.2f", ebiten.CurrentTPS()), 1)
	msgFPS = strings.Replace(msgFPS, "{fps}",
		fmt.Sprintf("%0.2f", ebiten.CurrentFPS()), 1)

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

	g.fonts["title"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.fonts["animated_size_1"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    33,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_1"])

	g.fonts["animated_size_2"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    34,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_2"])

	g.fonts["animated_size_3"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    35,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_3"])

	g.fonts["animated_size_4"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    36,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_4"])

	g.fonts["animated_size_5"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    37,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_5"])

	g.fonts["animated_size_6"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    38,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_6"])

	g.fonts["animated_size_7"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    39,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_7"])

	g.fonts["animated_size_8"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_8"])

	g.fonts["animated_size_9"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    41,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_9"])

	g.fonts["animated_size_10"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    42,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.fonts["animated_size_10"])

}

func (g *Game) InitGame() {
	g.gameState = MainMenu
	g.GenerateAssets()
	g.GenerateFonts()
	go g.processMainMenuAnimation()
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return WINDOW_W, WINDOW_H
}

func setAnimatedSize() {
	if listPointer == len(animatedFontList)-1 {
		reverseFontsList(animatedFontList)
		listPointer = 0
		animatedFont = animatedFontList[listPointer]
	} else {
		listPointer++
		animatedFont = animatedFontList[listPointer]
	}
}

func reverseFontsList(arr []font.Face) []font.Face {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

func (g *Game) processMainMenuAnimation() {
	for {
		setAnimatedSize()
		time.Sleep(50 * time.Millisecond)
	}
}

func DrawCenteredText(screen *ebiten.Image, s string, font font.Face, height int, color color.Color) {
	bounds := text.BoundString(font, s)
	x, y := WINDOW_W/2-bounds.Min.X-bounds.Dx()/2, height-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, color)
}

func main() {
	game := &Game{}
	game.InitGame()

	ebiten.SetWindowSize(WINDOW_W, WINDOW_H)
	ebiten.SetWindowTitle(utils.GetTranslation("game_window_name", lang))
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
