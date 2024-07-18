package crossword

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestWordsFromCSV(t *testing.T) {
	tests := []struct {
		name    string
		f       io.Reader
		want    []Word
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "empty reader returns empty words array",
			f:       strings.NewReader(""),
			want:    []Word{},
			wantErr: assert.NoError,
		},
		{
			name: "valid csv returns words",
			f:    strings.NewReader("foo, foo clue\nbar, bar clue"),
			want: []Word{
				{Word: "foo", Clue: "foo clue"},
				{Word: "bar", Clue: "bar clue"},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "invalid number of fields",
			f:       strings.NewReader("foo\nbar, bar clue"),
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WordsFromCSV(tt.f)
			if !tt.wantErr(t, err, fmt.Sprintf("WordsFromCSV(%v)", tt.f)) {
				return
			}
			assert.Equalf(t, tt.want, got, "WordsFromCSV(%v)", tt.f)
		})
	}
}
