package crossword

import (
	"bytes"
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
	"log"
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

func ResolveOptions(opts ...RenderOption) *Options {
	opt := &Options{}
	for _, v := range opts {
		v(opt)
	}
	return opt
}

type Options struct {
	solveAll bool
}

type RenderOption func(opts *Options)

func WithAllSolved() RenderOption {
	return func(opts *Options) {
		opts.solveAll = true
	}
}

func RenderText(cw *Crossword, opts ...RenderOption) string {
	options := ResolveOptions(opts...)

	out := &bytes.Buffer{}
	for y := range cw.Grid {
		for x := range cw.Grid[y] {
			if cw.Grid[y][x] == emptyCell {
				fmt.Fprintf(out, "#")
			} else {
				var solved bool
				for _, v := range cw.CellPlacements(x, y) {
					if v.Solved {
						solved = true
					}
				}
				if solved || options.solveAll {
					fmt.Fprintf(out, "%s", string(cw.Grid[y][x]))
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
	options := ResolveOptions(opts...)

	cellWidth := float64(width / len(c.Grid))
	cellHeight := float64(height / len(c.Grid))

	dc := gg.NewContext(width, height)
	dc.SetRGB(0, 0, 0)
	dc.Clear()

	for gridY := 0; gridY < len(c.Grid); gridY++ {
		for gridX, cell := range c.Grid[gridY] {

			dc.DrawRectangle(float64(gridX)*cellWidth, float64(gridY)*cellHeight, cellWidth, cellHeight)

			if cell != emptyCell {
				dc.SetRGB(1, 1, 1)
				dc.FillPreserve()
				dc.SetRGB(1, 0, 0)

				var solved bool
				placements := c.CellPlacements(gridX, gridY)
				if placements != nil {
					dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 12}))
					offset := 0.0
					for _, pl := range placements {
						if pl.X == gridX && pl.Y == gridY {
							// draw the word start identifier
							dc.DrawString(pl.ClueID(), float64(gridX)*cellWidth, float64(gridY)*cellHeight+12+offset)
							offset = cellHeight - 16
						}
						if pl.Solved {
							solved = true
						}
					}
				}

				dc.SetRGB(0, 0, 0)
				if solved || options.solveAll {
					dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 24}))
					dc.DrawStringAnchored(
						strings.ToUpper(string(cell)),
						float64(gridX)*cellWidth+cellWidth/2,
						float64(gridY)*cellHeight+cellHeight/2,
						0.5,
						0.5,
					)
				}

				dc.SetLineWidth(0.3)
				dc.Stroke()
			} else {
				dc.SetRGB(1, 1, 1)
				dc.SetLineWidth(0.3)
				dc.Stroke()
			}
		}
	}

	return dc, nil
}
