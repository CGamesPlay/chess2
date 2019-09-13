package chess2

import (
	"fmt"
)

// A Square represents a location on the board.
type Square struct {
	addr uint8
}

// InvalidSquare is not a real Square. It will cause undefined behavior if you
// use it as anything other than a sentinel.
var InvalidSquare = Square{127}

// SquareFromCoords takes an x, y pair and returns a Square.
func SquareFromCoords(x, y int) Square {
	if x < 0 || x > 7 || y < 0 || y > 7 {
		panic("Invalid coords for Square")
	}
	return Square{addr: uint8(y*8 + x)}
}

// SquareFromName takes a name like A1 and returns a Square.
func SquareFromName(name string) Square {
	x := name[0]
	y := name[1]
	if 'a' <= x && x <= 'h' {
		x &= ^uint8(0x20)
	}
	if x < 'A' || x > 'H' || y < '1' || y > '8' || len(name) != 2 {
		panic("Invalid Square name")
	}
	x = x - 'A'
	y = '8' - y
	return Square{addr: uint8(y*8 + x)}
}

func (s Square) mask() uint64 {
	return uint64(1 << s.addr)
}

func (s Square) String() string {
	x := s.addr % 8
	y := s.addr / 8
	return fmt.Sprintf("%c%d", x+'A', 8-y)
}

func pieceTypeIdx(p PieceType) int {
	return int(p) - 1
}

func idxPieceType(idx int) PieceType {
	return PieceType(idx + 1)
}

// Board represents the pieces on a chess board.
type Board struct {
	// Mask of where the pieces are
	pieces [6]uint64
	// Mask of which squares are occupied by which color
	white, black uint64
}

// PieceAt returns the piece at the given square in the receiver, and a boolean
// indicating whether the square is occupied
func (b *Board) PieceAt(s Square) (Piece, bool) {
	squareMask := s.mask()
	var (
		color     Color
		pieceType PieceType
	)
	if (b.white|b.black)&squareMask == 0 {
		return Piece{}, false
	} else if b.white&squareMask != 0 {
		color = ColorWhite
	} else {
		color = ColorBlack
	}
	for i := range b.pieces {
		if b.pieces[i]&squareMask != 0 {
			pieceType = idxPieceType(i)
			break
		}
	}
	result := NewPiece(pieceType, ArmyNone, color)
	return result, true
}

// SetPieceAt adjusts the reciever to have the provided piece at the provided
// space.
func (b *Board) SetPieceAt(s Square, p Piece) {
	b.ClearPieceAt(s)
	squareMask := s.mask()
	switch p.Color() {
	case ColorWhite:
		b.white |= squareMask
	case ColorBlack:
		b.black |= squareMask
	}
	b.pieces[pieceTypeIdx(p.Type())] |= squareMask
}

// ClearPieceAt adjusts the receiver to have an empty square at the provided
// space.
func (b *Board) ClearPieceAt(s Square) {
	squareMask := s.mask()
	for i := range b.pieces {
		b.pieces[i] &= ^squareMask
	}
	b.white &= ^squareMask
	b.black &= ^squareMask
}

// ReplacePieces modifies the board so that all pieces of the given color and
// find type are replaced by a corresponding piece of the same color of the
// replace type.
func (b *Board) ReplacePieces(color Color, find, replace PieceType) {
	var mask uint64
	if color == ColorWhite {
		mask = b.white
	} else {
		mask = b.black
	}
	idx := pieceTypeIdx(find)
	mask &= b.pieces[idx]
	b.pieces[idx] &= ^mask
	idx = pieceTypeIdx(replace)
	b.pieces[idx] |= mask
}
