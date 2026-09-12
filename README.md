# sudoku-stream-parser

A Go library for parsing and printing sudoku boards, built around one
requirement: it should be able to work through a file of a million
puzzles without holding more than one of them in memory at a time.

Puzzle datasets in the wild (the well-known 17-clue collections, for
example) are usually plain text, one 81-character line per board,
using `.` or `0` for blank cells. `Parser` reads that format line by
line off an `io.Reader`, so processing a multi-gigabyte file costs the
same handful of bytes as processing a single puzzle.

## Format

Each line is 81 characters, read left to right, top to bottom. A cell
is a digit `1`-`9`, or one of `.`, `0`, `_` for empty. Blank lines and
lines starting with `#` are skipped, so files can carry comments:

```
# easy puzzles
53..7....6..195....98....6.8...6...34..8.3..17...2...6.6....28....419..5....8..79
```

## Library usage

```go
package main

import (
	"fmt"
	"os"

	"github.com/alder-works/sudoku-stream-parser"
)

func main() {
	f, err := os.Open("puzzles.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := sudoku.NewParser(f)
	for {
		board, err := p.Next()
		if err != nil {
			break // io.EOF when the file is exhausted
		}
		fmt.Println("row 0:", board[0])
	}
}
```

`ParseLine` is available directly if you already have a single line
and don't need the streaming reader.

Validation checks the sudoku placement rules - no repeated digit in
any row, column, or 3x3 box - and returns a `*ConflictError` naming
the two cells that collide. A board doesn't need to be fully filled to
pass; a partially solved puzzle with no conflicts is valid.

## Command line

`cmd/sudokufmt` reads a puzzle file (or stdin) and pretty-prints each
board:

```
$ go run ./cmd/sudokufmt puzzles.txt
5 3 . | . 7 . | . . .
6 . . | 1 9 5 | . . .
. 9 8 | . . . | . 6 .
------+-------+------
8 . . | . 6 . | . . 3
4 . . | 8 . 3 | . . 1
7 . . | . 2 . | . . 6
------+-------+------
. 6 . | . . . | 2 8 .
. . . | 4 1 9 | . . 5
. . . | . 8 . | . 7 9
```

Any line that fails to parse or breaks a sudoku rule stops the run
with an error naming the offending line number.

## Status

Early skeleton: single-line parsing, validation, and the boxed printer
work. The multi-line grid format (nine rows of text per puzzle,
matching the printer's own output) isn't read back in yet - see the
roadmap.
