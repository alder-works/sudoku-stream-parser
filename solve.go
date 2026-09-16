package sudoku

import (
	"errors"
	"fmt"
	"math/bits"
)

// ErrUnsolvable is returned by Solve when a board, though free of rule
// conflicts, has no completion that satisfies the sudoku constraints.
var ErrUnsolvable = errors.New("sudoku: no solution")

// fullCandidateMask has bits 0-8 set, one per digit 1-9.
const fullCandidateMask = uint16(1<<Size) - 1

// Solve returns a completed copy of b with every blank cell filled in.
// b itself is left untouched. It first eliminates candidates that are
// forced by the current placements (constraint propagation), then
// falls back to backtracking search - trying each remaining candidate
// for the most constrained blank cell in turn - only where propagation
// alone can't finish the board. It returns ErrUnsolvable if no
// completion exists, or the underlying ConflictError if b already
// breaks the placement rules.
func Solve(b Board) (Board, error) {
	if err := b.Validate(); err != nil {
		return Board{}, fmt.Errorf("sudoku: cannot solve invalid board: %w", err)
	}
	solved, ok := solveBoard(b)
	if !ok {
		return Board{}, ErrUnsolvable
	}
	return solved, nil
}

// solveBoard tries to complete b, returning the completed board and
// true on success. b is passed and returned by value, so callers can
// try a candidate and simply discard the result on failure without any
// separate undo step.
func solveBoard(b Board) (Board, bool) {
	if !propagateSingles(&b) {
		return Board{}, false
	}
	r, c, found := mostConstrainedCell(b)
	if !found {
		return b, true
	}
	mask := candidateMask(b, r, c)
	for v := uint8(1); v <= Size; v++ {
		bit := uint16(1) << (v - 1)
		if mask&bit == 0 {
			continue
		}
		next := b
		next[r][c] = v
		if solved, ok := solveBoard(next); ok {
			return solved, true
		}
	}
	return Board{}, false
}

// propagateSingles repeatedly fills in any blank cell that has exactly
// one remaining candidate, until no such cell is left. It reports
// false if some blank cell is left with zero candidates, meaning the
// board (or the branch of a search that produced it) has no solution.
func propagateSingles(b *Board) bool {
	for {
		changed := false
		for r := 0; r < Size; r++ {
			for c := 0; c < Size; c++ {
				if b[r][c] != 0 {
					continue
				}
				mask := candidateMask(*b, r, c)
				switch bits.OnesCount16(mask) {
				case 0:
					return false
				case 1:
					b[r][c] = uint8(bits.TrailingZeros16(mask)) + 1
					changed = true
				}
			}
		}
		if !changed {
			return true
		}
	}
}

// mostConstrainedCell returns the blank cell with the fewest remaining
// candidates, the standard minimum-remaining-value heuristic for
// keeping backtracking search shallow. found is false if b has no
// blank cells at all.
func mostConstrainedCell(b Board) (row, col int, found bool) {
	best := Size + 1
	for r := 0; r < Size; r++ {
		for c := 0; c < Size; c++ {
			if b[r][c] != 0 {
				continue
			}
			if n := bits.OnesCount16(candidateMask(b, r, c)); n < best {
				best, row, col, found = n, r, c, true
			}
		}
	}
	return row, col, found
}

// candidateMask returns the digits 1-9 that could legally go in the
// blank cell (row, col), as a bitmask with bit (v-1) set for digit v.
func candidateMask(b Board, row, col int) uint16 {
	return fullCandidateMask &^ peerMask(b, row, col)
}

// peerMask returns the digits already placed in the row, column, or
// box that (row, col) belongs to.
func peerMask(b Board, row, col int) uint16 {
	var used uint16
	for i := 0; i < Size; i++ {
		used |= digitBit(b[row][i])
		used |= digitBit(b[i][col])
	}
	boxRow, boxCol := (row/boxSize)*boxSize, (col/boxSize)*boxSize
	for r := boxRow; r < boxRow+boxSize; r++ {
		for c := boxCol; c < boxCol+boxSize; c++ {
			used |= digitBit(b[r][c])
		}
	}
	return used
}

// digitBit returns the candidate-mask bit for v, or 0 for an empty cell.
func digitBit(v uint8) uint16 {
	if v == 0 {
		return 0
	}
	return 1 << (v - 1)
}
