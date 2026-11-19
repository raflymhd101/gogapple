package engine

import (
	"fmt"
	"math/rand/v2"
	"strings"
)

type Domino struct {
	Left  int
	Right int
}

func (d Domino) IsBalak() bool {
	return d.Left == d.Right
}

func (d Domino) Points() int {
	return d.Left + d.Right
}

var points = map[int][]string{
	0: {"       ", "       ", "       "},
	1: {"       ",
		"   \033[31m●\033[0m   ",
		"       "},
	2: {" ●     ", "       ", "     ● "},
	3: {" ●     ", "   ●   ", "     ● "},
	4: {" \033[31m●\033[0m   \033[31m●\033[0m ",
		"       ",
		" \033[31m●\033[0m   \033[31m●\033[0m "},
	5: {" ●   ● ", "   ●   ", " ●   ● "},
	6: {" ●   ● ", " ●   ● ", " ●   ● "},
}

func RenderDominos(dominos []Domino) string {
	chunkSize := 15
	var result strings.Builder

	for i := 0; i < len(dominos); i += chunkSize {
		end := i + chunkSize

		if end > len(dominos) {
			end = len(dominos)
		}

		chunk := dominos[i:end]

		lines := make([]string, 9)
		for _, d := range chunk {
			top := points[d.Left]
			bot := points[d.Right]

			lines[0] += "┌───────┐ "
			lines[1] += fmt.Sprintf("│%s│ ", top[0])
			lines[2] += fmt.Sprintf("│%s│ ", top[1])
			lines[3] += fmt.Sprintf("│%s│ ", top[2])
			lines[4] += "├───────┤ "
			lines[5] += fmt.Sprintf("│%s│ ", bot[0])
			lines[6] += fmt.Sprintf("│%s│ ", bot[1])
			lines[7] += fmt.Sprintf("│%s│ ", bot[2])
			lines[8] += "└───────┘ "
		}
		result.WriteString(strings.Join(lines, "\n"))
		result.WriteString("\n\n")
	}

	return result.String()
}

func (d Domino) Card() string {
	return RenderDominos([]Domino{d})
}

func GenerateCard() []Domino {
	var board []Domino

	for i := 0; i <= 6; i++ {
		for j := i; j <= 6; j++ {
			board = append(board, Domino{Left: i, Right: j})
		}
	}

	rand.Shuffle(len(board), func(i, j int) {
		board[i], board[j] = board[j], board[i]
	})

	return board
}
