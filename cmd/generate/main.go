package main

import (
	"encoding/json"
	"fmt"
	"github.com/warmans/go-crossword"
	"os"
	"path"
	"strings"
)

func main() {

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

	cw := crossword.Generate(25, words, 50)
	cw.Solve()
	fmt.Print(crossword.RenderText(cw))
	fmt.Printf("INPUT WORDS: %d OUTPUT WORDS: %d\n", len(words), len(cw.Words))

	canvas, err := crossword.RenderPNG(cw, 1200, 1200)
	if err != nil {
		panic(err.Error())
	}
	if err := canvas.SavePNG(strings.TrimSuffix(path.Base(os.Args[1]), ".json") + ".png"); err != nil {
		panic(err.Error())
	}
}
