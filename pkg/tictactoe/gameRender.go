// Copyright (c) 2022 Haute école d'ingénierie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

package tictactoe

import (
	"GoTicTacToe/pkg/api"
	"GoTicTacToe/resources"
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// Draw function only applies to a game instance.
// call Draw specific view function depending on game state
func (g *Game) Draw(screen *ebiten.Image) {

	switch g.GameState {
	case MainMenu:
		drawMainMenuView(screen, g)
	case Playing:
		drawPlayingView(screen, g)
	case Finished:
		drawFinishedView(screen, g)
	case LastGamesMenu:
		drawLastGamesMenuView(screen, g)
	case OldBoardView:
		drawOldBoardView(screen, g)
	}
	msgFPS := strings.Replace(api.GetTranslation("fps", g.Lang), "{fps}",
		fmt.Sprintf("%0.2f", ebiten.CurrentFPS()), 1)
	DrawLeftText(screen, msgFPS, g.Fonts["normal"], WINDOW_H-LINE_THICKNESS, color.White)

	DrawCenteredText(screen, "FR  |  EN  |  DE", g.Fonts["normal"], WINDOW_H-LINE_THICKNESS, color.White)
}

// drawOldBoardView draws view when player wants to see the gameboard of an old played game
func drawOldBoardView(screen *ebiten.Image, g *Game) {
	g.DrawGameBoard(screen)
	DrawCenteredText(screen, api.GetTranslation("details_of_game", g.Lang), g.Fonts["normal"], WINDOW_H-15*LINE_THICKNESS, color.White)
	DrawCenteredText(screen, api.GetTranslation("click_anywhere_to_dismiss", g.Lang), g.Fonts["normal"], WINDOW_H*0.93, color.White)
	var oldGameBoard [][]int
	var oldGameEntry OldGameEntry
	oldGameEntry = g.LastGameEntries[g.LastGameEntriesViewId]
	oldGameBoard = oldGameEntry.Board

	for i := 0; i < 3; i++ {
		row := oldGameBoard[i]
		for j := 0; j < 3; j++ {
			if row[j] == 1 {
				g.DrawSymbol(i, j, Cross, screen)
			} else if row[j] == 4 {
				g.DrawSymbol(i, j, Circle, screen)
			}

		}
	}
}

// drawLastGamesMenuView draws 5 last games played menu
func drawLastGamesMenuView(screen *ebiten.Image, g *Game) {
	DrawCenteredText(screen, api.GetTranslation("last_played_games", g.Lang), g.Fonts["button"], WINDOW_H/10, color.White)
	msgNumberOfGames := strings.Replace(api.GetTranslation("number_of_games", g.Lang), "{value}",
		strconv.Itoa(api.GetGamesCount()), 1)
	DrawCenteredText(screen, msgNumberOfGames, g.Fonts["normal"], WINDOW_H*0.75, color.White)
	DrawCenteredText(screen, api.GetTranslation("click_to_dismiss", g.Lang), g.Fonts["button"], WINDOW_H*0.85, color.White)
	DrawLeftText(screen, api.GetTranslation("last_games_mode", g.Lang), g.Fonts["button"], WINDOW_H*0.2, color.White)
	DrawCenteredText(screen, api.GetTranslation("last_games_victory", g.Lang), g.Fonts["button"], WINDOW_H*0.2, color.White)
	DrawRightText(screen, api.GetTranslation("last_games_board", g.Lang), g.Fonts["button"], WINDOW_H*0.2, color.White)

	var lastGames []string
	lastGames = api.GetLastGames()
	var lastGamesEntries []OldGameEntry
	offset := 0 // height offset to write last games entries
	for i := 0; i < 5; i++ {
		offset = 50 * i
		gameJson := lastGames[i]
		entry := OldGameEntry{}
		var gameMode string
		var winner string
		json.Unmarshal([]byte(gameJson), &entry)
		switch entry.Mode {
		case 0:
			gameMode = "Multiplayer"
		case 1:
			gameMode = "IARandom"
		default:
			gameMode = "IA"
		}

		winner = entry.Winner

		lastGamesEntries = append(lastGamesEntries, entry)

		DrawLeftText(screen, gameMode, g.Fonts["normal"], WINDOW_H*0.3+offset, color.White)
		DrawCenteredText(screen, winner, g.Fonts["normal"], WINDOW_H*0.3+offset, color.White)
		translation := api.GetTranslation("click_to_see", g.Lang)
		DrawRightText(screen, translation, g.Fonts["normal"], WINDOW_H*0.3+offset, color.White)
	}
	g.LastGameEntries = lastGamesEntries
}

// drawFinishedView
// It draws the game board, the win rod and the winner message
func drawFinishedView(screen *ebiten.Image, g *Game) {
	g.DrawGameBoard(screen)
	g.drawWinRod(screen)
	if g.Winner != "/" {
		winnerTraduction := strings.Replace(api.GetTranslation("game_message_win", g.Lang), "{winner}",
			g.Winner, 1)
		DrawCenteredText(screen, winnerTraduction, g.Fonts["button"], WINDOW_H-50, color.White)
	} else {
		DrawCenteredText(screen, api.GetTranslation("game_message_no_winner", g.Lang),
			g.Fonts["button"], WINDOW_H-50, color.White)
	}
}

// drawPlayingView draws the game board and the marks count
func drawPlayingView(screen *ebiten.Image, g *Game) {
	g.DrawGameBoard(screen)
	msgMarks := strings.Replace(api.GetTranslation("marks", g.Lang), "{xMarks}",
		fmt.Sprintf("%v", g.XMarks), 1)
	msgMarks = strings.Replace(msgMarks, "{oMarks}",
		fmt.Sprintf("%v", g.OMarks), 1)
	DrawCenteredText(screen, msgMarks, g.Fonts["normal"], WINDOW_H-10*LINE_THICKNESS, color.White)
}

// drawMainMenuView draws the main menu view
func drawMainMenuView(screen *ebiten.Image, g *Game) {
	text.Draw(screen, "TIC TAC TOE", g.Fonts["title"], WINDOW_W/4, WINDOW_H/3, color.White)
	DrawCenteredText(screen, api.GetTranslation("click_to_play", g.Lang), animatedFont, WINDOW_H/2.5, color.White)
	g.DrawSymbol(0, 0, Cross, screen)
	g.DrawSymbol(2, 0, Circle, screen)
	g.DrawGameModeSelection(screen)
	if g.SqlUsable {
		DrawCenteredText(screen, api.GetTranslation("recent_games_button", g.Lang), g.Fonts["button"], WINDOW_H*0.8, color.White)
	}
	DrawRightText(screen, api.GetTranslation("button_exiting", g.Lang), g.Fonts["normal"], WINDOW_H-LINE_THICKNESS, color.White)
}

// DrawGameBoard draws the game board.
func (g *Game) DrawGameBoard(screen *ebiten.Image) {
	screen.DrawImage(g.Assets["map"], nil)

	for x, array := range g.GameBoard {
		for y, sym := range array {
			g.DrawSymbol(x, y, sym, screen)
		}
	}
}

// DrawSymbol draws the given symbol at given location on given screen
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

// drawWinRod draws winrod depending on game winrod type
func (g *Game) drawWinRod(screen *ebiten.Image) {
	opSymbol := &ebiten.DrawImageOptions{}

	switch g.WinRod.rodType {
	case HRod:
		opSymbol.GeoM.Translate(0, float64(WINDOW_W/3)*float64(g.WinRod.location))
		screen.DrawImage(g.Assets["win_bar_h"], opSymbol)

	case VRod:
		opSymbol.GeoM.Translate(float64(WINDOW_W/3)*float64(g.WinRod.location), 0)
		screen.DrawImage(g.Assets["win_bar_v"], opSymbol)

	case D1Rod:
		screen.DrawImage(g.Assets["win_bar_d1"], opSymbol)

	case D2Rod:
		screen.DrawImage(g.Assets["win_bar_d2"], opSymbol)
	}
}

// GenerateAssets is a one-time called function that creates all assets used by game
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

	img = gg.NewContext(WINDOW_W/3, WINDOW_W)
	img.SetRGBA255(45, 255, 45, 255)
	img.DrawRectangle(symbolPos-LINE_THICKNESS/2, LINE_THICKNESS, LINE_THICKNESS, WINDOW_W-LINE_THICKNESS)
	img.Fill()

	g.Assets["win_bar_v"] = ebiten.NewImageFromImage(img.Image())

	img = gg.NewContext(WINDOW_W, WINDOW_W/3)
	img.SetRGBA255(45, 255, 45, 255)
	img.DrawRectangle(LINE_THICKNESS, symbolPos-(LINE_THICKNESS/2), WINDOW_W-LINE_THICKNESS, LINE_THICKNESS)
	img.Fill()

	g.Assets["win_bar_h"] = ebiten.NewImageFromImage(img.Image())

	img = gg.NewContext(WINDOW_W, WINDOW_W)
	img.SetRGBA255(45, 255, 45, 255)
	img.DrawLine(2*LINE_THICKNESS, 2*LINE_THICKNESS, WINDOW_W-(2*LINE_THICKNESS), WINDOW_W-(2*LINE_THICKNESS))
	img.SetLineWidth(float64(LINE_THICKNESS))
	img.Stroke()

	g.Assets["win_bar_d1"] = ebiten.NewImageFromImage(img.Image())

	img = gg.NewContext(WINDOW_W, WINDOW_W)
	img.SetRGBA255(45, 255, 45, 255)
	img.DrawLine(WINDOW_W-(2*LINE_THICKNESS), 2*LINE_THICKNESS, (2 * LINE_THICKNESS), WINDOW_W-(2*LINE_THICKNESS))
	img.SetLineWidth(float64(LINE_THICKNESS))
	img.Stroke()

	g.Assets["win_bar_d2"] = ebiten.NewImageFromImage(img.Image())

}

