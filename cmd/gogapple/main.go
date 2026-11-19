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
	fmt.Println("============= START ==============")
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

	fmt.Printf("==================================\n")

	validMoves := player.GetPlayableMoves(hand, g.BoardLeft, g.BoardRight)

	if len(validMoves) == 0 {
		fmt.Println("PASS")
		g.PassCount++

		if g.PassCount == 4 {
			fmt.Println("GOPLA")
			g.countGappleWinner()
			return true
		}

		g.Turn++
		return false
	}

	g.PassCount = 0
	var chosenMove *player.PlayableMove

	if level == 0 {

		if g.BoardLeft == -1 {
			fmt.Printf("\n======= TABLE STILL EMPTY ========\n")
		} else {
			fmt.Printf("\n=== CURRENT CARD === (Left: %d | Right: %d)\n", g.BoardLeft, g.BoardRight)
			fmt.Print(engine.RenderDominos(g.History))
		}

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
	fmt.Println("\nYour Card:")

	isValidIdx := make(map[int]bool)
	for _, vm := range validMoves {
		isValidIdx[vm.HandIndex] = true
	}

	for i := range hand {
		marker := " "
		if isValidIdx[i] {
			marker = "v"
		}

		fmt.Printf("    %s     ", marker)
	}
	fmt.Println()

	for i := range hand {
		fmt.Printf("    %d     ", i)
	}
	fmt.Println()

	fmt.Print(engine.RenderDominos(hand))

	for {
		fmt.Print("\nChoose number under the 'v': ")
		inputStr, _ := g.Reader.ReadString('\n')
		inputStr = strings.TrimSpace(inputStr)

		choice, err := strconv.Atoi(inputStr)

		if err == nil && isValidIdx[choice] {

			var matchedMoves []player.PlayableMove

			for _, vm := range validMoves {
				if vm.HandIndex == choice {
					matchedMoves = append(matchedMoves, vm)
				}
			}

			if len(matchedMoves) == 1 {
				temp := matchedMoves[0]
				return &temp
			}

			for {
				fmt.Print("Play on [1] Left / [2] Right: ")
				sideStr, _ := g.Reader.ReadString('\n')
				sideStr = strings.TrimSpace(sideStr)

				if sideStr == "1" {
					for _, vm := range matchedMoves {
						if vm.Side == "left" {
							temp := vm
							return &temp
						}
					}
				} else if sideStr == "2" {
					for _, vm := range matchedMoves {
						if vm.Side == "right" {
							temp := vm
							return &temp
						}
					}
				}

				fmt.Println("INVALID INPUT")
			}
		}
		fmt.Println("INVALID INPUT")
	}
}

func (g *GameState) executeMove(pIdx int, move *player.PlayableMove) {
	hand := g.Hands[pIdx]
	playedCard := hand[move.HandIndex]

	g.Hands[pIdx] = append(hand[:move.HandIndex], hand[move.HandIndex+1:]...)

	if pIdx != 0 {
		fmt.Printf("\n%s Play: 	         Remain: %d\n%s\n", g.Names[pIdx], len(g.Hands[pIdx]), playedCard.Card())
	}

	if g.BoardLeft == -1 {
		g.BoardLeft = playedCard.Left
		g.BoardRight = playedCard.Right
		g.History = append(g.History, playedCard)
	} else if move.Side == "left" {
		g.BoardLeft = move.LeaveNumber

		if playedCard.Left != move.LeaveNumber {
			playedCard = engine.Domino{Left: playedCard.Right, Right: playedCard.Left}
		}

		g.History = append([]engine.Domino{playedCard}, g.History...)
	} else {
		g.BoardRight = move.LeaveNumber

		if playedCard.Right != move.LeaveNumber {
			playedCard = engine.Domino{Left: playedCard.Right, Right: playedCard.Left}
		}

		g.History = append(g.History, playedCard)
	}
}

func (g *GameState) countGappleWinner() {
	fmt.Printf("\n===== COUNT REMAINING POINT ======\n")

	minPoints := 99
	winnerName := ""
	isDraw := false

	for i, hand := range g.Hands {
		total := 0

		for _, card := range hand {
			total += card.Points()
		}

		fmt.Printf("%s: %d point (%d card(s) left)\n", g.Names[i], total, len(hand))
		fmt.Print(engine.RenderDominos(hand))

		if total < minPoints {
			minPoints = total
			winnerName = g.Names[i]
			isDraw = false
		} else if total == minPoints {
			isDraw = true
		}
	}

	if isDraw {
		fmt.Print("\n===== DRAW =====\n")
	} else {
		fmt.Printf("\n%s WIN WITH %d POINT(s) \n", winnerName, minPoints)
	}
}
