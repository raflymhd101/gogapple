package player

import (
	"gogapple/internal/engine"
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
		return playMedium(moves)
	case LevelHard:
		return playHard(moves)
	default:
		return playEasy(moves)
	}
}

func playEasy(moves []PlayableMove) *PlayableMove {
	return &moves[0]
}

func playMedium(moves []PlayableMove) *PlayableMove {
	return &moves[len(moves)-1]
}

func playHard(moves []PlayableMove) *PlayableMove {
	return &moves[0]
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
