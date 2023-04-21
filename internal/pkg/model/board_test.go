package model

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type fixedCoordinatesProvider struct {
	points [][]int
}

func (f fixedCoordinatesProvider) coordinates(_, count int) ([]Point, error) {
	res := make([]Point, count)
	for i, p := range f.points {
		res[i] = Point{x: p[0], y: p[1]}
	}
	return res, nil
}

func TestBoard_NewBoard(t *testing.T) {
	type args struct {
		size           int
		count          int
		blackHoleCells [][]int
	}
	tests := []struct {
		args         args
		want         string
		wantError    bool
		errorMessage string
	}{
		{
			args: args{
				size:           3,
				count:          2,
				blackHoleCells: [][]int{{1, 0}, {0, 2}},
			},
			want: `
				1 * 1
				2 2 1
				* 1 0
				`,
		},
		{
			args: args{
				size:           5,
				count:          10,
				blackHoleCells: [][]int{{1, 1}, {4, 0}, {1, 4}, {4, 1}, {4, 3}, {0, 0}, {1, 2}, {3, 3}, {0, 3}, {2, 0}},
			},
			want: `
				* 3 * 3 *
				3 * 3 3 *
				3 * 3 3 3
				* 3 3 * *
				2 * 2 2 2
				`,
		},
		{
			args: args{
				size:           5,
				count:          3,
				blackHoleCells: [][]int{{1, 1}, {4, 0}, {1, 4}},
			},
			want: `
				1 1 1 1 *
				1 * 1 1 1
				1 1 1 0 0
				1 1 1 0 0
				1 * 1 0 0
				`,
		},
		{
			args: args{
				size:           7,
				count:          30,
				blackHoleCells: [][]int{{2, 3}, {0, 2}, {4, 5}, {3, 3}, {3, 4}, {1, 0}, {5, 2}, {3, 0}, {2, 1}, {3, 1}, {1, 5}, {3, 5}, {1, 4}, {4, 4}, {5, 4}, {0, 1}, {2, 5}, {2, 4}, {5, 3}, {0, 5}, {5, 5}, {4, 1}, {5, 1}, {4, 2}, {4, 0}, {2, 2}, {4, 3}, {3, 2}, {1, 2}, {0, 3}},
			},
			want: `
				2 * 4 * * 3 1
				* 6 * * * * 2
				* * * * * * 3
				* 7 * * * * 3
				4 * * * * * 3
				* * * * * * 2
				2 3 3 3 3 2 1
				`,
		},
		{
			args:         args{size: 10, count: 101},
			wantError:    true,
			errorMessage: "blackHoleCount should be less then or equal to board square (size*size)",
		},
		{
			args:         args{size: 10, count: 0},
			wantError:    true,
			errorMessage: "blackHoleCount should be greater then 0",
		},
		{
			args:         args{size: 10, count: -2},
			wantError:    true,
			errorMessage: "blackHoleCount should be greater then 0",
		},
		{
			args:         args{size: 0, count: 5},
			wantError:    true,
			errorMessage: "size should be greater then 0",
		},
		{
			args:         args{size: -1, count: 5},
			wantError:    true,
			errorMessage: "size should be greater then 0",
		},
		{
			args:         args{size: 0, count: 0},
			wantError:    true,
			errorMessage: "size should be greater then 0",
		},
		{
			args:         args{size: 75, count: 20},
			wantError:    true,
			errorMessage: "size should be less then 50",
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("board:%dx%d;count:%d", tt.args.size, tt.args.size, tt.args.count)
		t.Run(name, func(t *testing.T) {

			got, err := NewBoard(fixedCoordinatesProvider{points: tt.args.blackHoleCells}, tt.args.size, tt.args.count)

			actual := boardToString(got, false)
			if !tt.wantError && err == nil {
				if !equalIgnoreSpaces(actual, tt.want) || err != nil {
					t.Errorf("NewBoard(%v):\n%s\nWant:\n%s", tt.args, actual, tt.want)
				}
			} else {
				if err.Error() != tt.errorMessage {
					t.Errorf("NewBoard(%v):\n%s\nWant:\n%s", tt.args, err.Error(), tt.errorMessage)
				}
			}

		})
	}
}

func boardToString(b Board, hideNotOpenedCells bool) string {
	r := ""
	for y := 0; y < b.size; y++ {
		for x := 0; x < b.size; x++ {
			c := b.cells[x][y]
			if hideNotOpenedCells && !c.opened {
				r += "? "
				continue
			}
			if c.blackHole {
				r += "*"
			} else {
				r += strconv.Itoa(c.neighboursCount)
			}
			r += " "
		}
		r += "\n"
	}
	return r
}

func equalIgnoreSpaces(s1, s2 string) bool {
	return replaceWhiteSpaces(s1) == replaceWhiteSpaces(s2)
}

func replaceWhiteSpaces(s string) string {
	s = strings.ReplaceAll(s, "\t", "")
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

func TestBoard_Open(t *testing.T) {
	type args struct {
		size           int
		count          int
		blackHoleCells [][]int
	}
	tests := []struct {
		name             string
		args             args
		openedCells      [][]int
		wantState        State
		wantInitialBoard string
		wantOpenedBoard  string
	}{
		{
			name: "Open cell with a number",
			args: args{

				size:           3,
				count:          2,
				blackHoleCells: [][]int{{1, 0}, {0, 2}},
			},
			openedCells: [][]int{{0, 0}},
			wantState:   InProgress,
			wantInitialBoard: `
				1 * 1
				2 2 1
				* 1 0
				`,
			wantOpenedBoard: `
				1 ? ?
				? ? ?
				? ? ?
				`,
		},
		{
			name: "Open cell with 0 neighbours",
			args: args{

				size:           3,
				count:          2,
				blackHoleCells: [][]int{{1, 0}, {0, 2}},
			},
			openedCells: [][]int{{2, 2}},
			wantState:   InProgress,
			wantInitialBoard: `
				1 * 1
				2 2 1
				* 1 0
				`,
			wantOpenedBoard: `
				? ? ?
				? 2 1
				? 1 0
				`,
		},
		{
			name: "Open few cells",
			args: args{

				size:           5,
				count:          7,
				blackHoleCells: [][]int{{1, 3}, {3, 2}, {0, 0}, {0, 1}, {4, 3}},
			},
			openedCells: [][]int{{3, 0}, {1, 2}, {4, 4}},
			wantState:   InProgress,
			wantInitialBoard: `
				* 4 0 0 0 
				* 4 1 1 1 
				2 2 2 * 2 
				1 * 2 2 * 
				1 1 1 1 1
				`,
			wantOpenedBoard: `
				? 4 0 0 0 
				? 4 1 1 1 
				? 2 ? ? ? 
				? ? ? ? ? 
				? ? ? ? 1
				`,
		},
		{
			name: "Open black hole from first hit - Lost",
			args: args{

				size:           4,
				count:          5,
				blackHoleCells: [][]int{{1, 3}, {3, 2}, {0, 0}, {0, 1}},
			},
			openedCells: [][]int{{0, 0}},
			wantState:   Lost,
			wantInitialBoard: `
				* 3 0 0 
				* 3 1 1 
				2 2 2 * 
				1 * 2 1
				`,
			wantOpenedBoard: `
				* ? ? ? 
				? ? ? ? 
				? ? ? ? 
				? ? ? ? 
				`,
		},
		{
			name: "Open black hole - Lost",
			args: args{

				size:           4,
				count:          5,
				blackHoleCells: [][]int{{1, 3}, {3, 2}, {0, 0}, {0, 1}},
			},
			openedCells: [][]int{{2, 0}, {2, 3}, {0, 1}},
			wantState:   Lost,
			wantInitialBoard: `
				* 3 0 0 
				* 3 1 1 
				2 2 2 * 
				1 * 2 1
				`,
			wantOpenedBoard: `
				? 3 0 0 
				* 3 1 1 
				? ? ? ? 
				? ? 2 ? 
				`,
		},
		{
			name: "Open all non black holes - Win",
			args: args{

				size:           4,
				count:          5,
				blackHoleCells: [][]int{{1, 3}, {3, 2}, {0, 0}, {0, 1}},
			},
			openedCells: [][]int{{2, 0}, {0, 2}, {1, 2}, {2, 2}, {0, 3}, {2, 3}, {3, 3}},
			wantState:   Won,
			wantInitialBoard: `
				* 3 0 0 
				* 3 1 1 
				2 2 2 * 
				1 * 2 1
				`,
			wantOpenedBoard: `
				? 3 0 0 
				? 3 1 1 
				2 2 2 ? 
				1 ? 2 1
				`,
		},
		{
			name: "Open invalid cell",
			args: args{

				size:           3,
				count:          2,
				blackHoleCells: [][]int{{1, 0}, {0, 2}},
			},
			openedCells: [][]int{{10, 5}},
			wantState:   InProgress,
			wantInitialBoard: `
				1 * 1
				2 2 1
				* 1 0
				`,
			wantOpenedBoard: `
				? ? ?
				? ? ?
				? ? ?
				`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			board, err := NewBoard(fixedCoordinatesProvider{points: tt.args.blackHoleCells}, tt.args.size, tt.args.count)

			actualInitial := boardToString(board, false)
			if !equalIgnoreSpaces(actualInitial, tt.wantInitialBoard) || err != nil {
				t.Errorf("Got:\n%s\nWant:\n%s", actualInitial, tt.wantInitialBoard)
			}

			for _, p := range tt.openedCells {
				board.Open(p[0], p[1])
			}

			actualOpened := boardToString(board, true)
			if !equalIgnoreSpaces(actualOpened, tt.wantOpenedBoard) || err != nil {
				t.Errorf("Got:\n%s\nWant:\n%s", actualOpened, tt.wantOpenedBoard)
			}
		})
	}
}
