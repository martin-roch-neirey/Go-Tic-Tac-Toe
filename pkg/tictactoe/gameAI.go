package tictactoe

type Case struct {
	X, Y int
}

func (g *Game) AIGetEmptyCases() map[Case]bool {
	m := make(map[Case]bool)

	for x, array := range g.GameBoard {
		for y, _ := range array {
			if g.GameBoard[x][y] == None {
				m[Case{x, y}] = false
			}
		}
	}
	return m
}

func (g *Game) AIPlaceRandom() {

	m := g.AIGetEmptyCases()

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

func (g *Game) AIPlace() {

	bestScore := -1000
	bestMove := Case{-1, -1}
	m := g.AIGetEmptyCases()

	for k := range m {

		g.GameBoard[k.X][k.Y] = getNowPlaying(g)
		score := minimax(g, 0, false)
		g.GameBoard[k.X][k.Y] = None

		if score > bestScore {
			bestScore = score
			bestMove = Case{k.X, k.Y}
		}
	}

	g.GameBoard[bestMove.X][bestMove.Y] = getNowPlaying(g)

	if checkWinner(g, bestMove.X, bestMove.Y, getNowPlaying(g)) {
		g.GameState = Finished
		return
	}
	g.CurrentTurn++
}

func minimax(g *Game, depth int, isMax bool) int {
	val := evaluate(g)
	best := 0

	if val == 10 || val == -10 {
		return val
	}

	m := g.AIGetEmptyCases()
	if len(m) == 0 {
		return 0
	}

	if isMax {
		best = -1000
	} else {
		best = 1000
	}

	for x, array := range g.GameBoard {
		for y, _ := range array {
			if g.GameBoard[x][y] == None {
				if isMax {
					g.GameBoard[x][y] = Circle
					best = max(best, minimax(g, depth+1, false))
				} else {
					g.GameBoard[x][y] = Cross
					best = min(best, minimax(g, depth+1, true))
				}

				g.GameBoard[x][y] = None
			}
		}
	}

	return best
}

func evaluate(g *Game) int {

	for x, array := range g.GameBoard {

		if g.GameBoard[x][0] == g.GameBoard[x][1] && g.GameBoard[x][1] == g.GameBoard[x][2] {
			if g.GameBoard[x][0] == Circle {
				return 10
			} else if g.GameBoard[x][0] == Cross {
				return -10
			}
		}

		for y, _ := range array {
			if g.GameBoard[0][y] == g.GameBoard[1][y] && g.GameBoard[1][y] == g.GameBoard[2][y] {
				if g.GameBoard[0][y] == Circle {
					return 10
				} else if g.GameBoard[0][y] == Cross {
					return -10
				}
			}
		}
	}

	if g.GameBoard[0][0] == g.GameBoard[1][1] && g.GameBoard[1][1] == g.GameBoard[2][2] {
		if g.GameBoard[0][0] == Circle {
			return 10
		} else if g.GameBoard[0][0] == Cross {
			return -10
		}
	}

	if g.GameBoard[0][2] == g.GameBoard[1][1] && g.GameBoard[1][1] == g.GameBoard[2][2] {
		if g.GameBoard[0][2] == Circle {
			return 10
		} else if g.GameBoard[0][2] == Cross {
			return -10
		}
	}

	return 0
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
