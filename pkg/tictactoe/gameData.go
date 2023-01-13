package tictactoe

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
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

var sql = false // TODO delete later, set to false if DB is not online
var sqlProceed = false

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
	LastGamesMenu
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
	Winner      string
}
