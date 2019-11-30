package chess2

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseFen(t *testing.T) {
	board, err := ParseFen(FenDefault)
	require.NoError(t, err)
	tests := []struct {
		x, y int
		p    Piece
	}{
		{x: 0, y: 0, p: NewPiece(TypeRook, ArmyNone, ColorBlack)},
		{x: 3, y: 1, p: NewPiece(TypePawn, ArmyNone, ColorBlack)},
		{x: 7, y: 6, p: NewPiece(TypePawn, ArmyNone, ColorWhite)},
		{x: 4, y: 7, p: NewPiece(TypeKing, ArmyNone, ColorWhite)},
	}
	for _, data := range tests {
		square := SquareFromCoords(data.x, data.y)
		piece, _ := board.PieceAt(square)
		assert.Equal(t, data.p, piece, "Expected piece at %v to be %s, was %s", square, data.p, piece)
	}
}

func TestEncodeFen(t *testing.T) {
	var board Board
	board.SetPieceAt(SquareFromCoords(3, 3), NewPiece(TypeKing, ArmyTwoKings, ColorWhite))
	fen := EncodeFen(board)
	expected := "8/8/8/3K4/8/8/8/8"
	require.Equal(t, expected, fen)
}

func TestRoundTripEmptyFen(t *testing.T) {
	fen := FenEmpty
	board, err := ParseFen(fen)
	require.NoError(t, err)
	result := EncodeFen(board)
	require.Equal(t, fen, result)
}

func TestRoundTripStartingFen(t *testing.T) {
	fen := FenDefault
	board, err := ParseFen(fen)
	require.NoError(t, err)
	result := EncodeFen(board)
	require.Equal(t, fen, result)
}