// GenerateFonts is a one-time called function that creates all fonts used by game
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

	g.Fonts["subtitle"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	g.Fonts["button"], err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    28,
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

// Layout is a function used by ebiten
func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return WINDOW_W, WINDOW_H
}

// setAnimatedSize is called to switch animated size
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

// reverseFontsList reverts an array
func reverseFontsList(arr []font.Face) []font.Face {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
	return arr
}

// processMainMenuAnimation starts the animation of shiny text
func (g *Game) processMainMenuAnimation() {
	for {
		setAnimatedSize()
		time.Sleep(40 * time.Millisecond)
	}
}

// DrawCenteredText is a help function used to draw text centered on screen
func DrawCenteredText(screen *ebiten.Image, s string, font font.Face, height int, color color.Color) {
	bounds := text.BoundString(font, s)
	x, y := WINDOW_W/2-bounds.Min.X-bounds.Dx()/2, height-bounds.Min.Y-bounds.Dy()/2
	text.Draw(screen, s, font, x, y, color)
}

// DrawLeftText is a help function used to draw text on left part of screen
func DrawLeftText(screen *ebiten.Image, s string, font font.Face, height int, color color.Color) {
	bounds := text.BoundString(font, s)
	x, y := 20, height-bounds.Min.Y-bounds.Dy()/2 // 20 of left padding
	text.Draw(screen, s, font, x, y, color)
}

