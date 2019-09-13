package chess2

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGameFromArmiesTwoKings(t *testing.T) {
	var (
		game Game
		fen  string
	)
	game = GameFromArmies(ArmyTwoKings, ArmyNemesis)
	fen = EncodeFen(game.board)
	require.Equal(t, "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBKKBNR", fen)
	game = GameFromArmies(ArmyNemesis, ArmyTwoKings)
	fen = EncodeFen(game.board)
	require.Equal(t, "rnbkkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR", fen)
}
