package crossword

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"image/color"
	"log"
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
	}
	for _, v := range opts {
		v(opt)
	}
	return opt
}

type renderOpts struct {
	solveAll            bool
	borderWidth         int
	backgroundColor     color.Color
	wordBackgroundColor color.Color
	wordColor           color.Color
	labelColor          color.Color
}

type RenderOption func(opts *renderOpts)

func WithAllSolved(solveAll bool) RenderOption {
	return func(opts *renderOpts) {
		opts.solveAll = solveAll
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

	gridWidth := width - options.borderWidth
	gridHeight := height - options.borderWidth

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
					dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 12}))
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
				if solved || options.solveAll {
					dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 24}))
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
				//	dc.SetRGB(1, 1, 1)
				dc.SetLineWidth(0)
				dc.Stroke()
			}
		}
	}

	return dc, nil
}
