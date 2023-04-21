package model

import "fmt"

type State int

const (
	InProgress State = iota
	Won              = iota
	Lost             = iota
)

type Board struct {
	cells                        [][]cell
	size                         int
	state                        State
	closedNonBlackHoleCellsCount int
}

func (b *Board) GetSize() int {
	return b.size
}

func (b *Board) GetState() State {
	return b.state
}

func (b *Board) IsOpened(x, y int) bool {
	return b.cells[x][y].opened
}

func (b *Board) Open(x, y int) {
	// TODO: check if it's the first opened cell and rebuild board in case of we hit black hole.
	// TODO: It should not be possible to loose with the first move

	if outsideOfBoard(x, y, b.size) {
		return
	}
	c := &b.cells[x][y]
	if c.opened {
		return
	}
	if c.blackHole {
		c.opened = true
		b.state = Lost
		return
	}
	b.openCell(x, y)
	if b.closedNonBlackHoleCellsCount == 0 {
		b.state = Won
	}
}

func (b *Board) IsBlackHole(x, y int) bool {
	return b.cells[x][y].blackHole
}

func (b *Board) GetNeighboursCount(x, y int) int {
	return b.cells[x][y].neighboursCount
}

func outsideOfBoard(x, y, size int) bool {
	return x < 0 || x >= size || y < 0 || y >= size
}

func (b *Board) openCell(x, y int) {
	if outsideOfBoard(x, y, b.size) || b.cells[x][y].opened {
		return
	}
	b.cells[x][y].opened = true
	b.closedNonBlackHoleCellsCount--
	if b.cells[x][y].neighboursCount == 0 {
		b.openCell(x-1, y)
		b.openCell(x+1, y)
		b.openCell(x, y-1)
		b.openCell(x, y+1)
		b.openCell(x-1, y-1)
		b.openCell(x-1, y+1)
		b.openCell(x+1, y-1)
		b.openCell(x+1, y+1)
	}

}

type Point struct {
	x int
	y int
}

type CoordinatesProvider interface {
	coordinates(size, count int) ([]Point, error)
}

func NewBoard(cp CoordinatesProvider, size, blackHoleCount int) (Board, error) {
	if size <= 0 {
		return Board{}, fmt.Errorf("size should be greater then 0")
	}
	if size > 50 {
		return Board{}, fmt.Errorf("size should be less then 50")
	}
	if blackHoleCount <= 0 {
		return Board{}, fmt.Errorf("blackHoleCount should be greater then 0")
	}
	if blackHoleCount > size*size {
		return Board{}, fmt.Errorf("blackHoleCount should be less then or equal to board square (size*size)")
	}

	cells := initCells(size)

	blackHoleCoordinates, err := cp.coordinates(size, blackHoleCount)
	if err != nil {
		return Board{}, err
	}

	for _, p := range blackHoleCoordinates {
		cells[p.x][p.y].turnToBlackHole()

		markAsBlackHoleNeighbour(cells, p.x, p.y-1)
		markAsBlackHoleNeighbour(cells, p.x, p.y+1)
		markAsBlackHoleNeighbour(cells, p.x-1, p.y-1)
		markAsBlackHoleNeighbour(cells, p.x-1, p.y+1)
		markAsBlackHoleNeighbour(cells, p.x+1, p.y-1)
		markAsBlackHoleNeighbour(cells, p.x+1, p.y+1)
		markAsBlackHoleNeighbour(cells, p.x-1, p.y)
		markAsBlackHoleNeighbour(cells, p.x+1, p.y)
	}

	return Board{
		cells:                        cells,
		size:                         size,
		state:                        InProgress,
		closedNonBlackHoleCellsCount: size*size - blackHoleCount,
	}, nil
}

func markAsBlackHoleNeighbour(cells [][]cell, x, y int) {
	if outsideOfBoard(x, y, len(cells)) {
		return
	}
	cells[x][y].addNeighbour()
}

func initCells(size int) [][]cell {
	cells := make([][]cell, 0, size)
	for y := 0; y < size; y++ {
		row := make([]cell, 0, size)
		for x := 0; x < size; x++ {
			row = append(row, cell{})
		}
		cells = append(cells, row)
	}
	return cells
}
