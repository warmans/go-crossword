package crossword

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
)

const emptyCell = Cell(0)

type Cell rune

func NewGrid(size int) Grid {
	grid := make(Grid, size)
	for y := range size {
		grid[y] = make([]Cell, size)
	}
	return grid
}

type Grid [][]Cell

type Placement struct {
	id       int
	Word     Word
	X        int
	Y        int
	Vertical bool
	Solved   bool
}

func (p Placement) ID() string {
	if p.Vertical {
		return fmt.Sprintf("D%d", p.id)
	}
	return fmt.Sprintf("A%d", p.id)
}

type Word struct {
	Word string
	Clue string
}

type Crossword struct {
	Grid  Grid
	Words []Placement
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

func Generate(gridSize int, words []Word, attempts int) *Crossword {
	return NewGenerator(gridSize).Generate(words, attempts)
}

func NewGenerator(gridSize int) *Generator {
	return &Generator{gridSize: gridSize, grid: NewGrid(gridSize)}
}

type Generator struct {
	gridSize    int
	grid        Grid
	placedWords []Placement
}

func (g *Generator) Generate(words []Word, attempts int) *Crossword {

	for k := range words {
		words[k].Word = strings.ToUpper(words[k].Word)
	}

	var bestCrossword *Crossword
	for range attempts {
		if attempts == 0 {
			// first attempt sort words by length
			slices.SortFunc(words, func(a, b Word) int {
				if len(a.Word) > len(b.Word) {
					return -1
				}
				return 1
			})
		} else {
			// remaining attempts should randomize the words instead
			slices.SortFunc(words, func(a, b Word) int {
				if rand.Float64() > 0.5 {
					return -1
				}
				return 1
			})
		}
		for startWord := range len(words) - 1 {
			// place the first word
			g.placeWord(Placement{
				id:       1,
				Word:     words[startWord],
				X:        0,
				Y:        0,
				Vertical: false,
			})
			for k, word := range words {
				if k == startWord {
					continue
				}
				if slices.IndexFunc(g.placedWords, func(placement Placement) bool {
					return word == placement.Word
				}) > -1 {
					continue
				}
				placements := g.suggestPlacements(word)
				if placements == nil {
					continue
				}
				var bestPlacement *Placement
				var bestScore int
				for k, pl := range placements {
					if score := g.scorePlacement(pl); bestPlacement == nil || score > bestScore {
						bestPlacement = &placements[k]
						bestScore = score
					}
				}
				if bestPlacement == nil || bestScore < 2 {
					continue
				}

				g.placeWord(*bestPlacement)
			}

			if bestCrossword == nil || len(g.placedWords) > len(bestCrossword.Words) {
				bestCrossword = &Crossword{Words: g.placedWords, Grid: g.grid}
			}
			*g = *NewGenerator(g.gridSize)
		}
		if len(words) == len(bestCrossword.Words) {
			return bestCrossword
		}
	}

	return bestCrossword
}

func (g *Generator) placeWord(placement Placement) {
	for c := range len(placement.Word.Word) {
		// don't bother checking if the word fits since this should already happen
		// in suggestPlacements
		if !placement.Vertical {
			g.grid[placement.Y][placement.X+c] = Cell(placement.Word.Word[c])
		} else {
			g.grid[placement.Y+c][placement.X] = Cell(placement.Word.Word[c])
		}
	}
	placement.id = len(g.placedWords) + 1
	g.placedWords = append(g.placedWords, placement)
}

func (g *Generator) suggestPlacements(word Word) []Placement {
	var placements []Placement
	for charIdx := range len(word.Word) {
		for y := range g.gridSize {
			for x := range g.gridSize {
				// word intersects existing cell
				if g.grid[y][x] == Cell(word.Word[charIdx]) {
					// check vertical fit.
					{
						if y-charIdx >= 0 && y+(len(word.Word)-(charIdx+1)) < g.gridSize {
							placements = append(placements, Placement{
								Word:     word,
								X:        x,
								Y:        y - charIdx,
								Vertical: true,
							})
						}
					}
					// check horizontal fit.
					if x-charIdx >= 0 && x+(len(word.Word)-(charIdx+1)) < g.gridSize {
						placements = append(placements, Placement{
							Word: word,
							X:    x - charIdx,
							Y:    y,
						})
					}
				}
			}
		}
	}
	return placements
}

func (g *Generator) scorePlacement(pl Placement) int {
	score := 1
	// word overflows grid
	if (!pl.Vertical && pl.X+len(pl.Word.Word)-1 > g.gridSize) || (pl.Vertical && pl.Y+len(pl.Word.Word)-1 > g.gridSize) {
		return 0
	}
	// horizontal checking
	if !pl.Vertical {
		for charIdx := range len(pl.Word.Word) {

			// if the word doesn't start at the edge of the board...
			if charIdx == 0 && pl.X > 0 {
				// check preceding cell for collision
				if g.grid[pl.Y][pl.X-1] != emptyCell {
					return 0
				}
			}
			// if the word doesn't end at the edge of the board...
			if charIdx == len(pl.Word.Word)-1 && (pl.X+len(pl.Word.Word)) < len(g.grid[pl.Y]) {
				// check following cell for collision
				if g.grid[pl.Y][pl.X+len(pl.Word.Word)] != emptyCell {
					return 0
				}
			}

			// increase score for any valid overlaps
			nextCharInGrid := g.grid[pl.Y][pl.X+charIdx]
			if Cell(pl.Word.Word[charIdx]) == nextCharInGrid {
				score += 1
			} else if nextCharInGrid != emptyCell {
				return 0
			} else {
				// check the word has space above and below if it is not intersecting a vertical word
				if pl.Y > 0 && g.grid[pl.Y-1][pl.X+charIdx] != emptyCell {
					return 0
				}
				if pl.Y < g.gridSize-1 && g.grid[pl.Y+1][pl.X+charIdx] != emptyCell {
					return 0
				}
			}
			// check the next cell to the last char
			if charIdx == (len(pl.Word.Word) - 1) {
				nextCellIdx := pl.X + charIdx + 1
				if nextCellIdx < len(g.grid[pl.Y]) && g.grid[pl.Y][nextCellIdx] != emptyCell {
					return 0
				}
			}
		}
	} else {
		for charIdx := range len(pl.Word.Word) {
			// if the word doesn't start at the top of the board...
			if charIdx == 0 && pl.Y > 0 {
				// check preceding cell for collision
				if g.grid[pl.Y-1][pl.X] != emptyCell {
					return 0
				}
			}
			// if the word doesn't end at the edge of the board...
			if charIdx == len(pl.Word.Word)-1 && (pl.Y+len(pl.Word.Word)-1) < len(g.grid[pl.X])-1 {
				// check following cell for collision
				if g.grid[pl.Y+len(pl.Word.Word)][pl.X] != emptyCell {
					return 0
				}
			}

			// increase score for any valid overlaps
			nextCharInGrid := g.grid[pl.Y+charIdx][pl.X]
			if Cell(pl.Word.Word[charIdx]) == nextCharInGrid {
				score += 1
			} else if nextCharInGrid != emptyCell {
				return 0
			} else {
				// check the word has space to the left and right if it is not intersecting a vertical word

				// left
				if pl.X > 0 && g.grid[pl.Y+charIdx][pl.X-1] != emptyCell {
					return 0
				}

				// right
				if pl.X < g.gridSize-1 && g.grid[pl.Y+charIdx][pl.X+1] != emptyCell {
					return 0
				}
			}
		}
	}
	return score
}
