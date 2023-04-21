package main

import (
	"github.com/k-sever/galaxy_tramp/cli"
	"log"
	"os"
)

func main() {
	// TODO: add custom mode with arbitrary board size and holes count
	boardSize := 8
	blackHolesCount := 10
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "medium":
			boardSize = 16
			blackHolesCount = 40
		case "hard":
			boardSize = 24
			blackHolesCount = 99
		}
	}
	game, err := cli.NewGame(boardSize, blackHolesCount)
	if err != nil {
		log.Fatalf("%+v", err)
	}

	game.Start()
}
