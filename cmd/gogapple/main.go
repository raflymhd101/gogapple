package main

import (
	"bufio"
	"fmt"
	"gogapple/internal/engine"
	"gogapple/internal/player"
	"os"
	"strconv"
	"strings"
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
	game := initGame()
	game.Start()
}

func initGame() *GameState {
	card := engine.GenerateCard()

	return &GameState{
		Hands:      [][]engine.Domino{card[0:7], card[7:14], card[14:21], card[21:28]},
		Names:      []string{"You", "Bot 1", "Bot 2", "Bot 3"},
		Levels:     []player.BotLevel{0, player.LevelEasy, player.LevelMedium, player.LevelHard},
		BoardLeft:  -1,
		BoardRight: -1,
		History:    []engine.Domino{},
		Turn:       0,
		PassCount:  0,
		Reader:     *bufio.NewReader(os.Stdin),
	}
}

func (g *GameState) Start() {
	fmt.Println("=== START ===")
	for {
		isGameOver := g.playTurn()
		if isGameOver {
			break
		}
	}
}

func (g *GameState) playTurn() bool {
	pIdx := g.Turn % 4
	hand := g.Hands[pIdx]
	name := g.Names[pIdx]
	level := g.Levels[pIdx]

	fmt.Print("\n=====================\n")
	fmt.Printf("%s Turn\n", name)

	if g.BoardLeft != -1 {
		fmt.Printf("Current table is: (Left: %d) ... (Right: %d)", g.BoardLeft, g.BoardRight)
	}

	validMoves := player.GetPlayableMoves(hand, g.BoardLeft, g.BoardRight)

	if len(validMoves) == 0 {
		fmt.Println("PASS")
		g.PassCount++

		if g.PassCount == 4 {
			fmt.Println("GOPLA")
			return true
		}

		g.Turn++
		return false
	}

	g.PassCount = 0
	var chosenMove *player.PlayableMove

	if level == 0 {
		chosenMove = g.promptHumanMove(hand, validMoves)
	} else {
		chosenMove = player.CalculateBotMove(hand, g.BoardLeft, g.BoardRight, g.History, level)
	}

	g.executeMove(pIdx, chosenMove)

	if len(g.Hands[pIdx]) == 0 {
		fmt.Printf("%s WIN", name)
		return true
	}

	g.Turn++
	return false
}

func (g *GameState) promptHumanMove(hand []engine.Domino, validMoves []player.PlayableMove) *player.PlayableMove {
	fmt.Println("Your Card:")

	isValidIdx := make(map[int]bool)
	for _, vm := range validMoves {
		isValidIdx[vm.HandIndex] = true
	}

	for i := range hand {
		marker := " "
		if isValidIdx[i] {
			marker = "v"
		}

		fmt.Printf("   %s%d    ", marker, i)
	}

	fmt.Println()
	for _, c := range hand {
		fmt.Printf(" %s ", c.TopPoint())
	}

	fmt.Println()
	for _, c := range hand {
		fmt.Printf(" %s ", c.BottomPoint())
	}

	fmt.Println()

	for {
		fmt.Print("\nChoose number under the 'v': ")
		inputStr, _ := g.Reader.ReadString('\n')
		inputStr = strings.TrimSpace(inputStr)

		choice, err := strconv.Atoi(inputStr)

		if err == nil && isValidIdx[choice] {
			for _, vm := range validMoves {
				if vm.HandIndex == choice {
					temp := vm
					return &temp
				}
			}
		}

		fmt.Println("INVALID INPUT")
	}
}

func (g *GameState) executeMove(pIdx int, move *player.PlayableMove) {
	hand := g.Hands[pIdx]
	playedCard := hand[move.HandIndex]

	g.Hands[pIdx] = append(hand[:move.HandIndex], hand[move.HandIndex+1:]...)
	g.History = append(g.History, playedCard)

	fmt.Printf("%s Playing card:\n%s\n", g.Names, playedCard.Card())

	if g.BoardLeft == -1 {
		g.BoardLeft = playedCard.Left
		g.BoardRight = playedCard.Right
	} else if move.Side == "left" {
		g.BoardLeft = move.LeaveNumber
	} else {
		g.BoardRight = move.LeaveNumber
	}
}
