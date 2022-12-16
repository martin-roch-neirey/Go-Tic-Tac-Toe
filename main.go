package main

import (
	"GoTicTacToe/utils"
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"io/ioutil"
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
	EventType Event
	MouseX    int
	MouseY    int
}

// change to generate image
//
//go:embed images/*
var imageFS embed.FS

type Game struct {
	Assets      map[string]*ebiten.Image
	Fonts       map[string]font.Face
	GameBoard   [3][3]Symbol
	GameState   State
	GameMode    Mode
	PlayerInput InputEvent
	CurrentTurn uint
	XMarks      uint
	OMarks      uint
	XWins       uint
	OWins       uint
}

func (g *Game) Update() error {

	inputAction(g)

	switch g.GameState {
	case MainMenu:
		refreshMainMenu(g)
	case Playing:
		refreshInGame(g)

		file, _ := json.MarshalIndent(g, "", " ")
		_ = ioutil.WriteFile("test.json", file, 0644)
	case Finished:
		proceedEndGame(g)
	case Pause:
		//refreshPauseMenu(g)
	}

	g.PlayerInput.EventType = Void
	return nil
}

func proceedEndGame(g *Game) {

}

func inputAction(g *Game) {

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.PlayerInput.EventType = Mouse
		g.PlayerInput.MouseX, g.PlayerInput.MouseY = ebiten.CursorPosition()
	}

	if inpututil.KeyPressDuration(ebiten.KeyR) == KEY_PRESS_TIME {
		g.PlayerInput.EventType = Restart
	}
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == KEY_PRESS_TIME {
		g.PlayerInput.EventType = Quit
	}
}

func getNowPlaying(g *Game) Symbol {
	if g.CurrentTurn%2 == 0 {
		return Cross
	}
	return Circle
}

func refreshMainMenu(g *Game) {
	if g.PlayerInput.EventType == Mouse {
		g.GameState = Playing
	}
}

func refreshInGame(g *Game) {

	if g.PlayerInput.EventType == Mouse {
		// check if on game area
		if g.PlayerInput.MouseX > 0 &&
			g.PlayerInput.MouseX < WINDOW_W &&
			g.PlayerInput.MouseY > 0 &&
			g.PlayerInput.MouseY < WINDOW_W {

			pX := g.PlayerInput.MouseX / (WINDOW_W / 3)
			pY := g.PlayerInput.MouseY / (WINDOW_W / 3)

			// place a symbol
			if g.GameBoard[pX][pY] == None {
				g.GameBoard[pX][pY] = getNowPlaying(g)
				if checkWinner(g, pX, pY, getNowPlaying(g)) {

					var newBoard [3][3]Symbol
					g.GameBoard = newBoard

					g.GameState = MainMenu
					return
				}
				g.CurrentTurn++
			}

		}
	}
}

