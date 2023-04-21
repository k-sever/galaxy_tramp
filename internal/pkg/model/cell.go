package model

type cell struct {
	opened          bool
	blackHole       bool
	neighboursCount int
}

func (c *cell) turnToBlackHole() {
	c.blackHole = true
}

func (c *cell) addNeighbour() {
	c.neighboursCount++
}
