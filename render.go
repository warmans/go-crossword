package crossword

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"slices"
	"strings"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
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
		wordFontSizePcnt:    0.5,
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
	wordFontSizePcnt    float64
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

// WithWordFontSizePcnt sets the font size as a proportion of the containing square on the board.
func WithWordFontSizePcnt(size float64) RenderOption {
	return func(opts *renderOpts) {
		opts.wordFontSizePcnt = size
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

	if options.solveRandom {
		for k := range c.Words {
			if rand.Float64() < 0.5 {
				c.Words[k].Solved = true
			}
		}
	}

	B := float64(options.borderWidth)

	var gridWidth float64
	if !options.renderClues {
		gridWidth = float64(width) - 2*B
	} else {
		gridWidth = (float64(width) - 3*B) / 2
	}

	// ensure the grid is square and fits in the vertical space
	if gridWidth > float64(height)-2*B {
		gridWidth = float64(height) - 2*B
	}

	cellWidth := gridWidth / float64(len(c.Grid))
	cellHeight := cellWidth
	cellOffset := B

	dc := gg.NewContext(width, height)

	clueFontSize := 25.0
	checkboxSize := 10.0
	checkboxSpace := 15.0
	var leftPos, maxClueWidth float64
	if options.renderClues {
		leftPos = cellOffset + gridWidth + B
		maxClueWidth = float64(width) - leftPos - B

		// try to find a font size that fits.
		for clueFontSize > 4 {
			dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: clueFontSize}))
			checkboxSize = dc.FontHeight() * 0.8
			checkboxSpace = checkboxSize + 5
			if measureCluesHeight(c, dc, clueFontSize, maxClueWidth-checkboxSpace) <= float64(height)-2*B {
				break
			}
			clueFontSize -= 0.5
		}
	}

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
					clueIDFontSize := cellHeight * 0.25
					dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: clueIDFontSize}))
					offset := 0.0
					for _, pl := range placements {
						if pl.X == gridX && pl.Y == gridY {
							// draw the word start identifier
							dc.DrawString(pl.ClueID(), cellOffset+float64(gridX)*cellWidth, cellOffset+float64(gridY)*cellHeight+clueIDFontSize+offset)
							offset = cellHeight - (clueIDFontSize * 1.4)
						}
						if pl.Solved || slices.Contains(pl.Word.CharacterHints, cell.CharIdx) {
							solved = true
						}
					}
				}

				dc.SetColor(options.wordColor)
				if solved || options.solveAll {
					dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: cellHeight * options.wordFontSizePcnt}))
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
		dc.SetColor(options.clueColor)
		dc.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: clueFontSize}))
		dc.SetLineWidth(0.3)

		offset := B + clueFontSize

		// DOWN
		dc.DrawStringAnchored("DOWN", leftPos, offset, 0, 0)
		offset += clueFontSize
		for _, w := range c.Words {
			if w.Vertical {
				s := fmt.Sprintf("%s: %s [%s]", w.ClueID(), w.Word.Clue, w.Word.LetterCountStr())
				height := measureWrappedHeight(dc, s, maxClueWidth-checkboxSpace)

				dc.DrawRectangle(leftPos, offset+(height/2)-(checkboxSize/2), checkboxSize, checkboxSize)
				dc.StrokePreserve()
				if w.Solved {
					dc.Fill()
				}
				dc.ClearPath()
				drawStringWrapped(dc, s, leftPos+checkboxSpace, offset, maxClueWidth-checkboxSpace)
				offset += height + 5
			}
		}

		offset += 32
		dc.DrawStringAnchored("ACROSS", leftPos, offset, 0, 0)
		offset += clueFontSize
		for _, w := range c.Words {
			if !w.Vertical {
				s := fmt.Sprintf("%s: %s [%s]", w.ClueID(), w.Word.Clue, w.Word.LetterCountStr())
				height := measureWrappedHeight(dc, s, maxClueWidth-checkboxSpace)

				dc.DrawRectangle(leftPos, offset+(height/2)-(checkboxSize/2), checkboxSize, checkboxSize)
				dc.StrokePreserve()
				if w.Solved {
					dc.Fill()
				}
				dc.ClearPath()
				drawStringWrapped(dc, s, leftPos+checkboxSpace, offset, maxClueWidth-checkboxSpace)
				offset += height + 5
			}
		}
	}

	return dc, nil
}

func drawStringWrapped(dc *gg.Context, s string, x, y float64, maxWidth float64) {
	dc.DrawStringWrapped(s, x, y, 0, 0, maxWidth, 1.0, gg.AlignLeft)
}

func measureWrappedHeight(dc *gg.Context, s string, maxWidth float64) float64 {
	_, height := dc.MeasureMultilineString(strings.Join(dc.WordWrap(s, maxWidth), "\n"), 1.0)
	return height
}

func measureCluesHeight(c *Crossword, dc *gg.Context, fontSize float64, maxWidth float64) float64 {
	offset := fontSize // DOWN header
	offset += fontSize // DOWN header space
	for _, w := range c.Words {
		if w.Vertical {
			offset += measureWrappedHeight(dc, fmt.Sprintf("%s: %s [%s]", w.ClueID(), w.Word.Clue, w.Word.LetterCountStr()), maxWidth) + 5
		}
	}
	offset += 32       // middle space
	offset += fontSize // ACROSS header
	offset += fontSize // ACROSS header space
	for _, w := range c.Words {
		if !w.Vertical {
			offset += measureWrappedHeight(dc, fmt.Sprintf("%s: %s [%s]", w.ClueID(), w.Word.Clue, w.Word.LetterCountStr()), maxWidth) + 5
		}
	}
	return offset
}
