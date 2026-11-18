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

func formatSide(value int) string {
	if value == 0 {
		return "[      ]"
	}

	points := strings.Repeat("●", value)
	return fmt.Sprintf("[%-6s]", points)
}

func (d Domino) TopPoint() string {
	return formatSide(d.Left)
}

func (d Domino) BottomPoint() string {
	return formatSide(d.Right)
}

func (d Domino) Card() string {
	return fmt.Sprintf("%s\n %s", d.TopPoint(), d.BottomPoint())
}

func GenerateCard() []Domino {
	var board []Domino

	for i := 0; i <= 6; i++ {
		for j := 0; j <= 6; j++ {
			board = append(board, Domino{Left: i, Right: j})
		}
	}

	rand.Shuffle(len(board), func(i, j int) {
		board[i], board[j] = board[j], board[i]
	})

	return board
}
