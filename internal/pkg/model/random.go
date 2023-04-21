package model

import (
	"fmt"
	"math/rand"
)

type RandomCoordinatesProvider struct {
	Seed int64
}

func (r RandomCoordinatesProvider) coordinates(size, count int) ([]Point, error) {
	if size <= 0 {
		return nil, fmt.Errorf("size should be greater then 0")
	}
	if count <= 0 {
		return nil, fmt.Errorf("count should be greater then 0")
	}
	if count > size*size {
		return nil, fmt.Errorf("count should be less then or equal to board square (size*size)")
	}
	c := initCoordinates(size)
	rnd := rand.New(rand.NewSource(r.Seed))

	rnd.Shuffle(len(c), func(i, j int) {
		c[i], c[j] = c[j], c[i]
	})

	return c[:count], nil
}

func initCoordinates(size int) []Point {
	c := make([]Point, 0, size*size)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			c = append(c, Point{x: x, y: y})
		}
	}
	return c
}
