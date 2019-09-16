package chess2

import (
	"fmt"
)

const (
	// MaskEmpty is a mask of an empty board.
	MaskEmpty = uint64(0)
	// MaskFull is a mask of an entirely occupied board.
	MaskFull = ^uint64(0)
)

var (
	// MaskRank is a mask of the row at the given y-coordinate.
	MaskRank = []uint64{
		0x00000000000000ff,
		0x000000000000ff00,
		0x0000000000ff0000,
		0x00000000ff000000,
		0x000000ff00000000,
		0x0000ff0000000000,
		0x00ff000000000000,
		0xff00000000000000,
	}
)

// A Square represents a location on the board.
type Square struct {
	Address uint8
}

// InvalidSquare is not a real Square. It will cause undefined behavior if you
// use it as anything other than a sentinel.
var InvalidSquare = Square{127}

// SquareFromCoords takes an x, y pair and returns a Square.
func SquareFromCoords(x, y int) Square {
	if x < 0 || x > 7 || y < 0 || y > 7 {
		panic("Invalid coords for Square")
	}
	return Square{Address: uint8(y*8 + x)}
}

// SquareFromName takes a name like A1 and returns a Square. Will return
// InvalidSquare if the name is not valid.
func SquareFromName(name string) Square {
	x := name[0]
	y := name[1]
	if 'a' <= x && x <= 'h' {
		x &= ^uint8(0x20)
	}
	if x < 'A' || x > 'H' || y < '1' || y > '8' || len(name) != 2 {
		return InvalidSquare
	}
	x = x - 'A'
	y = '8' - y
	return Square{Address: uint8(y*8 + x)}
}

// Mask returns a bitmask that selects only this square.
func (s Square) Mask() uint64 {
	return uint64(1 << s.Address)
}

// Y value of the receiver. Corresponds to y = 8 - rank of the square.
func (s Square) Y() int {
	return int(s.Address / 8)
}

// X value of the receiver. Corresponds to x = file of the square - 'a'.
func (s Square) X() int {
	return int(s.Address % 8)
}

func (s Square) String() string {
	if s == InvalidSquare {
		return "--"
	}
	return fmt.Sprintf("%c%d", s.X()+'a', 8-s.Y())
}

// SquareDistance calculates the greater of the horizontal or vertical distance
// between the squares.
func SquareDistance(a, b Square) int {
	dx := abs(a.X() - b.X())
	dy := abs(a.Y() - b.Y())
	if dx > dy {
		return dx
	}
	return dy
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
	colors [2]uint64
}

// PieceAt returns the piece at the given square in the receiver, and a
// boolean indicating whether the square is occupied. The returned piece will
// always have ArmyNone as the Army.
func (b *Board) PieceAt(s Square) (Piece, bool) {
	squareMask := s.Mask()
	var (
		color     Color
		pieceType PieceType
	)
	if (b.colors[0]|b.colors[1])&squareMask == 0 {
		return Piece{}, false
	} else if b.colors[0]&squareMask != 0 {
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
	squareMask := s.Mask()
	b.colors[ColorIdx(p.Color())] |= squareMask
	b.pieces[pieceTypeIdx(p.Type())] |= squareMask
}

// ClearPieceAt adjusts the receiver to have an empty square at the provided
// space.
func (b *Board) ClearPieceAt(s Square) {
	squareMask := s.Mask()
	for i := range b.pieces {
		b.pieces[i] &= ^squareMask
	}
	b.colors[0] &= ^squareMask
	b.colors[1] &= ^squareMask
}

// ReplacePieces modifies the board so that all pieces of the given color and
// find type are replaced by a corresponding piece of the same color of the
// replace type.
func (b *Board) ReplacePieces(color Color, find, replace PieceType) {
	mask := b.colors[ColorIdx(color)]
	idx := pieceTypeIdx(find)
	mask &= b.pieces[idx]
	b.pieces[idx] &= ^mask
	idx = pieceTypeIdx(replace)
	b.pieces[idx] |= mask
}

func (b *Board) occupiedMask() uint64 {
	return b.colors[0] | b.colors[1]
}

func (b *Board) pieceMask(p PieceType) uint64 {
	return b.pieces[pieceTypeIdx(p)]
}

func (b *Board) colorMask(c Color) uint64 {
	return b.colors[ColorIdx(c)]
}
