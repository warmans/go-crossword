package crossword

import (
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

// WordsFromCSV creates a word list from a CSV with 2 columns (word, clue)
func WordsFromCSV(f io.Reader) ([]Word, error) {
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}
	words := []Word{}
	for _, r := range rows {
		if len(r) != 2 {
			return nil, fmt.Errorf("csv should have exactly 2 columns (word, clue)")
		}
		words = append(words, Word{Word: strings.TrimSpace(r[0]), Clue: strings.TrimSpace(r[1])})
	}
	return words, nil
}
