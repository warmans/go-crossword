package crossword

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func word(chars string) Word {
	return Word{Word: chars}
}

func TestGenerator_suggestCoordinates(t *testing.T) {
	type fields struct {
		gridSize int
		grid     Grid
	}
	type args struct {
		word string
	}
	tests := []struct {
		name          string
		existingWords []Placement
		fields        fields
		args          args
		want          []Placement
	}{
		{
			name: "no placements if grid is empty",
			fields: fields{
				gridSize: 3,
				grid:     NewGrid(3),
			},
			args: args{
				word: "foo",
			},
			want: nil,
		}, {
			name: "intersecting horizontal word is valid placement",
			existingWords: []Placement{{
				Word: word("fud"),
				X:    0,
				Y:    0,
			}},
			fields: fields{
				gridSize: 3,
				grid:     NewGrid(3),
			},
			args: args{
				word: "foo",
			},
			want: []Placement{{
				Word:     word("foo"),
				X:        0,
				Y:        0,
				Vertical: true,
			}, {
				Word:     word("foo"),
				X:        0,
				Y:        0,
				Vertical: false,
			}},
		}, {
			name: "intersecting vertical word is valid placement",
			existingWords: []Placement{{
				Word:     word("fud"),
				X:        0,
				Y:        0,
				Vertical: true,
			}},
			fields: fields{
				gridSize: 3,
				grid:     NewGrid(3),
			},
			args: args{
				word: "foo",
			},
			want: []Placement{{
				Word:     word("foo"),
				X:        0,
				Y:        0,
				Vertical: true,
			}, {
				Word:     word("foo"),
				X:        0,
				Y:        0,
				Vertical: false,
			}},
		}, {
			name: "check last letter vertical intersection",
			existingWords: []Placement{{
				Word:     word("fud"),
				X:        0,
				Y:        0,
				Vertical: false,
			}},
			fields: fields{
				gridSize: 3,
				grid:     NewGrid(3),
			},
			args: args{
				word: "doo",
			},
			want: []Placement{{
				Word:     word("doo"),
				X:        2,
				Y:        0,
				Vertical: true,
			}},
		}, {
			name: "check last letter horizontal intersection",
			existingWords: []Placement{{
				Word:     word("fud"),
				X:        0,
				Y:        0,
				Vertical: true,
			}},
			fields: fields{
				gridSize: 3,
				grid:     NewGrid(3),
			},
			args: args{
				word: "doo",
			},
			want: []Placement{{
				Word:     word("doo"),
				X:        0,
				Y:        2,
				Vertical: false,
			}},
		}, {
			name: "midpoint intersection",
			existingWords: []Placement{{
				Word:     word("fop"),
				X:        0,
				Y:        1,
				Vertical: false,
			}},
			fields: fields{
				gridSize: 3,
				grid:     NewGrid(3),
			},
			args: args{
				word: "dox",
			},
			want: []Placement{{
				Word:     word("dox"),
				X:        1,
				Y:        0,
				Vertical: true,
			}, {
				Word:     word("dox"),
				X:        0,
				Y:        1,
				Vertical: false,
			}},
		}, {
			name: "no placement if word is too long",
			existingWords: []Placement{{
				Word:     word("foo"),
				X:        0,
				Y:        0,
				Vertical: true,
			}},
			fields: fields{
				gridSize: 3,
				grid:     NewGrid(3),
			},
			args: args{
				word: "food",
			},
			want: nil,
		}, {
			name: "no placement if word is too high",
			existingWords: []Placement{{
				Word:     word("foo"),
				X:        0,
				Y:        0,
				Vertical: false,
			}},
			fields: fields{
				gridSize: 3,
				grid:     NewGrid(3),
			},
			args: args{
				word: "food",
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				gridSize: tt.fields.gridSize,
				grid:     tt.fields.grid,
			}
			for _, v := range tt.existingWords {
				g.placeWord(v)
			}
			require.EqualValues(t, tt.want, g.suggestPlacements(word(tt.args.word)))
		})
	}
}

