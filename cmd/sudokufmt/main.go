// Command sudokufmt reads sudoku puzzles, one per line, from stdin or
// a file, validates each one, and pretty-prints it to stdout. It is
// meant to work on puzzle collections of any size: input is streamed
// through sudoku.Parser rather than read into memory up front.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/alder-works/sudoku-stream-parser"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "sudokufmt:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	r := os.Stdin
	if len(args) > 0 {
		f, err := os.Open(args[0])
		if err != nil {
			return err
		}
		defer f.Close()
		r = f
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	p := sudoku.NewParser(r)
	count := 0
	for {
		board, err := p.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		if count > 0 {
			fmt.Fprintln(out)
		}
		if err := sudoku.Fprint(out, board); err != nil {
			return err
		}
		count++
	}
	if count == 0 {
		return errors.New("no puzzles found in input")
	}
	return nil
}
