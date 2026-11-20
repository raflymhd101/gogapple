package ui

import (
	"bufio"
	"fmt"
	"gogapple/internal/engine"
	"gogapple/internal/player"
	"strconv"
	"strings"
)

func printBoard(boardLeft, boardRight int, history []engine.Domino) {
	if boardLeft == -1 {
		fmt.Printf("\n======= TABLE STILL EMPTY ========\n")
	} else {
		fmt.Printf("\n=== CURRENT CARD === (Left: %d | Right: %d)\n", boardLeft, boardRight)
		fmt.Print(engine.RenderDominos(history))
	}
}

func PromptHumanMove(reader *bufio.Reader, hand []engine.Domino, validMoves []player.PlayableMove) *player.PlayableMove {
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
		inputStr, _ := reader.ReadString('\n')
		choice, err := strconv.Atoi(strings.TrimSpace(inputStr))

		if err == nil && isValidIdx[choice] {

			var matchedMoves []player.PlayableMove

			for _, vm := range validMoves {
				if vm.HandIndex == choice {
					matchedMoves = append(matchedMoves, vm)
				}
			}

			if len(matchedMoves) == 1 {
				return &matchedMoves[0]
			}

			promptSide(reader, matchedMoves)
		}
		fmt.Println("INVALID INPUT")
	}
}

func promptSide(reader *bufio.Reader, matchedMoves []player.PlayableMove) *player.PlayableMove {
	for {
		fmt.Print("Play on [1] Left / [2] Right: ")
		sideStr, _ := reader.ReadString('\n')
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