// DrawRightText is a help function used to draw text on right part of screen
func DrawRightText(screen *ebiten.Image, s string, font font.Face, height int, color color.Color) {
	bounds := text.BoundString(font, s)
	x, y := WINDOW_W-bounds.Min.X-bounds.Dx()-20, height-bounds.Min.Y-bounds.Dy()/2 // 20 of right padding
	text.Draw(screen, s, font, x, y, color)
}

// DrawGameModeSelection draws gamemode selection buttons
func (g *Game) DrawGameModeSelection(screen *ebiten.Image) {
	selectedColor := color.RGBA{45, 255, 45, 200}
	heightOffset := int(WINDOW_H * 0.6)

	switch g.GameMode {
	case MultiPlayer:
		DrawLeftText(screen, "Multiplayer", g.Fonts["button"], heightOffset, selectedColor)
		DrawCenteredText(screen, "IA", g.Fonts["button"], heightOffset, color.White)
		DrawRightText(screen, "IARandom", g.Fonts["button"], heightOffset, color.White)
	case IA:
		DrawLeftText(screen, "Multiplayer", g.Fonts["button"], heightOffset, color.White)
		DrawCenteredText(screen, "IA", g.Fonts["button"], heightOffset, selectedColor)
		DrawRightText(screen, "IARandom", g.Fonts["button"], heightOffset, color.White)
	case IARandom:
		DrawLeftText(screen, "Multiplayer", g.Fonts["button"], heightOffset, color.White)
		DrawCenteredText(screen, "IA", g.Fonts["button"], heightOffset, color.White)
		DrawRightText(screen, "IARandom", g.Fonts["button"], heightOffset, selectedColor)
	}
}

// setupWindow is a function used to set up game window
func setupWindow(g *Game) {
	ebiten.SetWindowSize(WINDOW_W, WINDOW_H)
	var favicon []image.Image

	imageBytes, err := resources.ImageFS.ReadFile("images/tic-tac-toe.png")
	if err != nil {
		log.Fatal(err)
	}
	decoded, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	favicon = append(favicon, decoded)

	ebiten.SetWindowIcon(favicon)
	ebiten.SetWindowTitle(api.GetTranslation("game_window_name", g.Lang))
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
