package sudoku

import (
	"errors"
	"testing"
)

func TestBoardValidate(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(b *Board)
		wantUnit string // "" means no conflict expected
	}{
		{
			name:  "empty board",
			setup: func(b *Board) {},
		},
		{
			name: "row conflict",
			setup: func(b *Board) {
				b[0][0] = 5
				b[0][8] = 5
			},
			wantUnit: "row",
		},
		{
			name: "column conflict",
			setup: func(b *Board) {
				b[0][0] = 7
				b[8][0] = 7
			},
			wantUnit: "column",
		},
		{
			name: "box conflict, different row and column",
			setup: func(b *Board) {
				b[0][0] = 3
				b[2][2] = 3
			},
			wantUnit: "box",
		},
		{
			name: "same digit in different row, column, and box is fine",
			setup: func(b *Board) {
				b[0][0] = 4
				b[1][3] = 4
			},
		},
		{
			name: "repeated blanks never conflict",
			setup: func(b *Board) {
				b[0][0] = 0
				b[0][1] = 0
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var b Board
			tt.setup(&b)
			err := b.Validate()
			if tt.wantUnit == "" {
				if err != nil {
					t.Fatalf("Validate() = %v, want nil", err)
				}
				return
			}
			var ce *ConflictError
			if !errors.As(err, &ce) {
				t.Fatalf("Validate() = %v, want *ConflictError", err)
			}
			if ce.Unit != tt.wantUnit {
				t.Errorf("ConflictError.Unit = %q, want %q", ce.Unit, tt.wantUnit)
			}
		})
	}
}
