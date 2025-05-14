package crossword

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
	"math/rand"
	"slices"
	"strings"
)

var font *truetype.Font

func init() {
	var err error
	font, err = truetype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err)
	}
}

func resolveRenderOptions(opts ...RenderOption) *renderOpts {
	opt := &renderOpts{
		backgroundColor:     color.Black,
		wordBackgroundColor: color.White,
		wordColor:           color.Black,
		labelColor:          color.RGBA{R: 200, G: 10, B: 10, A: 255},
		clueColor:           color.White,
	}
	for _, v := range opts {
		v(opt)
	}
	return opt
}

type renderOpts struct {
	solveAll            bool
	solveRandom         bool
	borderWidth         int
	backgroundColor     color.Color
	wordBackgroundColor color.Color
	wordColor           color.Color
	labelColor          color.Color
	clueColor           color.Color
	renderClues         bool
}

type RenderOption func(opts *renderOpts)

func WithClues(clues bool) RenderOption {
	return func(opts *renderOpts) {
		opts.renderClues = clues
	}
}

func WithClueColor(color color.Color) RenderOption {
	return func(opts *renderOpts) {
		opts.clueColor = color
	}
}

func WithAllSolved(solveAll bool) RenderOption {
	return func(opts *renderOpts) {
		opts.solveAll = solveAll
	}
}

func WithRandomSolved() RenderOption {
	return func(opts *renderOpts) {
		opts.solveRandom = true
	}
}

func WithBorder(width int) RenderOption {
	return func(opts *renderOpts) {
		opts.borderWidth = width
	}
}

func WithBackgroundColor(cl color.Color) RenderOption {
	return func(opts *renderOpts) {
		opts.backgroundColor = cl
	}
}

func WithWordBackgroundColor(cl color.Color) RenderOption {
	return func(opts *renderOpts) {
		opts.wordBackgroundColor = cl
	}
}

func WithWordColor(cl color.Color) RenderOption {
	return func(opts *renderOpts) {
		opts.wordColor = cl
	}
}

func WithLabelColor(cl color.Color) RenderOption {
	return func(opts *renderOpts) {
		opts.labelColor = cl
	}
}

func RenderText(cw *Crossword, opts ...RenderOption) string {
	options := resolveRenderOptions(opts...)

	out := &bytes.Buffer{}
	for y := range cw.Grid {
		for x := range cw.Grid[y] {
			if cw.Grid[y][x].Empty() {
				fmt.Fprintf(out, "#")
			} else {
				var solved bool
				for _, v := range cw.CellPlacements(x, y) {
					if v.Solved || slices.Contains(v.Word.CharacterHints, cw.Grid[y][x].CharIdx) {
						solved = true
					}
				}
				if solved || options.solveAll {
					fmt.Fprintf(out, "%s", cw.Grid[y][x])
				} else {
					fmt.Fprintf(out, "?")
				}
			}
		}
		fmt.Fprintf(out, "\n")
	}
	return out.String()
}

func RenderPNG(c *Crossword, width, height int, opts ...RenderOption) (*gg.Context, error) {
	options := resolveRenderOptions(opts...)

	if options.solveAll {
		c.Solve()
	}
	if options.solveRandom {
		for k := range c.Words {
			if rand.Float64() < 0.5 {
				c.Words[k].Solved = true
			}
		}
	}

	var gridWidth, gridHeight int
	if !options.renderClues {
		gridWidth = width - options.borderWidth
		gridHeight = height - options.borderWidth
	} else {
		gridWidth = (width - options.borderWidth) / 2
		gridHeight = height - options.borderWidth
	}

	cellWidth := float64(gridWidth / len(c.Grid))
	cellHeight := float64(gridHeight / len(c.Grid))
	cellOffset := 0.0
	if options.borderWidth > 0 {
		cellOffset = float64(options.borderWidth) / 2
	}

	dc := gg.NewContext(width, height)
	dc.SetColor(options.backgroundColor)
	dc.Clear()

	for gridY := 0; gridY < len(c.Grid); gridY++ {
		for gridX, cell := range c.Grid[gridY] {

			dc.DrawRectangle(cellOffset+(float64(gridX)*cellWidth), cellOffset+(float64(gridY)*cellHeight), cellWidth, cellHeight)

			if !cell.Empty() {
				dc.SetColor(options.wordBackgroundColor)
				dc.FillPreserve()

				dc.SetColor(options.labelColor)
				var solved bool
				placements := c.CellPlacements(gridX, gridY)
				if placements != nil {
					dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 10}))
					offset := 0.0
					for _, pl := range placements {
						if pl.X == gridX && pl.Y == gridY {
							// draw the word start identifier
							dc.DrawString(pl.ClueID(), cellOffset+float64(gridX)*cellWidth, cellOffset+float64(gridY)*cellHeight+12+offset)
							offset = cellHeight - 16
						}
						if pl.Solved || slices.Contains(pl.Word.CharacterHints, cell.CharIdx) {
							solved = true
						}
					}
				}

				dc.SetColor(options.wordColor)
				if solved {
					dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 20}))
					dc.DrawStringAnchored(
						strings.ToUpper(cell.String()),
						cellOffset+float64(gridX)*cellWidth+cellWidth/2,
						cellOffset+float64(gridY)*cellHeight+cellHeight/2,
						0.5,
						0.5,
					)
				}

				dc.SetLineWidth(0.3)
				dc.Stroke()
			} else {
				dc.SetLineWidth(0)
				dc.Stroke()
			}
		}
	}

	if options.renderClues {
		leftPos := float64(gridWidth) + (float64(options.borderWidth) * 2)
		maxClueWidth := float64(gridWidth) - (float64(options.borderWidth) * 2)
		checkboxSpace := 15.0

		dc.SetColor(options.clueColor)
		dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 12}))
		dc.SetLineWidth(0.3)

		offset := float64(options.borderWidth) / 2

		// DOWN
		dc.DrawString("DOWN", leftPos, offset)
		offset += 10
		for _, w := range c.Words {

			if w.Vertical {
				dc.DrawRectangle(leftPos, offset+2, 10, 10)
				dc.StrokePreserve()
				if w.Solved {
					dc.Fill()
				}
				dc.ClearPath()
				offset += drawStringWrapped(dc, fmt.Sprintf("%s: %s", w.ClueID(), w.Word.Clue), leftPos+checkboxSpace, offset, maxClueWidth)
			}
		}

		offset += 32
		dc.DrawString("ACROSS", leftPos, offset)
		offset += 10
		for _, w := range c.Words {
			if !w.Vertical {
				dc.DrawRectangle(leftPos, offset+2, 10, 10)
				dc.StrokePreserve()
				if w.Solved {
					dc.Fill()
				}
				dc.ClearPath()
				offset += drawStringWrapped(dc, fmt.Sprintf("%s: %s", w.ClueID(), w.Word.Clue), leftPos+checkboxSpace, offset, maxClueWidth)
			}
		}

	}

	return dc, nil
}

func drawStringWrapped(dc *gg.Context, s string, x, y float64, maxWidth float64) float64 {
	var lineSpacing = 1.0
	_, height := dc.MeasureMultilineString(strings.Join(dc.WordWrap(s, maxWidth), "\n"), lineSpacing)
	dc.DrawStringWrapped(s, x, y, 0, 0, maxWidth, lineSpacing, gg.AlignLeft)
	return height + 5 // add some extra space
}
