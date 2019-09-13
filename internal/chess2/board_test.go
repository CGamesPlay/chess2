package chess2

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSquareString(t *testing.T) {
	require.Equal(t, "A1", SquareFromCoords(0, 7).String())
	require.Equal(t, "A8", SquareFromCoords(0, 0).String())
	require.Equal(t, "H8", SquareFromCoords(7, 0).String())
	require.Equal(t, "H1", SquareFromCoords(7, 7).String())
}

func TestSquareFromName(t *testing.T) {
	require.Equal(t, "A1", SquareFromName("A1").String())
	require.Equal(t, "A8", SquareFromName("A8").String())
	require.Equal(t, "H8", SquareFromName("H8").String())
	require.Equal(t, "H1", SquareFromName("H1").String())
}

func TestReplacePieces(t *testing.T) {
	board, err := ParseFen(FenDefault)
	require.NoError(t, err)
	board.ReplacePieces(ColorBlack, TypePawn, TypeQueen)
	expected := "rnbqkbnr/qqqqqqqq/8/8/8/8/PPPPPPPP/RNBQKBNR"
	fen := EncodeFen(board)
	require.Equal(t, expected, fen)
}
