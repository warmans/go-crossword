package main

import (
	"encoding/json"
	"fmt"
	"github.com/warmans/go-crossword"
	"image/color"
	"os"
	"strconv"
)

func main() {

	solveAll := os.Getenv("SOLVE_ALL") == "true"
	attempts := 10
	if intVal, err := strconv.ParseInt(os.Getenv("ATTEMPTS"), 10, 64); err == nil {
		attempts = int(intVal)
	}
	fmt.Printf("Running %d attempts\n", attempts)

	if len(os.Args) != 2 {
		fmt.Println("Expected first argument to contain word file path")
		os.Exit(1)
	}
	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	var words []crossword.Word
	if err := json.NewDecoder(f).Decode(&words); err != nil {
		panic(err.Error())
	}

	cw := crossword.Generate(25, words, attempts, crossword.WithAllAttempts(true))
	fmt.Print(crossword.RenderText(cw, crossword.WithAllSolved(solveAll)))
	fmt.Printf("INPUT WORDS: %d OUTPUT WORDS: %d TOTAL SCORE: %d\n", len(words), len(cw.Words), cw.TotalScore)

	canvas, err := crossword.RenderPNG(
		cw,
		1500,
		750,
		crossword.WithRandomSolved(),
		crossword.WithBorder(50),
		crossword.WithBackgroundColor(color.RGBA{R: 30, G: 30, B: 50, A: 255}),
		crossword.WithWordBackgroundColor(color.RGBA{R: 255, G: 255, B: 255, A: 255}),
		crossword.WithWordColor(color.RGBA{R: 10, G: 10, B: 10, A: 255}),
		crossword.WithLabelColor(color.RGBA{R: 200, G: 10, B: 10, A: 255}),
		crossword.WithClues(true),
	)
	if err != nil {
		panic(err.Error())
	}
	if err := canvas.SavePNG(fmt.Sprintf("example/simple/crossword-%d.png", attempts)); err != nil {
		panic(err.Error())
	}
}
