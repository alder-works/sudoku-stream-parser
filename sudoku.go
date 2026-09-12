// Package sudoku parses and prints 9x9 sudoku boards, and is built to
// process files that contain many boards (one per line) without ever
// holding more than one board in memory at a time. See Parser for the
// streaming entry point.
package sudoku

import "fmt"

// Size is the side length of a standard sudoku board.
const Size = 9

// boxSize is the side length of the 3x3 sub-grids.
const boxSize = 3

// Board holds one 9x9 sudoku grid. A value of 0 means the cell is empty;
// values 1-9 are filled cells. The zero value is a fully empty board.
type Board [Size][Size]uint8

// At returns the value at (row, col), 0 for empty. It panics if row or
// col is outside [0, Size), same as slice indexing would.
func (b Board) At(row, col int) uint8 {
	return b[row][col]
}

// Filled reports whether every cell has a value, i.e. the board has no
// blanks left. It does not check that the board is a valid solution.
func (b Board) Filled() bool {
	for _, row := range b {
		for _, v := range row {
			if v == 0 {
				return false
			}
		}
	}
	return true
}

// ConflictError describes a rule violation found while validating a
// board: two identical digits sharing a row, column, or 3x3 box.
type ConflictError struct {
	Value      uint8
	RowA, ColA int
	RowB, ColB int
	Unit       string // "row", "column", or "box"
}

func (e *ConflictError) Error() string {
	return fmt.Sprintf("duplicate %d in %s: (%d,%d) and (%d,%d)",
		e.Value, e.Unit, e.RowA, e.ColA, e.RowB, e.ColB)
}

// Validate checks that no row, column, or 3x3 box contains the same
// digit twice. Empty cells never conflict. It returns the first
// conflict found, or nil if the board obeys the sudoku placement rules
// (the board need not be fully filled to pass).
func (b Board) Validate() error {
	if err := checkUnits(b, "row", rowUnits()); err != nil {
		return err
	}
	if err := checkUnits(b, "column", colUnits()); err != nil {
		return err
	}
	return checkUnits(b, "box", boxUnits())
}

// cellPos is a (row, col) coordinate within a board.
type cellPos struct{ row, col int }

func rowUnits() [][]cellPos {
	units := make([][]cellPos, Size)
	for r := 0; r < Size; r++ {
		unit := make([]cellPos, Size)
		for c := 0; c < Size; c++ {
			unit[c] = cellPos{r, c}
		}
		units[r] = unit
	}
	return units
}

func colUnits() [][]cellPos {
	units := make([][]cellPos, Size)
	for c := 0; c < Size; c++ {
		unit := make([]cellPos, Size)
		for r := 0; r < Size; r++ {
			unit[r] = cellPos{r, c}
		}
		units[c] = unit
	}
	return units
}

func boxUnits() [][]cellPos {
	units := make([][]cellPos, 0, Size)
	for br := 0; br < Size; br += boxSize {
		for bc := 0; bc < Size; bc += boxSize {
			unit := make([]cellPos, 0, Size)
			for r := br; r < br+boxSize; r++ {
				for c := bc; c < bc+boxSize; c++ {
					unit = append(unit, cellPos{r, c})
				}
			}
			units = append(units, unit)
		}
	}
	return units
}

func checkUnits(b Board, name string, units [][]cellPos) error {
	for _, unit := range units {
		var seen [Size + 1]*cellPos // index by digit 1-9
		for i := range unit {
			pos := unit[i]
			v := b[pos.row][pos.col]
			if v == 0 {
				continue
			}
			if prev := seen[v]; prev != nil {
				return &ConflictError{
					Value: v,
					RowA:  prev.row, ColA: prev.col,
					RowB: pos.row, ColB: pos.col,
					Unit: name,
				}
			}
			seen[v] = &cellPos{pos.row, pos.col}
		}
	}
	return nil
}
