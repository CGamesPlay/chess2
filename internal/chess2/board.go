package chess2

import (
	"fmt"
)

// A Square represents a location on the board.
type Square struct {
	addr uint8
}

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

// Board represents the pieces on a chess board.
type Board struct {
	// Mask of where the pieces are
	kings, queens, bishops, knights, rooks, pawns uint64
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
	switch {
	case b.kings&squareMask != 0:
		pieceType = TypeKing
	case b.queens&squareMask != 0:
		pieceType = TypeQueen
	case b.bishops&squareMask != 0:
		pieceType = TypeBishop
	case b.knights&squareMask != 0:
		pieceType = TypeKnight
	case b.rooks&squareMask != 0:
		pieceType = TypeRook
	case b.pawns&squareMask != 0:
		pieceType = TypePawn
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
	switch p.Type() {
	case TypeKing:
		b.kings |= squareMask
	case TypeQueen:
		b.queens |= squareMask
	case TypeBishop:
		b.bishops |= squareMask
	case TypeKnight:
		b.knights |= squareMask
	case TypeRook:
		b.rooks |= squareMask
	case TypePawn:
		b.pawns |= squareMask
	}
}

// ClearPieceAt adjusts the receiver to have an empty square at the provided
// space.
func (b *Board) ClearPieceAt(s Square) {
	squareMask := s.mask()
	b.kings &= ^squareMask
	b.queens &= ^squareMask
	b.bishops &= ^squareMask
	b.knights &= ^squareMask
	b.rooks &= ^squareMask
	b.pawns &= ^squareMask
	b.white &= ^squareMask
	b.black &= ^squareMask
}
