package crossword

import (
	"fmt"
	"slices"
)

type Cell struct {
	Char    rune
	CharIdx int
}

func (c Cell) String() string {
	return string(c.Char)
}

func (c Cell) Empty() bool {
	return c.Char == rune(0)
}

func NewGrid(size int) Grid {
	grid := make(Grid, size)
	for y := range size {
		grid[y] = make([]Cell, size)
	}
	return grid
}

type Grid [][]Cell

type Placement struct {
	ID       int
	Word     Word
	X        int
	Y        int
	Vertical bool
	// Solved reveals the characters of the word when it's rendered.
	Solved bool
}

func (p Placement) ClueID() string {
	label := fmt.Sprintf("%d", p.ID)
	if p.Word.Label != nil {
		label = *p.Word.Label
	}
	if p.Vertical {
		return fmt.Sprintf("D%s", label)
	}
	return fmt.Sprintf("A%s", label)
}

type Word struct {
	Word  string
	Clue  string
	Label *string

	// CharacterHints allows subset of characters to be revealed (e.g. []int{0} would reveal
	// the first char of a word by default)
	CharacterHints []int
}

type Crossword struct {
	Grid       Grid
	Words      []Placement
	TotalScore int
}

func (cw *Crossword) Solve() {
	for k := range cw.Words {
		cw.Words[k].Solved = true
	}
}

func (cw *Crossword) CellPlacements(cellX, cellY int) []Placement {
	var placements []Placement
	for _, pl := range cw.Words {
		if pl.Vertical {
			if pl.X == cellX && cellY >= pl.Y && cellY < pl.Y+len(pl.Word.Word) {
				placements = append(placements, pl)
			}
		} else {
			if pl.Y == cellY && cellX >= pl.X && cellX < pl.X+len(pl.Word.Word) {
				placements = append(placements, pl)
			}
		}
	}
	slices.SortFunc(placements, func(a, b Placement) int {
		if a.Vertical {
			return 1
		}
		return -1
	})
	return placements
}