func TestGenerator_scorePlacement_horizontalWords(t *testing.T) {
	tests := []struct {
		name          string
		existingWords []Placement
		generator     *Generator
		placement     Placement
		want          int
	}{
		{
			name:      "score 1 if no X collisions",
			generator: NewGenerator(4),
			placement: Placement{
				Word:     word("foo"),
				X:        0,
				Y:        0,
				Vertical: false,
			},
			want: 1,
		},
		{
			name: "score 0 if preceding cell collides",
			existingWords: []Placement{{
				Word:     word("abcd"),
				X:        0,
				Y:        0,
				Vertical: false,
			}},
			generator: NewGenerator(4),
			placement: Placement{
				Word:     word("bcd"),
				X:        1,
				Y:        0,
				Vertical: false,
			},
			want: 0,
		},
		{
			name: "score 0 if following cell collides",
			existingWords: []Placement{{
				Word:     word("abcd"),
				X:        0,
				Y:        0,
				Vertical: false,
			}},
			generator: NewGenerator(4),
			placement: Placement{
				Word:     word("abc"),
				X:        1,
				Y:        0,
				Vertical: false,
			},
			want: 0,
		},
		{
			name: "score 2 if match one char",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: true,
			}},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: false,
			},
			want: 2,
		},
		{
			name: "score 3 if match two chars",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: true,
			}, {
				Word:     word("abc"),
				X:        2,
				Y:        0,
				Vertical: true,
			}},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("aba"),
				X:        0,
				Y:        0,
				Vertical: false,
			},
			want: 3,
		},
		{
			name: "score 0 if second char overlap is not a match",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: true,
			}, {
				Word:     word("xxx"),
				X:        2,
				Y:        0,
				Vertical: true,
			}},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("aba"),
				X:        0,
				Y:        0,
				Vertical: false,
			},
			want: 0,
		},
		{
			name: "score 0 if there is no gap to the next char",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: true,
			}, {
				Word:     word("xxx"),
				X:        2,
				Y:        0,
				Vertical: true,
			}},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("ab"),
				X:        0,
				Y:        0,
				Vertical: false,
			},
			want: 0,
		},
		{
			name: "score 0 if there is no gap above a word",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: true,
			}, {
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: false,
			},
			},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("bxx"),
				X:        0,
				Y:        0,
				Vertical: false,
			},
			want: 0,
		},
		{
			name: "score 0 if there is no gap below a word",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: true,
			}, {
				Word:     word("abc"),
				X:        0,
				Y:        2,
				Vertical: false,
			},
			},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("bxx"),
				X:        0,
				Y:        0,
				Vertical: false,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run("HORIZONTAL: "+tt.name, func(t *testing.T) {
			for _, v := range tt.existingWords {
				tt.generator.placeWord(v)
			}
			require.EqualValues(t, tt.want, tt.generator.scorePlacement(tt.placement))
		})
	}
}

