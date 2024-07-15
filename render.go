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

func RenderText(cw *Crossword) string {
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
				if solved {
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

func RenderPNG(c *Crossword, width, height int) (*gg.Context, error) {

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
							dc.DrawString(pl.ID(), float64(gridX)*cellWidth, float64(gridY)*cellHeight+12+offset)
							offset = cellHeight - 16
						}
						if pl.Solved {
							solved = true
						}
					}
				}

				dc.SetRGB(0, 0, 0)
				if solved {
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
