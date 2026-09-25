package sudoku

import (
	"errors"
	"io"
	"strings"
	"testing"
)

// classicPuzzle is the well-known example puzzle used throughout the
// README, kept here so tests don't depend on that doc staying in sync.
const classicPuzzle = "530070000600195000098000060800060003400803001700020006060000280000419005000080079"

func TestParseLine(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
	}{
		{name: "valid puzzle", line: classicPuzzle},
		{name: "valid, all dots", line: strings.Repeat(".", Size*Size)},
		{name: "valid, all zeros", line: strings.Repeat("0", Size*Size)},
		{name: "valid, all underscores", line: strings.Repeat("_", Size*Size)},
		{name: "empty line", line: "", wantErr: true},
		{name: "too short", line: strings.Repeat(".", Size*Size-1), wantErr: true},
		{name: "too long", line: strings.Repeat(".", Size*Size+1), wantErr: true},
		{name: "letter instead of digit", line: strings.Repeat(".", 40) + "x" + strings.Repeat(".", 40), wantErr: true},
		{name: "punctuation instead of digit", line: strings.Repeat(".", 40) + "-" + strings.Repeat(".", 40), wantErr: true},
		{name: "valid, single digit among blanks", line: strings.Repeat(".", 40) + "9" + strings.Repeat(".", 40)},
		{
			name:    "row conflict",
			line:    "11" + strings.Repeat(".", Size*Size-2),
			wantErr: true,
		},
		{
			// position 0 and position 9 are both column 0, one row apart.
			name:    "column conflict",
			line:    "5" + strings.Repeat(".", 8) + "5" + strings.Repeat(".", Size*Size-10),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := ParseLine(tt.line)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseLine(%q) = %v, want error", tt.line, b)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseLine(%q) returned error: %v", tt.line, err)
			}
		})
	}
}

func TestParseLineValues(t *testing.T) {
	b, err := ParseLine(classicPuzzle)
	if err != nil {
		t.Fatalf("ParseLine returned error: %v", err)
	}
	want := [Size]uint8{5, 3, 0, 0, 7, 0, 0, 0, 0}
	if b[0] != want {
		t.Errorf("row 0 = %v, want %v", b[0], want)
	}
	if b[8][8] != 9 {
		t.Errorf("cell (8,8) = %d, want 9", b[8][8])
	}
}

func TestParserNext(t *testing.T) {
	input := "# a comment\n\n" + classicPuzzle + "\n" + strings.Repeat(".", Size*Size) + "\n"
	p := NewParser(strings.NewReader(input))

	b, err := p.Next()
	if err != nil {
		t.Fatalf("first Next() returned error: %v", err)
	}
	if b[0][0] != 5 {
		t.Errorf("first board cell (0,0) = %d, want 5", b[0][0])
	}
	if got, want := p.Line(), 3; got != want {
		t.Errorf("Line() after first board = %d, want %d", got, want)
	}

	b, err = p.Next()
	if err != nil {
		t.Fatalf("second Next() returned error: %v", err)
	}
	if b != (Board{}) {
		t.Errorf("second board = %v, want all blank", b)
	}
	if got, want := p.Line(), 4; got != want {
		t.Errorf("Line() after second board = %d, want %d", got, want)
	}

	if _, err := p.Next(); !errors.Is(err, io.EOF) {
		t.Errorf("third Next() = %v, want io.EOF", err)
	}
}

func TestParserNextMalformedLine(t *testing.T) {
	p := NewParser(strings.NewReader("not a puzzle\n"))
	_, err := p.Next()
	if err == nil {
		t.Fatal("Next() = nil, want error")
	}
	if !strings.Contains(err.Error(), "line 1") {
		t.Errorf("error = %q, want it to name line 1", err.Error())
	}
}

func TestParserNextLineTooLong(t *testing.T) {
	p := NewParser(strings.NewReader(strings.Repeat("1", maxLineLen+1) + "\n"))
	_, err := p.Next()
	if err == nil {
		t.Fatal("Next() = nil, want error")
	}
}

func TestParserNextEmptyInput(t *testing.T) {
	p := NewParser(strings.NewReader(""))
	if _, err := p.Next(); !errors.Is(err, io.EOF) {
		t.Errorf("Next() on empty input = %v, want io.EOF", err)
	}
}