func checkLine(g *Game, x int, y int, dx int, dy int, len int) bool {

	var sum int

	for i := 0; i < len; i++ {
		sum += int(g.GameBoard[x+(dx*i)][y+(dy*i)])
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

	switch g.GameState {
	case MainMenu:
		DrawCenteredText(screen, utils.GetTranslation("game_name", lang), g.Fonts["title"], WINDOW_H/2.5, color.White)
		DrawCenteredText(screen, utils.GetTranslation("click_to_play", lang), animatedFont, WINDOW_H/2, color.White)
		g.DrawSymbol(0, 0, Cross, screen)
		g.DrawSymbol(2, 2, Circle, screen)
	case Playing:
		screen.DrawImage(g.Assets["map"], nil)

		for x, array := range g.GameBoard {
			for y, sym := range array {
				g.DrawSymbol(x, y, sym, screen)
			}
		}
	}

	msgFPS := strings.Replace(utils.GetTranslation("tps_fps", lang), "{tps}",
		fmt.Sprintf("%0.2f", ebiten.CurrentTPS()), 1)
	msgFPS = strings.Replace(msgFPS, "{fps}",
		fmt.Sprintf("%0.2f", ebiten.CurrentFPS()), 1)

	text.Draw(screen, msgFPS, g.Fonts["normal"], 0, WINDOW_H-LINE_THICKNESS, color.White)

}

func (g *Game) DrawSymbol(x int, y int, sym Symbol, screen *ebiten.Image) {

	opSymbol := &ebiten.DrawImageOptions{}
	opSymbol.GeoM.Translate(float64(WINDOW_W/3)*float64(x), float64(WINDOW_W/3)*float64(y))

	switch sym {
	case Circle:
		screen.DrawImage(g.Assets["circle"], opSymbol)
	case Cross:
		screen.DrawImage(g.Assets["cross"], opSymbol)
	}
}

func (g *Game) GenerateAssets() {
	g.Assets = make(map[string]*ebiten.Image)

	// Generate MAP
	img := gg.NewContext(WINDOW_W, WINDOW_W)
	img.SetRGB(1, 1, 1)

	img.DrawLine((WINDOW_W / 3), 0, (WINDOW_W / 3), WINDOW_W)
	img.DrawLine((WINDOW_W/3)*2, 0, (WINDOW_W/3)*2, WINDOW_W)
	img.DrawLine(0, (WINDOW_W / 3), WINDOW_W, (WINDOW_W / 3))
	img.DrawLine(0, (WINDOW_W/3)*2, WINDOW_W, (WINDOW_W/3)*2)
	img.SetLineWidth(float64(LINE_THICKNESS))
	img.Stroke()

	g.Assets["map"] = ebiten.NewImageFromImage(img.Image())

	// Generate Cross
	symbolPos := float64((WINDOW_W / 3) / 2)
	img = gg.NewContext(WINDOW_W/3, WINDOW_W/3)
	img.SetRGB(1, 1, 1)

	img.DrawLine(symbolPos-SYMBOL_SIZE, symbolPos-SYMBOL_SIZE, symbolPos+SYMBOL_SIZE, symbolPos+SYMBOL_SIZE)
	img.DrawLine(symbolPos+SYMBOL_SIZE, symbolPos-SYMBOL_SIZE, symbolPos-SYMBOL_SIZE, symbolPos+SYMBOL_SIZE)
	img.SetLineWidth(float64(LINE_THICKNESS))
	img.Stroke()

	g.Assets["cross"] = ebiten.NewImageFromImage(img.Image())

	// Generate Circle
	img = gg.NewContext(WINDOW_W/3, WINDOW_W/3)
	img.SetRGB(1, 1, 1)

	img.DrawCircle(symbolPos, symbolPos, SYMBOL_SIZE)
	img.SetLineWidth(float64(LINE_THICKNESS))
	img.Stroke()

	g.Assets["circle"] = ebiten.NewImageFromImage(img.Image())
}

func (g *Game) GenerateFonts() {
	g.Fonts = make(map[string]font.Face)

	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	g.Fonts["normal"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    15,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.Fonts["title"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.Fonts["animated_size_1"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    37.5,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_1"])

	g.Fonts["animated_size_2"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    38,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_2"])

	g.Fonts["animated_size_3"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    38.5,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_3"])

	g.Fonts["animated_size_4"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    39,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_4"])

	g.Fonts["animated_size_5"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    39.5,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_5"])

	g.Fonts["animated_size_6"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_6"])

	g.Fonts["animated_size_7"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    40.5,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_7"])

	g.Fonts["animated_size_8"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    41,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_8"])

	g.Fonts["animated_size_9"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    41.5,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_9"])

	g.Fonts["animated_size_10"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    42,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	animatedFontList = append(animatedFontList, g.Fonts["animated_size_10"])

}

func (g *Game) InitGame() {
	g.GameState = MainMenu
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
		time.Sleep(40 * time.Millisecond)
	}
}

func DrawCenteredText(screen *ebiten.Image, s string, font font.Face, height int, color color.Color) {
	bounds := text.BoundString(font, s)
	x, y := WINDOW_W/2-bounds.Min.X-bounds.Dx()/2, height-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, color)
}

func setupWindow(g *Game) {
	ebiten.SetWindowSize(WINDOW_W, WINDOW_H)
	var favicon []image.Image

	imageBytes, err := imageFS.ReadFile("images/tic-tac-toe.png")
	if err != nil {
		log.Fatal(err)
	}
	decoded, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	favicon = append(favicon, decoded)

	ebiten.SetWindowIcon(favicon)
	ebiten.SetWindowTitle(utils.GetTranslation("game_window_name", lang))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func main() {
	game := &Game{}
	game.InitGame()
	setupWindow(game)

}
