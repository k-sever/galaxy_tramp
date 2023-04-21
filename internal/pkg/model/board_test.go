package model

import (
	"fmt"
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
		size          int
		count         int
		bhCoordinates [][]int
	}
	tests := []struct {
		args         args
		want         string
		wantError    bool
		errorMessage string
	}{
		{
			args: args{
				size:          3,
				count:         2,
				bhCoordinates: [][]int{{1, 0}, {0, 2}},
			},
			want: `
				1 * 1
				2 2 1
				* 1 0
				`,
		},
		{
			args: args{
				size:          5,
				count:         10,
				bhCoordinates: [][]int{{1, 1}, {4, 0}, {1, 4}, {4, 1}, {4, 3}, {0, 0}, {1, 2}, {3, 3}, {0, 3}, {2, 0}},
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
				size:          5,
				count:         3,
				bhCoordinates: [][]int{{1, 1}, {4, 0}, {1, 4}},
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
				size:          7,
				count:         30,
				bhCoordinates: [][]int{{2, 3}, {0, 2}, {4, 5}, {3, 3}, {3, 4}, {1, 0}, {5, 2}, {3, 0}, {2, 1}, {3, 1}, {1, 5}, {3, 5}, {1, 4}, {4, 4}, {5, 4}, {0, 1}, {2, 5}, {2, 4}, {5, 3}, {0, 5}, {5, 5}, {4, 1}, {5, 1}, {4, 2}, {4, 0}, {2, 2}, {4, 3}, {3, 2}, {1, 2}, {0, 3}},
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

			got, err := NewBoard(fixedCoordinatesProvider{points: tt.args.bhCoordinates}, tt.args.size, tt.args.count)

			if !tt.wantError && err == nil {
				if !equalIgnoreSpaces(got.String(), tt.want) || err != nil {
					t.Errorf("NewBoard(%d):\n%s\nWant:\n%s", tt.args, got.String(), tt.want)
				}
			} else {
				if err.Error() != tt.errorMessage {
					t.Errorf("NewBoard(%d):\n%s\nWant:\n%s", tt.args, err.Error(), tt.errorMessage)
				}
			}

		})
	}
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
