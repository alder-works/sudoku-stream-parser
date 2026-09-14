package sudoku

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// GridParser reads sudoku boards from the boxed, multi-line layout that
// Fprint writes: nine row lines, with a divider line after every third
// row. It is the round-trip counterpart to Fprint, the way Parser is
// the round-trip counterpart to Compact. Boards in the stream may be
// separated by blank lines or '#' comments, matching what a file made
// by running Fprint over several boards in a row looks like.
type GridParser struct {
	scanner *bufio.Scanner
	lineNo  int
}

// NewGridParser wraps r for streaming reads of the boxed grid format.
func NewGridParser(r io.Reader) *GridParser {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, maxLineLen), maxLineLen)
	return &GridParser{scanner: sc}
}

// Next reads and validates the next boxed board in the stream. Blank
// lines and lines starting with '#' are skipped wherever a new board
// is expected to start. It returns io.EOF once the stream is
// exhausted.
func (p *GridParser) Next() (Board, error) {
	first, ok := p.nextContentLine()
	if !ok {
		if err := p.scanner.Err(); err != nil {
			return Board{}, fmt.Errorf("line %d: %w", p.lineNo+1, err)
		}
		return Board{}, io.EOF
	}

	var b Board
	row, err := parseGridRow(first)
	if err != nil {
		return Board{}, fmt.Errorf("line %d: %w", p.lineNo, err)
	}
	b[0] = row

	for r := 1; r < Size; r++ {
		if r%boxSize == 0 {
			line, ok := p.nextLine()
			if !ok {
				return Board{}, fmt.Errorf("line %d: unexpected end of input, want divider", p.lineNo+1)
			}
			if line != divider {
				return Board{}, fmt.Errorf("line %d: want divider %q, got %q", p.lineNo, divider, line)
			}
		}
		line, ok := p.nextLine()
		if !ok {
			return Board{}, fmt.Errorf("line %d: unexpected end of input, want row %d of 9", p.lineNo+1, r+1)
		}
		row, err := parseGridRow(line)
		if err != nil {
			return Board{}, fmt.Errorf("line %d: %w", p.lineNo, err)
		}
		b[r] = row
	}

	if err := b.Validate(); err != nil {
		return Board{}, err
	}
	return b, nil
}

// Line reports the 1-based line number of the last input line consumed
// by Next, for callers that want to report their own errors against
// the source file.
func (p *GridParser) Line() int {
	return p.lineNo
}

// nextLine returns the next line, trimmed of surrounding whitespace,
// with no skipping. It reports false at end of input.
func (p *GridParser) nextLine() (string, bool) {
	if !p.scanner.Scan() {
		return "", false
	}
	p.lineNo++
	return strings.TrimSpace(p.scanner.Text()), true
}

// nextContentLine skips blank lines and '#' comments and returns the
// next line that has content. It reports false at end of input.
func (p *GridParser) nextContentLine() (string, bool) {
	for {
		line, ok := p.nextLine()
		if !ok {
			return "", false
		}
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		return line, true
	}
}

// parseGridRow parses one row of the boxed layout, e.g.
// "5 3 . | . 7 . | . . .", into the nine cell values it encodes.
func parseGridRow(line string) ([Size]uint8, error) {
	fields := strings.Fields(strings.ReplaceAll(line, "|", " "))
	if len(fields) != Size {
		return [Size]uint8{}, fmt.Errorf("sudoku: want %d cells, got %d", Size, len(fields))
	}
	var row [Size]uint8
	for i, f := range fields {
		if len(f) != 1 {
			return [Size]uint8{}, fmt.Errorf("sudoku: invalid cell %q", f)
		}
		ch := f[0]
		switch {
		case ch >= '1' && ch <= '9':
			row[i] = uint8(ch - '0')
		case ch == '.' || ch == '0' || ch == '_':
			row[i] = 0
		default:
			return [Size]uint8{}, fmt.Errorf("sudoku: invalid character %q", ch)
		}
	}
	return row, nil
}
