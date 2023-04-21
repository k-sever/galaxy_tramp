package model

import "strconv"

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

func (c *cell) String() string {
	if c.blackHole {
		return "*"
	}
	return strconv.Itoa(c.neighboursCount)
}
