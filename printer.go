package sudoku

import (
	"bufio"
	"io"
)

// divider separates the 3x3 boxes vertically. It matches the width of
// a row rendered by Fprint: three digits and two spaces per box,
// joined by " | ".
const divider = "------+-------+------"

// Fprint writes b to w in a boxed layout, using '.' for empty cells:
//
//	5 3 . | . 7 . | . . .
//	6 . . | 1 9 5 | . . .
//	. 9 8 | . . . | . 6 .
//	------+-------+------
//	8 . . | . 6 . | . . 3
//	...
//
// It writes directly to w rather than building the whole board as a
// string first, so printing many boards in sequence costs no more
// memory than printing one.
func Fprint(w io.Writer, b Board) error {
	bw, ok := w.(*bufio.Writer)
	if !ok {
		bw = bufio.NewWriter(w)
	}
	for r := 0; r < Size; r++ {
		if r > 0 && r%boxSize == 0 {
			if _, err := bw.WriteString(divider + "\n"); err != nil {
				return err
			}
		}
		if err := writeRow(bw, b[r]); err != nil {
			return err
		}
	}
	return bw.Flush()
}

func writeRow(bw *bufio.Writer, row [Size]uint8) error {
	for c := 0; c < Size; c++ {
		if c > 0 && c%boxSize == 0 {
			if err := bw.WriteByte('|'); err != nil {
				return err
			}
			if err := bw.WriteByte(' '); err != nil {
				return err
			}
		}
		if err := bw.WriteByte(cellGlyph(row[c])); err != nil {
			return err
		}
		if c < Size-1 {
			if err := bw.WriteByte(' '); err != nil {
				return err
			}
		}
	}
	return bw.WriteByte('\n')
}

func cellGlyph(v uint8) byte {
	if v == 0 {
		return '.'
	}
	return '0' + v
}

// Compact writes b as a single 81-character line, the same format
// ParseLine reads, using '.' for empty cells. It is the round-trip
// counterpart to ParseLine and is the form to use when writing a
// stream of many boards back out one per line.
func Compact(w io.Writer, b Board) error {
	var buf [Size*Size + 1]byte
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			buf[r*Size+c] = cellGlyph(b[r][c])
		}
	}
	buf[Size*Size] = '\n'
	_, err := w.Write(buf[:])
	return err
}
