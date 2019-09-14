package chess2

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSquareString(t *testing.T) {
	require.Equal(t, "a1", SquareFromCoords(0, 7).String())
	require.Equal(t, "a8", SquareFromCoords(0, 0).String())
	require.Equal(t, "h8", SquareFromCoords(7, 0).String())
	require.Equal(t, "h1", SquareFromCoords(7, 7).String())
}

func TestSquareFromName(t *testing.T) {
	require.Equal(t, "a1", SquareFromName("A1").String())
	require.Equal(t, "a8", SquareFromName("a8").String())
	require.Equal(t, "h8", SquareFromName("H8").String())
	require.Equal(t, "h1", SquareFromName("h1").String())
}

func TestReplacePieces(t *testing.T) {
	board, err := ParseFen(FenDefault)
	require.NoError(t, err)
	board.ReplacePieces(ColorBlack, TypePawn, TypeQueen)
	expected := "rnbqkbnr/qqqqqqqq/8/8/8/8/PPPPPPPP/RNBQKBNR"
	fen := EncodeFen(board)
	require.Equal(t, expected, fen)
}
