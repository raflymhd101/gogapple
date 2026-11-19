package player

import (
	"gogapple/internal/engine"
	"math/rand/v2"
)

type BotLevel int

const (
	LevelEasy   BotLevel = 1
	LevelMedium BotLevel = 2
	LevelHard   BotLevel = 3
)

type PlayableMove struct {
	HandIndex   int
	Side        string
	LeaveNumber int
}

func CalculateBotMove(hand []engine.Domino, boardLeft, boardRight int, history []engine.Domino, level BotLevel) *PlayableMove {
	moves := GetPlayableMoves(hand, boardLeft, boardRight)

	if len(moves) == 0 {
		return nil
	}

	switch level {
	case LevelEasy:
		return playEasy(moves)
	case LevelMedium:
		return playMedium(hand, moves)
	case LevelHard:
		return playHard(hand, moves, history)
	default:
		return playEasy(moves)
	}
}

func playEasy(moves []PlayableMove) *PlayableMove {
	idx := rand.IntN(len(moves))

	return &moves[idx]
}

func playMedium(hand []engine.Domino, moves []PlayableMove) *PlayableMove {

	bestIdx := 0
	maxPoints := -1

	for i, move := range moves {
		card := hand[move.HandIndex]

		if card.Points() > maxPoints {
			maxPoints = card.Points()
			bestIdx = i
		}
	}

	return &moves[bestIdx]
}

func playHard(hand []engine.Domino, moves []PlayableMove, history []engine.Domino) *PlayableMove {
	freq := countPointFrequency(history)

	bestIdx := 0
	highestScore := -1

	for i, move := range moves {
		card := hand[move.HandIndex]
		score := card.Points() * 2

		if card.IsBalak() {
			score += 100
		}

		score += freq[move.LeaveNumber] * 10

		if score > highestScore {
			highestScore = score
			bestIdx = i
		}
	}

	return &moves[bestIdx]
}

func GetPlayableMoves(hand []engine.Domino, boardLeft, boardRight int) []PlayableMove {
	var moves []PlayableMove

	for i, card := range hand {
		if boardLeft == -1 {
			moves = append(moves, PlayableMove{HandIndex: i, Side: "left", LeaveNumber: card.Right})
			continue
		}

		if card.Left == boardLeft {
			moves = append(moves, PlayableMove{HandIndex: i, Side: "left", LeaveNumber: card.Right})
		} else if card.Right == boardLeft {
			moves = append(moves, PlayableMove{HandIndex: i, Side: "left", LeaveNumber: card.Left})
		}

		if card.Left == boardRight {
			moves = append(moves, PlayableMove{HandIndex: i, Side: "right", LeaveNumber: card.Right})
		} else if card.Right == boardRight {
			moves = append(moves, PlayableMove{HandIndex: i, Side: "right", LeaveNumber: card.Left})
		}
	}

	return moves
}

func countPointFrequency(history []engine.Domino) map[int]int {
	freq := make(map[int]int)

	for _, card := range history {
		freq[card.Left]++

		if !card.IsBalak() {
			freq[card.Right]++
		}
	}

	return freq
}
