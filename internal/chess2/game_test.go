package chess2

import (
	"github.com/stretchr/testify/assert"
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

type AttackMaskTest struct {
	epd     string
	squares []string
	concern string
}

func TestAttackMask(t *testing.T) {
	cases := map[string]AttackMaskTest{
		"basic king": {
			epd:     "8/8/8/8/8/2K5/8/8 w - - 0 1 cc 33",
			squares: []string{"b4", "c4", "d4", "b3", "d3", "b2", "c2", "d2"},
		},
		"basic queen": {
			epd: "8/8/8/8/8/8/4Q3/8 w - - 0 1 cc 33",
			squares: []string{
				"e8", "e7", "a6", "e6", "b5", "e5", "h5",
				"c4", "e4", "g4", "d3", "e3", "f3",
				"a2", "b2", "c2", "d2", "f2", "g2", "h2",
				"d1", "e1", "f1",
			},
		},
		"basic bishop": {
			epd: "8/8/8/8/8/8/4B3/8 w - - 0 1 cc 33",
			squares: []string{
				"a6", "b5", "h5", "c4", "g4", "d3", "f3", "d1", "f1",
			},
		},
		"basic knight": {
			epd:     "8/8/8/8/8/4N3/8/8 w - - 0 1 cc 33",
			squares: []string{"d5", "f5", "c4", "g4", "c2", "g2", "d1", "f1"},
		},
		"basic rook": {
			epd: "8/8/8/8/8/8/8/R7 w - - 0 1 cc 33",
			squares: []string{
				"a8", "a7", "a6", "a5", "a4", "a3", "a2",
				"b1", "c1", "d1", "e1", "f1", "g1", "h1",
			},
		},
		"basic pawn": {
			epd:     "8/p3p2p/8/8/8/8/P2P3P/8 w - - 0 1 cn 33",
			squares: []string{"b6", "d6", "f6", "g6", "b3", "c3", "e3", "g3"},
		},
		"nemesis queen": {
			epd: "3k4/8/8/8/5p2/2p5/3Q4/8 w - - 0 1 nc 33",
			squares: []string{
				"d8", "d7", "d6", "d5", "d4", "d3", "e3",
				"a2", "b2", "c2", "e2", "f2", "g2", "h2",
				"c1", "d1", "e1",
			},
			concern: "d2",
		},
		"empowered rook-knight": {
			epd: "8/8/8/3rn3/8/8/8/8 w - - 0 1 ce 33",
			squares: []string{
				"d8", "e8", "c7", "d7", "e7", "f7",
				"b6", "c6", "d6", "e6", "f6", "g6",
				"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
				"b4", "c4", "d4", "e4", "f4", "g4",
				"c3", "d3", "e3", "f3", "d2", "e2", "d1", "e1",
			},
		},
		"empowered rook-bishop": {
			epd: "8/8/8/3RB3/8/8/8/8 w - - 0 1 ec 33",
			squares: []string{
				"a8", "b8", "d8", "e8", "g8", "h8",
				"b7", "c7", "d7", "e7", "f7", "g7",
				"c6", "d6", "e6", "f6",
				"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
				"c4", "d4", "e4", "f4",
				"b3", "c3", "d3", "e3", "f3", "g3",
				"a2", "b2", "d2", "e2", "g2", "h2",
				"a1", "d1", "e1", "h1",
			},
		},
		"empowered knight-bishop": {
			epd: "8/8/3r4/2bp4/2n5/8/8/8 w - - 0 1 ce 33",
			squares: []string{
				"a7", "b7", "d7", "a6", "b6", "d6", "e6",
				"a4", "b4", "d4", "e4", "a3", "b3", "d3", "e3",
				"f2", "g1",
			},
			concern: "c5",
		},
		"reaper queen": {
			epd: "rnbq1bnr/pppkpppp/R2p3R/8/8/8/PPPPPPPP/1NBQKBN1 w KQkq - 0 1 rc 33",
			squares: []string{
				"a7", "b7", "c7", "e7", "f7", "g7", "h7",
				"a6", "b6", "c6", "d6", "e6", "f6", "g6", "h6",
				"a5", "b5", "c5", "d5", "e5", "f5", "g5", "h5",
				"a4", "b4", "c4", "d4", "e4", "f4", "g4", "h4",
				"a3", "b3", "c3", "d3", "e3", "f3", "g3", "h3",
				"a2", "b2", "c2", "d2", "e2", "f2", "g2", "h2",
				"a1", "b1", "c1", "d1", "e1", "f1", "g1", "h1",
			},
			concern: "d1",
		},
		"reaper rook": {
			epd: "r1bqkb1r/pppp1ppp/2n2n2/4N3/4P3/2N5/PPPP1PPP/R1BQKB1R b KQkq - 0 1 rc 33",
			squares: []string{
				"b8", "g8", "e7", "a6", "b6", "d6", "e6", "g6", "h6",
				"a5", "b5", "c5", "d5", "f5", "g5", "h5",
				"a4", "b4", "c4", "d4", "f4", "g4", "h4",
				"a3", "b3", "d3", "e3", "f3", "g3", "h3", "e2", "b1", "g1",
			},
			concern: "a1",
		},
		"animals queen": {
			epd: "8/8/8/8/8/8/8/Q7 w - - 0 1 ac 33",
			squares: []string{
				"a8", "a7", "a6", "a5", "a4", "a3", "b3", "a2", "c2",
				"b1", "c1", "d1", "e1", "f1", "g1", "h1",
			},
		},
		"animals bishop": {
			epd: "8/8/8/8/3B4/8/8/8 w - - 0 1 ac 33",
			squares: []string{
				"b6", "f6", "c5", "e5", "c3", "e3", "b2", "f2",
			},
		},
		"animals rook": {
			epd: "8/8/8/8/3R4/8/8/8 w - - 0 1 ac 33",
			squares: []string{
				"d7", "d6", "d5", "a4", "b4", "c4",
				"e4", "f4", "g4", "d3", "d2", "d1",
			},
		},
	}
	for name, config := range cases {
		game, err := ParseEpd(config.epd)
		require.NoError(t, err, "EPD: %s  Name: %s", config.epd, name)
		fromMask := MaskFull
		if config.concern != "" {
			fromMask = SquareFromName(config.concern).Mask()
		}
		mask := uint64(0)
		eachSquareInMask(fromMask, func(from Square) {
			mask |= game.attackMask(from)
		})
		sqNames := make([]string, 0, 64)
		eachSquareInMask(mask, func(to Square) {
			sqNames = append(sqNames, to.String())
		})
		assert.Equal(t, config.squares, sqNames, "EPD: %s  Name: %s", config.epd, name)
	}
}
