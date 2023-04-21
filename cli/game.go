package cli

import (
	"github.com/gdamore/tcell/v2"
	"github.com/k_sever/galaxy_tramp/internal/pkg/model"
	"os"
	"strconv"
	"time"
)

const XAxisStep = 2
const YAxisStep = 1

const BannerWidth = 70
const BannerHeight = 4
const BannerPadding = 10

type Game struct {
	screen   tcell.Screen
	board    model.Board
	location point
	cursor   point
}

func NewGame(boardSize, blackHolesCount int) (Game, error) {

	rcp := model.RandomCoordinatesProvider{Seed: time.Now().UnixMilli()}
	board, err := model.NewBoard(rcp, boardSize, blackHolesCount)
	if err != nil {
		return Game{}, err
	}

	s, err := tcell.NewScreen()
	if err != nil {
		return Game{}, err
	}
	if err := s.Init(); err != nil {
		return Game{}, err
	}

	return Game{
		screen:   s,
		board:    board,
		location: point{x: (BannerWidth - boardSize) / 2, y: BannerHeight},
		cursor:   point{x: (BannerWidth - boardSize) / 2, y: BannerHeight},
	}, nil
}

type point struct {
	x int
	y int
}

func (g *Game) Start() {

	defStyle := tcell.StyleDefault.Background(tcell.ColorWhiteSmoke).Foreground(tcell.ColorBlack)
	g.screen.SetStyle(defStyle)

	go g.printScreen(defStyle)

	for {
		switch event := g.screen.PollEvent().(type) {
		case *tcell.EventResize:
			g.screen.Sync()
		case *tcell.EventKey:
			g.handleEventKey(event)
		}
	}
}

func (g *Game) handleEventKey(event *tcell.EventKey) {

	if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
		g.screen.Fini()
		os.Exit(0)
	}
	if g.board.GetState() != model.InProgress {
		return
	}

	g.handleMoves(event)
}

func (g *Game) handleMoves(event *tcell.EventKey) {
	leftBoundary := g.location.x
	rightBoundary := g.location.x + g.board.GetSize()*XAxisStep
	topBoundary := g.location.y
	bottomBoundary := g.location.y + g.board.GetSize()*YAxisStep

	switch event.Key() {
	case tcell.KeyRight:
		if g.cursor.x < rightBoundary-XAxisStep {
			g.cursor.x += XAxisStep
		}
	case tcell.KeyLeft:
		if g.cursor.x >= leftBoundary+XAxisStep {
			g.cursor.x -= XAxisStep
		}
	case tcell.KeyDown:
		if g.cursor.y < bottomBoundary-YAxisStep {
			g.cursor.y += YAxisStep
		}
	case tcell.KeyUp:
		if g.cursor.y >= topBoundary+YAxisStep {
			g.cursor.y -= YAxisStep
		}
	case tcell.KeyRune:
		switch event.Rune() {
		case ' ':
			g.board.Open((g.cursor.x-g.location.x)/XAxisStep, (g.cursor.y-g.location.y)/YAxisStep)
		}
	}
}

func (g *Game) printScreen(s tcell.Style) {

	for {
		g.screen.Clear()
		g.printBanner(s, "Use arrows to navigate, space to open a cell, esc to quit")
		g.printBoard(s)
		g.printCursor(s)
		g.screen.Show()

		time.Sleep(50 * time.Millisecond)
	}
}

func (g *Game) printBoard(s tcell.Style) {
	for y := 0; y < g.board.GetSize(); y++ {
		for x := 0; x < g.board.GetSize(); x++ {
			symbol := getSymbol(&g.board, x, y)
			g.screen.SetContent(g.location.x+x*XAxisStep, g.location.y+y*YAxisStep, symbol, nil, s)
		}
	}
	if g.board.GetState() == model.Lost {
		g.printMessage(s.Foreground(tcell.ColorRed), "Oops, that was a black hole. You Lost :(")
	}
	if g.board.GetState() == model.Won {
		g.printMessage(s.Foreground(tcell.ColorDarkGreen), "Great job! You've avoided all the black holes!")
	}
}

func (g *Game) printCursor(s tcell.Style) {
	x := (g.cursor.x - g.location.x) / XAxisStep
	y := (g.cursor.y - g.location.y) / YAxisStep
	var symbol rune
	if x >= 0 && x < g.board.GetSize() && y >= 0 && y < g.board.GetSize() {
		symbol = getSymbol(&g.board, x, y)
		symbol = highlight(symbol)
	} else {
		symbol = '〇'
	}
	g.screen.SetContent(g.cursor.x, g.cursor.y, symbol, nil, s)
}

func (g *Game) printMessage(s tcell.Style, message string) {
	for i, r := range message {
		g.screen.SetContent(i+BannerPadding, 2, r, nil, s)
	}
}

func (g *Game) printBanner(s tcell.Style, info string) {
	g.screen.SetContent(0, 0, '╔', nil, s)
	g.screen.SetContent(0, BannerHeight-1, '╚', nil, s)

	g.screen.SetContent(BannerWidth, 0, '╗', nil, s)
	g.screen.SetContent(BannerWidth, BannerHeight-1, '╝', nil, s)
	for i := 1; i < BannerHeight-1; i++ {
		g.screen.SetContent(0, i, '║', nil, s)
		g.screen.SetContent(BannerWidth, i, '║', nil, s)
	}
	for i := 0; i < BannerWidth-1; i++ {
		g.screen.SetContent(i+1, 0, '═', nil, s)
		g.screen.SetContent(i+1, 3, '═', nil, s)
	}
	for i, r := range info {
		g.screen.SetContent(i+BannerPadding, 1, r, nil, s)
	}
}

func highlight(symbol rune) rune {
	switch symbol {
	case '·':
		return '⊙'
	case ' ':
		return '〇'
	case '1':
		return '①'
	case '2':
		return '②'
	case '3':
		return '③'
	case '4':
		return '④'
	case '5':
		return '⑤'
	case '6':
		return '⑥'
	case '7':
		return '⑦'
	case '8':
		return '⑧'
	}
	return symbol
}

func getSymbol(board *model.Board, x, y int) (symbol rune) {
	switch {
	case board.IsOpened(x, y) && board.IsBlackHole(x, y):
		return '⨂'
	case !board.IsOpened(x, y):
		return '·'
	default:
		return toRune(board.GetNeighboursCount(x, y))
	}
}

func toRune(i int) rune {
	if i == 0 {
		return ' '
	}
	n := strconv.Itoa(i)
	return []rune(n)[0]
}
