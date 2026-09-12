package sudoku

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// maxLineLen bounds a single input line. A valid puzzle line is 81
// characters; this leaves generous room for trailing whitespace or a
// puzzle ID column some datasets append, while still catching garbage
// input (or the wrong file entirely) as a clear error instead of
// letting the scanner's buffer grow to swallow it.
const maxLineLen = 256

// Parser reads sudoku boards one at a time from a stream. It is built
// for files that hold many puzzles, one per line - a common shape for
// sudoku datasets - and it never buffers more than a single line, so
// memory use stays flat no matter how large the input is.
type Parser struct {
	scanner *bufio.Scanner
	lineNo  int
}

// NewParser wraps r for streaming reads. r is read incrementally as
// Next is called; nothing is read up front.
func NewParser(r io.Reader) *Parser {
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 0, maxLineLen), maxLineLen)
	return &Parser{scanner: sc}
}

// Next reads and validates the next board in the stream. Blank lines
// and lines starting with '#' are skipped so a puzzle file can carry
// comments or spacing between entries. It returns io.EOF once the
// stream is exhausted.
func (p *Parser) Next() (Board, error) {
	for p.scanner.Scan() {
		p.lineNo++
		line := strings.TrimSpace(p.scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		b, err := ParseLine(line)
		if err != nil {
			return Board{}, fmt.Errorf("line %d: %w", p.lineNo, err)
		}
		return b, nil
	}
	if err := p.scanner.Err(); err != nil {
		return Board{}, fmt.Errorf("line %d: %w", p.lineNo+1, err)
	}
	return Board{}, io.EOF
}

// Line reports the 1-based line number of the board most recently
// returned by Next, for callers that want to report their own errors
// against the source file.
func (p *Parser) Line() int {
	return p.lineNo
}

// ParseLine parses one board from an 81-character line and validates
// it. A cell is a digit 1-9 for a filled square, or one of '.', '0',
// '_' for an empty one. Any other character, or a line whose length
// isn't 81 once surrounding whitespace is trimmed, is an error.
func ParseLine(line string) (Board, error) {
	if len(line) != Size*Size {
		return Board{}, fmt.Errorf("sudoku: want %d characters, got %d", Size*Size, len(line))
	}
	var b Board
	for i := 0; i < len(line); i++ {
		ch := line[i]
		row, col := i/Size, i%Size
		switch {
		case ch >= '1' && ch <= '9':
			b[row][col] = uint8(ch - '0')
		case ch == '.' || ch == '0' || ch == '_':
			b[row][col] = 0
		default:
			return Board{}, fmt.Errorf("sudoku: invalid character %q at position %d", ch, i)
		}
	}
	if err := b.Validate(); err != nil {
		return Board{}, err
	}
	return b, nil
}
