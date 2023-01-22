// Copyright (c) 2022 Haute école d'ingénierie et d'architecture de Fribourg
// SPDX-License-Identifier: Apache-2.0
// Author:  William Margueron & Martin Roch-Neirey

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
)

var animatedFont font.Face
var animatedFontList []font.Face
var listPointer = 0

type Symbol uint

const (
	None   Symbol = 0
	Cross  Symbol = 1
	Circle Symbol = 4
)

func (s Symbol) toString() rune {
	switch s {
	case Cross:
		return 'X'
	case Circle:
		return 'O'
	default:
		return '/'
	}
}

type Mode uint

const (
	MultiPlayer Mode = iota
	IARandom
	IA
)

type State uint

const (
	MainMenu State = iota
	Playing
	Finished
	LastGamesMenu
	OldBoardView
)

type Event uint

const (
	Void Event = iota
	Quit
	Mouse
)

type InputEvent struct {
	eventType Event
	mouseX    int
	mouseY    int
}

type RodType uint

const (
	NORod RodType = iota
	HRod
	VRod
	D1Rod
	D2Rod
)

type WinRod struct {
	rodType  RodType
	location uint
}

type Game struct {
	Assets                map[string]*ebiten.Image
	Fonts                 map[string]font.Face
	Lang                  string
	GameBoard             [3][3]Symbol
	GameState             State
	GameMode              Mode
	WinRod                WinRod
	CurrentTurn           uint
	XMarks                uint
	OMarks                uint
	Winner                string
	SqlUsable             bool
	SqlProceed            bool
	LastGameEntries       []OldGameEntry
	LastGameEntriesViewId int
}

type OldGameEntry struct {
	Mode   int     `json:"GameMode"`
	Winner string  `json:"Winner"`
	Board  [][]int `json:"GameBoard"`
}
