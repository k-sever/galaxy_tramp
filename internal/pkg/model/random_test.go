package model

import (
	"fmt"
	"reflect"
	"testing"
)

const SEED = 123

func TestRandomCoordinatesProvider_coordinates(t *testing.T) {
	type fields struct {
		seed int64
	}
	type args struct {
		size  int
		count int
	}
	tests := []struct {
		fields       fields
		args         args
		want         []Point
		wantError    bool
		errorMessage string
	}{
		{
			fields: fields{seed: SEED},
			args:   args{size: 2, count: 1},
			want:   []Point{{x: 1, y: 0}},
		},
		{
			fields: fields{seed: SEED},
			args:   args{size: 3, count: 2},
			want:   []Point{{x: 1, y: 1}, {x: 2, y: 2}},
		},
		{
			fields: fields{seed: SEED},
			args:   args{size: 3, count: 4},
			want:   []Point{{x: 1, y: 1}, {x: 2, y: 2}, {x: 1, y: 2}, {x: 2, y: 0}},
		},
		{
			fields: fields{seed: SEED},
			args:   args{size: 3, count: 9},
			want:   []Point{{x: 1, y: 1}, {x: 2, y: 2}, {x: 1, y: 2}, {x: 2, y: 0}, {x: 0, y: 2}, {x: 0, y: 1}, {x: 1, y: 0}, {x: 0, y: 0}, {x: 2, y: 1}},
		},
		{
			fields: fields{seed: SEED},
			args:   args{size: 5, count: 10},
			want:   []Point{{x: 1, y: 3}, {x: 0, y: 4}, {x: 3, y: 0}, {x: 2, y: 3}, {x: 4, y: 3}, {x: 2, y: 4}, {x: 0, y: 3}, {x: 3, y: 3}, {x: 4, y: 1}, {x: 0, y: 2}},
		},
		{
			fields:       fields{seed: SEED},
			args:         args{size: 10, count: 101},
			wantError:    true,
			errorMessage: "count should be less then or equal to board square (size*size)",
		},
		{
			fields:       fields{seed: SEED},
			args:         args{size: 10, count: 0},
			wantError:    true,
			errorMessage: "count should be greater then 0",
		},
		{
			fields:       fields{seed: SEED},
			args:         args{size: 10, count: -2},
			wantError:    true,
			errorMessage: "count should be greater then 0",
		},
		{
			fields:       fields{seed: SEED},
			args:         args{size: 0, count: 5},
			wantError:    true,
			errorMessage: "size should be greater then 0",
		},
		{
			fields:       fields{seed: SEED},
			args:         args{size: -1, count: 5},
			wantError:    true,
			errorMessage: "size should be greater then 0",
		},
		{
			fields:       fields{seed: SEED},
			args:         args{size: 0, count: 0},
			wantError:    true,
			errorMessage: "size should be greater then 0",
		},
	}
	for _, tt := range tests {
		name := fmt.Sprintf("board:%dx%d;count:%d", tt.args.size, tt.args.size, tt.args.count)
		t.Run(name, func(t *testing.T) {
			r := RandomCoordinatesProvider{
				Seed: tt.fields.seed,
			}
			got, err := r.coordinates(tt.args.size, tt.args.count)
			if !tt.wantError && err == nil {
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("coordinates() = %v, want %v", got, tt.want)
				}
			} else {
				if err.Error() != tt.errorMessage {
					t.Errorf("coordinates() = %v, want %v", err, tt.errorMessage)
				}
			}

		})
	}
}
