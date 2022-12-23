package tictactoe

type Case struct {
	X, Y int
}

func AIPlaceRandom(g *Game) {
	// place a symbol
	m := make(map[Case]int)

	for x, array := range g.GameBoard {
		for y, _ := range array {
			if g.GameBoard[x][y] == None {
				m[Case{x, y}] = 0
			}
		}
	}

	for k := range m {
		g.GameBoard[k.X][k.Y] = getNowPlaying(g)
		incrementMarksCounter(g)
		if checkWinner(g, k.X, k.Y, getNowPlaying(g)) {
			g.GameState = Finished
			return
		}
		g.CurrentTurn++
		return
	}
}