func TestGenerator_scorePlacement_verticalWords(t *testing.T) {
	tests := []struct {
		name          string
		existingWords []Placement
		generator     *Generator
		placement     Placement
		want          int
	}{
		{
			name:      "score 1 if no X collisions",
			generator: NewGenerator(4),
			placement: Placement{
				Word:     word("foo"),
				X:        0,
				Y:        0,
				Vertical: true,
			},
			want: 1,
		},
		{
			name: "score 0 if preceding cell collides",
			existingWords: []Placement{{
				Word:     word("abcd"),
				X:        0,
				Y:        0,
				Vertical: true,
			}},
			generator: NewGenerator(4),
			placement: Placement{
				Word:     word("bcd"),
				X:        0,
				Y:        1,
				Vertical: true,
			},
			want: 0,
		},
		{
			name: "score 0 if following cell collides",
			existingWords: []Placement{{
				Word:     word("abcd"),
				X:        0,
				Y:        0,
				Vertical: false,
			}},
			generator: NewGenerator(4),
			placement: Placement{
				Word:     word("abc"),
				X:        1,
				Y:        0,
				Vertical: false,
			},
			want: 0,
		},
		{
			name: "score 2 if match one char",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: false,
			}},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: true,
			},
			want: 2,
		},
		{
			name: "score 3 if match two chars",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: false,
			}, {
				Word:     word("abc"),
				X:        0,
				Y:        2,
				Vertical: false,
			}},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("aba"),
				X:        0,
				Y:        0,
				Vertical: true,
			},
			want: 3,
		},
		{
			name: "score 0 if second char overlap is not a match",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: false,
			}, {
				Word:     word("xxx"),
				X:        0,
				Y:        2,
				Vertical: false,
			}},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("aba"),
				X:        0,
				Y:        0,
				Vertical: true,
			},
			want: 0,
		},
		{
			name: "score 0 if there is no gap to the next char",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: false,
			}, {
				Word:     word("xxx"),
				X:        0,
				Y:        2,
				Vertical: false,
			}},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("ab"),
				X:        0,
				Y:        0,
				Vertical: true,
			},
			want: 0,
		},
		{
			name: "score 0 if there is no gap to the left of the word",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: false,
			}, {
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: true,
			},
			},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("bxx"),
				X:        1,
				Y:        0,
				Vertical: true,
			},
			want: 0,
		},
		{
			name: "score 0 if there is no gap to the right of the word",
			existingWords: []Placement{{
				Word:     word("abc"),
				X:        0,
				Y:        0,
				Vertical: false,
			}, {
				Word:     word("cba"),
				X:        2,
				Y:        0,
				Vertical: true,
			},
			},
			generator: NewGenerator(3),
			placement: Placement{
				Word:     word("bxx"),
				X:        1,
				Y:        0,
				Vertical: true,
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run("VERTICAL: "+tt.name, func(t *testing.T) {
			for _, v := range tt.existingWords {
				tt.generator.placeWord(v)
			}
			require.EqualValues(t, tt.want, tt.generator.scorePlacement(tt.placement))
		})
	}
}

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name            string
		generator       *Generator
		words           []Word
		wantCrossword   string
		wantPlacedWords []Placement
		solve           bool
	}{
		{
			name:      "single word placed",
			generator: NewGenerator(3),
			words:     []Word{{Word: "foo"}},
			solve:     true,
			wantCrossword: `
FOO
###
###`,
		}, {
			name:      "single word placed, but unsolved",
			generator: NewGenerator(3),
			words:     []Word{{Word: "foo"}},
			solve:     false,
			wantCrossword: `
???
###
###`,
		}, {
			name:      "two words placed",
			generator: NewGenerator(4),
			words:     []Word{{Word: "food"}, {Word: "fud"}},
			solve:     true,
			wantCrossword: `
FOOD
U###
D###
####`,
		}, {
			name:      "three words placed",
			generator: NewGenerator(4),
			words:     []Word{{Word: "food"}, {Word: "fud"}, {Word: "duff"}},
			solve:     true,
			wantCrossword: `
FOOD
U##U
D##F
###F`,
		}, {
			name:      "four words placed",
			generator: NewGenerator(4),
			words:     []Word{{Word: "food"}, {Word: "fud"}, {Word: "duff"}, {Word: "dxxf"}},
			solve:     true,
			wantCrossword: `
FOOD
U##U
DXXF
###F`,
			wantPlacedWords: []Placement{
				{
					id:       1,
					Word:     Word{Word: "FOOD"},
					X:        0,
					Y:        0,
					Vertical: false,
					Solved:   true,
				}, {
					id:       2,
					Word:     Word{Word: "DUFF"},
					X:        3,
					Y:        0,
					Vertical: true,
					Solved:   true,
				}, {
					id:       3,
					Word:     Word{Word: "DXXF"},
					X:        0,
					Y:        2,
					Vertical: false,
					Solved:   true,
				}, {
					id:       4,
					Word:     Word{Word: "FUD"},
					X:        0,
					Y:        0,
					Vertical: true,
					Solved:   true,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cw := tt.generator.Generate(tt.words)
			if tt.solve {
				cw.Solve()
			}
			got := strings.TrimSpace(RenderText(cw))
			if !assert.EqualValues(t, strings.TrimSpace(tt.wantCrossword), got) {
				fmt.Println("WANT")
				fmt.Println(tt.wantCrossword)
				fmt.Println("GOT")
				fmt.Println(got)
			}
			if tt.wantPlacedWords != nil {
				require.EqualValues(t, tt.wantPlacedWords, cw.Words)
			}
		})
	}
}
