package main

import (
	"bufio"
	"gogapple/internal/engine"
	"gogapple/internal/game"
	"gogapple/internal/player"
)

type GameState struct {
	Hands      [][]engine.Domino
	Names      []string
	Levels     []player.BotLevel
	BoardLeft  int
	BoardRight int
	History    []engine.Domino
	Turn       int
	PassCount  int
	Reader     bufio.Reader
}

func main() {
	game := game.NewGame()
	game.Start()
}
