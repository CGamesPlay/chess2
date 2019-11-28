package chess2

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSingleStepMask(t *testing.T) {
	cases := []struct {
		pair    string
		squares []string
	}{
		{
			pair:    "b5d5",
			squares: []string{"c5"},
		},
		{
			pair:    "b3d5",
			squares: []string{"b4", "c4", "c3"},
		},
		{
			pair:    "d3d5",
			squares: []string{"d4"},
		},
		{
			pair:    "f3d5",
			squares: []string{"e4", "f4", "e3"},
		},
		{
			pair:    "f5d5",
			squares: []string{"e5"},
		},
		{
			pair:    "f7d5",
			squares: []string{"e7", "e6", "f6"},
		},
		{
			pair:    "d7d5",
			squares: []string{"d6"},
		},
		{
			pair:    "b7d5",
			squares: []string{"c7", "b6", "c6"},
		},
	}
	for _, config := range cases {
		move, err := ParseUci(config.pair)
		require.NoError(t, err, "UCI: %s", config.pair)
		mask := singleStepMask(move.From, move.To.Mask())
		results := make([]string, 0, 3)
		eachSquareInMask(mask, func(sq Square) {
			results = append(results, sq.String())
		})
		require.Equal(t, config.squares, results, "UCI: %s", config.pair)
	}
}

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

type IsInCheckTest struct {
	epd     string
	color   Color
	inCheck bool
}

func TestIsInCheck(t *testing.T) {
	cases := map[string]IsInCheckTest{
		"not in check": {
			epd:     "4k3/8/8/8/8/8/8/4K3 w - - 0 1 cc 33",
			color:   ColorWhite,
			inCheck: false,
		},
		"regular threat": {
			epd:     "4k3/8/8/8/1b6/8/8/4K3 w - - 0 1 cc 33",
			color:   ColorWhite,
			inCheck: true,
		},
		"rampage would kill opponent": {
			epd:     "4k3/8/8/8/8/8/8/2rK4 b - - 0 1 ca 33",
			color:   ColorWhite,
			inCheck: true,
		},
		"rampage would kill own king": {
			epd:     "8/8/8/8/8/8/8/2rK1k2 b - - 0 1 ca 33",
			color:   ColorWhite,
			inCheck: false,
		},
	}
	for name, config := range cases {
		if name == "rampage would kill own king" {
			if testing.Verbose() {
				t.Logf("skipping case %s\n", name)
			}
			continue
		}
		game, err := ParseEpd(config.epd)
		require.NoError(t, err, "EPD: %s  Name: %s", config.epd, name)
		inCheck := game.IsInCheck(config.color)
		assert.Equal(t, inCheck, config.inCheck, "Case: %s", name)
	}
}

type ValidateMoveTest struct {
	epd  string
	move string
	err  error
}

func TestValidatePseudoLegalMove(t *testing.T) {
	cases := map[string]ValidateMoveTest{
		"illegal pass": {
			epd:  "4k3/8/8/8/8/8/8/4K3 b - - 0 1 nn 33",
			move: "0000",
			err:  IllegalPassError,
		},
		"legal pass": {
			epd:  "4k3/8/8/8/8/8/8/4K3 K - - 0 1 kn 33",
			move: "0000",
		},
		"illegal drop": {
			epd:  "4k3/8/8/8/8/8/8/4K3 b - - 0 1 nn 33",
			move: "Q@a8",
			err:  IllegalDropError,
		},
		"out of turn": {
			epd:  "4k3/8/8/8/8/8/8/4K3 b - - 0 1 nn 33",
			move: "e1e2",
			err:  NotMovablePieceError,
		},
		"missing piece": {
			epd:  "4k3/8/8/8/8/8/8/4K3 b - - 0 1 nn 33",
			move: "d2d4",
			err:  NotMovablePieceError,
		},
		"non-king during king turn": {
			epd:  "4k3/8/8/8/8/8/8/3QK3 K - - 0 1 kk 33",
			move: "d1d2",
			err:  IllegalKingTurnError,
		},
		"valid promotion": {
			epd:  "4k3/P7/8/8/8/8/8/4K3 w - - 0 1 cc 33",
			move: "a7a8q",
		},
		"promotion before last rank": {
			epd:  "4k3/8/P7/8/8/8/8/4K3 w - - 0 1 cc 33",
			move: "a6a7q",
			err:  IllegalPromotionError,
		},
		"promotion of other piece": {
			epd:  "4k3/8/8/8/8/8/8/R3K3 w - - 0 1 cc 33",
			move: "a1a8q",
			err:  IllegalPromotionError,
		},
		"missing promotion": {
			epd:  "4k3/P7/8/8/8/8/8/4K3 w - - 0 1 cc 33",
			move: "a7a8",
			err:  IllegalPromotionError,
		},
		"basic king": {
			epd:  "4k3/8/8/8/8/8/8/4K3 w - - 0 1 nn 33",
			move: "e1e2",
		},
		"illegal capture of own piece": {
			epd:  "2B1k3/8/8/8/8/8/8/2R1K3 w - - 0 1 cc 33",
			move: "c1c8",
			err:  IllegalCaptureError,
		},
		"legal capture of own piece": {
			epd:  "4k3/8/8/8/8/8/3P4/1N2K3 w - - 0 1 ac 33",
			move: "b1d2",
		},
		"illegal capture of own king": {
			epd:  "4k3/8/8/8/8/8/2N5/4K3 w - - 0 1 ac 33",
			move: "c2e1",
			err:  IllegalCaptureError,
		},
		"illegal capture of nemesis queen": {
			epd:  "4k3/8/8/8/3q4/8/8/3RK3 w - - 0 1 cn 33",
			move: "d1d4",
			err:  IllegalCaptureError,
		},
		"legal capture of nemesis queen": {
			epd:  "4k3/8/8/8/8/8/8/3qK3 w - - 0 1 cn 33",
			move: "e1d1",
		},
		"illegal capture of nemesis rook": {
			epd:  "4k3/8/8/8/8/8/8/3rK3 w - - 0 1 cn 33",
			move: "e1d1",
			err:  IllegalCaptureError,
		},
		"illegal capture of elephant": {
			epd:  "3rk3/8/8/8/8/8/8/3RK3 w - - 0 1 ca 33",
			move: "d1d8",
			err:  IllegalCaptureError,
		},
		"legal capture of elephant": {
			epd:  "3rk3/8/1B6/8/8/8/8/4K3 w - - 0 1 ca 33",
			move: "b6d8",
		},
		"pawn double move": {
			epd:  "4k3/8/8/8/8/8/4P3/4K3 w - - 0 1 cc 33",
			move: "e2e4",
		},
		"illegal pawn double move": {
			epd:  "4k3/8/8/8/8/4P3/8/4K3 w - - 0 1 cc 33",
			move: "e3e5",
			err:  UnreachableSquareError,
		},
		"pawn single move": {
			epd:  "4k3/8/8/8/4p3/8/4P3/4K3 w - - 0 1 cc 33",
			move: "e2e3",
		},
		"illegal pawn single move": {
			epd:  "4k3/8/8/8/8/4p3/4P3/4K3 w - - 0 1 cc 33",
			move: "e2e3",
			err:  IllegalCaptureError,
		},
		"pawn move backwards": {
			epd:  "4k3/8/8/8/8/8/4P3/4K3 w - - 0 1 cc 33",
			move: "e2e1",
			err:  UnreachableSquareError,
		},
		"nemesis pawn move backwards": {
			epd:  "8/2P5/4k3/8/8/8/8/4K3 w KQkq - 0 1 nc 33",
			move: "c7c6",
		},
		"illegal nemesis move": {
			epd:  "4k3/8/8/8/8/8/3P4/4K3 w - - 0 1 nc 33",
			move: "d2c2",
			err:  UnreachableSquareError,
		},
		"legal nemesis move": {
			epd:  "4k3/8/8/8/8/8/3P4/4K3 w - - 0 1 nc 33",
			move: "d2e2",
		},
		"illegal capturing nemesis move": {
			epd:  "4k3/8/8/8/8/8/3Pp3/4K3 w - - 0 1 nc 33",
			move: "d2e2",
			err:  IllegalCaptureError,
		},
		"pawn double move capture": {
			epd:  "4k3/8/8/8/4p3/8/4P3/4K3 w - - 0 1 cc 33",
			move: "e2e4",
			err:  IllegalCaptureError,
		},
		"pawn capture wraparound": {
			epd:  "4k3/7p/8/P7/8/8/8/4K3 w - - 0 1 cc 33",
			move: "a5h7",
			err:  UnreachableSquareError,
		},
		"pawn capture en passant": {
			epd:  "4k3/8/8/2Pp4/8/8/8/4K3 w - d6 0 1 cc 33",
			move: "c5d6",
		},
		"illegal sliding attack": {
			epd:  "4k3/8/8/4r3/8/8/8/4K3 w - - 0 1 cc 33",
			move: "e1e5",
			err:  UnreachableSquareError,
		},
		"whirlwind of illegal army": {
			epd:  "4k3/8/8/8/8/8/8/4K3 w - - 0 1 cc 33",
			move: "e1e1",
			err:  IllegalWhirlwindAttackError,
		},
		"whirlwind of illegal piece": {
			epd:  "4k3/8/8/8/8/3P4/8/3KK3 w - - 0 1 kc 33",
			move: "d3d3",
			err:  IllegalWhirlwindAttackError,
		},
		"illegal whirlwind during normal turn": {
			epd:  "4k3/8/8/8/8/8/8/2K1K3 w - - 0 1 kc 33",
			move: "c1c1",
			err:  IllegalWhirlwindAttackError,
		},
		"whirlwind capturing own king": {
			epd:  "4k3/8/8/8/8/8/8/3KK3 K - - 0 1 kc 33",
			move: "d1d1",
			err:  IllegalCaptureError,
		},
		"legal whirlwind": {
			epd:  "4k3/8/8/8/8/8/8/2K1K3 K - - 0 1 kc 33",
			move: "c1c1",
		},
		"illegal castle after moving pieces": {
			epd:  "4k3/8/8/8/8/8/8/R3K3 w kq - 0 1 cn 33",
			move: "e1c1",
			err:  IllegalCastleError,
		},
		"illegal castle due to blockage": {
			epd:  "4k3/8/8/8/8/8/8/R2QK3 w KQkq - 0 1 cn 33",
			move: "e1c1",
			err:  IllegalCastleError,
		},
		"castle by illegal army": {
			epd:  "4k3/8/8/8/8/8/8/R3K3 w KQkq - 0 1 nn 33",
			move: "e1c1",
			err:  IllegalCastleError,
		},
		"illegal castle through check": {
			epd:  "3rk3/8/8/8/8/8/8/R3K3 w Q - 0 1 cc 33",
			move: "e1c1",
			err:  IllegalCastleError,
		},
		"legal castle": {
			epd:  "4k3/8/8/8/8/8/8/R3K3 w KQkq - 0 1 cn 33",
			move: "e1c1",
		},
		"elephant rampage": {
			epd:  "4k3/8/8/8/RnpP4/8/8/4K3 w KQkq - 0 1 ac 33",
			move: "a4d4",
		},
		"rampage at left": {
			epd:  "4k3/8/8/8/8/8/1pR5/4K3 w - - 0 1 ac 33",
			move: "c2a2",
		},
		"rampage stops short": {
			epd:  "4k3/8/8/8/8/2p5/2R5/4K3 w - - 0 1 ac 33",
			move: "c2c3",
			err:  IllegalRampageError,
		},
		"rampage kills own king": {
			epd:  "4k3/8/8/8/8/8/8/4K2R w K - 0 1 ac 33",
			move: "h1e1",
			err:  IllegalRampageError,
		},
		"illegal elephant capture": {
			epd:  "4k3/8/8/8/Rn6/8/8/4K3 w KQkq - 0 1 ac 33",
			move: "a4b4",
			err:  IllegalRampageError,
		},
	}
	for name, config := range cases {
		game, err := ParseEpd(config.epd)
		require.NoError(t, err, "EPD: %s  Name: %s", config.epd, name)
		move, err := ParseUci(config.move)
		require.NoError(t, err, "Move: %s  Name: %s", config.move, name)
		err = game.ValidatePseudoLegalMove(move)
		if config.err != nil {
			assert.EqualError(t, err, config.err.Error(), "Case: %s", name)
		} else {
			assert.NoError(t, err, "Case: %s", name)
		}
	}
}

func TestValidateLegalMove(t *testing.T) {
	cases := map[string]ValidateMoveTest{
		"legal move": {
			epd:  "4k3/4p3/8/8/8/8/8/4K3 b - - 0 1 cc 33",
			move: "e7e5",
			err:  nil,
		},
		"not pseudo legal": {
			epd:  "4k3/8/8/8/8/8/8/4K3 b - - 0 1 nn 33",
			move: "0000",
			err:  IllegalPassError,
		},
		"into check": {
			epd:  "3rk3/8/8/8/8/8/8/4K3 w - - 0 1 nn 33",
			move: "e1d1",
			err:  MoveIntoCheckError,
		},
	}
	for name, config := range cases {
		game, err := ParseEpd(config.epd)
		require.NoError(t, err, "EPD: %s  Name: %s", config.epd, name)
		move, err := ParseUci(config.move)
		require.NoError(t, err, "Move: %s  Name: %s", config.move, name)
		err = game.ValidateLegalMove(move)
		if config.err != nil {
			assert.EqualError(t, err, config.err.Error(), "Case: %s", name)
		} else {
			assert.NoError(t, err, "Case: %s", name)
		}
	}
}

func TestValidateDuels(t *testing.T) {
	cases := map[string]ValidateMoveTest{
		"legal duel": {
			epd:  "4k3/8/8/4p3/3P4/8/8/4K3 w - - 0 1 cc 33",
			move: "d4e5:22",
		},
		"duel without capture": {
			epd:  "4k3/8/8/4p3/3P4/8/8/4K3 w - - 0 1 cc 33",
			move: "d4d5:00+",
			err:  TooManyDuelsError,
		},
		"challenge too expensive": {
			epd:  "4k3/8/8/4p3/3P4/8/8/4K3 w - - 0 1 cc 13",
			move: "d4e5:22",
			err:  NotEnoughStonesError,
		},
		"response too expensive": {
			epd:  "4k3/8/8/4p3/3P4/8/8/4K3 w - - 0 1 cc 31",
			move: "d4e5:22",
			err:  NotEnoughStonesError,
		},
		"cannot pay to initiate": {
			epd:  "4k3/8/8/4b3/3P4/8/8/4K3 w - - 0 1 cc 03",
			move: "d4e5:12",
			err:  NotEnoughStonesError,
		},
		"king is dueling": {
			epd:  "4k3/8/8/8/8/8/4p3/4K3 w - - 0 1 cc 33",
			move: "e1e2:22",
			err:  NotDuelableError,
		},
	}
	for name, config := range cases {
		game, err := ParseEpd(config.epd)
		require.NoError(t, err, "EPD: %s  Name: %s", config.epd, name)
		move, err := ParseUci(config.move)
		require.NoError(t, err, "Move: %s  Name: %s", config.move, name)
		err = game.ValidateDuels(move)
		if config.err != nil {
			assert.EqualError(t, err, config.err.Error(), "Case: %s", name)
		} else {
			assert.NoError(t, err, "Case: %s", name)
		}
	}
}

func TestApplyMove(t *testing.T) {
	cases := map[string]struct {
		before string
		move   string
		after  string
	}{
		"drop": {
			before: "8/8/8/8/8/8/8/8 w KQkq - 0 1 nk 33",
			move:   "K@e1",
			after:  "8/8/8/8/8/8/8/4K3 w KQkq - 0 1 nk 33",
		},
		"pass without two kings": {
			before: "8/8/8/8/8/8/8/8 w KQkq - 0 1 nk 33",
			move:   "0000",
			after:  "8/8/8/8/8/8/8/8 b KQkq - 1 1 nk 33",
		},
		"pass with two kings": {
			before: "8/8/8/8/8/8/8/8 w KQkq - 0 1 ka 33",
			move:   "0000",
			after:  "8/8/8/8/8/8/8/8 K KQkq - 1 1 ka 33",
		},
		"pass on king-turn": {
			before: "8/8/8/8/8/8/8/8 K KQkq - 0 1 ka 33",
			move:   "0000",
			after:  "8/8/8/8/8/8/8/8 b KQkq - 0 1 ka 33",
		},
		"normal move": {
			before: "4k3/8/8/4p3/8/8/8/4K3 w KQkq e6 0 1 cc 33",
			move:   "e1e2",
			after:  "4k3/8/8/4p3/8/8/4K3/8 b kq - 1 1 cc 33",
		},
		"capturing move": {
			before: "4k3/8/8/8/8/8/4b3/4K3 w - - 0 1 cc 33",
			move:   "e1e2",
			after:  "4k3/8/8/8/8/8/4K3/8 b - - 0 1 cc 33",
		},
		"capturing pawn": {
			before: "4k3/8/8/8/8/8/4p3/4K3 w KQkq - 0 1 cc 33",
			move:   "e1e2",
			after:  "4k3/8/8/8/8/8/4K3/8 b kq - 0 1 cc 43",
		},
		"capture en-passant": {
			before: "4k3/8/8/3Pp3/8/8/8/4K3 w KQkq e6 0 1 cc 33",
			move:   "d5e6",
			after:  "4k3/8/4P3/8/8/8/8/4K3 b KQkq - 0 1 cc 43",
		},
		"non-pawn to ep square": {
			before: "4k3/8/8/3Bp3/8/8/8/4K3 w - e6 0 1 cc 33",
			move:   "d5e6",
			after:  "4k3/8/4B3/4p3/8/8/8/4K3 b - - 1 1 cc 33",
		},
		"nemesis move to ep square": {
			before: "4k3/8/3P4/4p3/8/8/8/4K3 w - - 0 1 nc 33",
			move:   "d6e6",
			after:  "4k3/8/4P3/4p3/8/8/8/4K3 b - - 1 1 nc 33",
		},
		"advance fullmove number": {
			before: "8/8/8/8/8/8/8/8 b KQkq - 0 1 cc 33",
			move:   "0000",
			after:  "8/8/8/8/8/8/8/8 w KQkq - 1 2 cc 33",
		},
		"set ep square": {
			before: "4k3/8/8/8/8/8/4P3/4K3 w KQkq - 0 1 cc 33",
			move:   "e2e4",
			after:  "4k3/8/8/8/4P3/8/8/4K3 b KQkq e3 0 1 cc 33",
		},
		"preserve ep square on king-turn": {
			before: "4k3/8/8/4p3/8/8/8/4K3 k - e6 0 1 ck 33",
			move:   "e8e7",
			after:  "8/4k3/8/4p3/8/8/8/4K3 w - e6 0 2 ck 33",
		},
		"clear castling rights on rook move": {
			before: "4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1 cc 33",
			move:   "a1a8",
			after:  "R3k3/8/8/8/8/8/8/4K2R b K - 1 1 cc 33",
		},
		"clear castling rights on king move": {
			before: "4k3/8/8/8/8/8/8/R3K2R w KQ - 0 1 cc 33",
			move:   "e1e2",
			after:  "4k3/8/8/8/8/8/4K3/R6R b - - 1 1 cc 33",
		},
		"castling queenside": {
			before: "4k3/8/8/8/8/8/8/R3K3 w KQ - 0 1 cc 33",
			move:   "e1c1",
			after:  "4k3/8/8/8/8/8/8/2KR4 b - - 1 1 cc 33",
		},
		"castling kingside": {
			before: "4k2r/8/8/8/8/8/8/4K3 b kq - 0 1 cc 33",
			move:   "e8g8",
			after:  "5rk1/8/8/8/8/8/8/4K3 w - - 1 2 cc 33",
		},
		"animals bishop noncapture": {
			before: "4k3/8/8/8/8/5p2/8/3BK3 w - - 0 1 ac 33",
			move:   "d1e2",
			after:  "4k3/8/8/8/8/5p2/4B3/4K3 b - - 1 1 ac 33",
		},
		"animals bishop capture": {
			before: "4k3/8/8/8/8/5p2/8/3BK3 w - - 0 1 ac 33",
			move:   "d1f3",
			after:  "4k3/8/8/8/8/8/8/3BK3 b - - 0 1 ac 43",
		},
		"whirlwind attack": {
			before: "4k3/8/8/2Prp3/2bKn3/2pBP3/8/4K3 w - - 0 1 kr 33",
			move:   "d4d4",
			after:  "4k3/8/8/3r4/3K4/8/8/4K3 K - - 0 1 kr 53",
		},
		"rampage": {
			before: "4k3/Rppp4/8/8/8/8/8/4K3 w - - 0 1 ac 33",
			move:   "a7d7",
			after:  "4k3/3R4/8/8/8/8/8/4K3 b - - 0 1 ac 63",
		},
		"winning challenge": {
			before: "4k3/8/8/4p3/3P4/8/8/4K3 w - - 0 1 cc 33",
			move:   "d4e5:10+",
			after:  "4k3/8/8/8/8/8/8/4K3 b - - 0 1 cc 42",
		},
		"losing challenge": {
			before: "4k3/8/8/4p3/3P4/8/8/4K3 w - - 0 1 cc 33",
			move:   "d4e5:11",
			after:  "4k3/8/8/4P3/8/8/8/4K3 b - - 0 1 cc 32",
		},
		"pay to duel": {
			before: "4k3/8/8/4b3/3P4/8/8/4K3 w - - 0 1 cc 33",
			move:   "d4e5:22",
			after:  "4k3/8/8/4P3/8/8/8/4K3 b - - 0 1 cc 01",
		},
		"call bluff, gain": {
			before: "4k3/8/8/4p3/3P4/8/8/4K3 w - - 0 1 cc 33",
			move:   "d4e5:00+",
			after:  "4k3/8/8/4P3/8/8/8/4K3 b - - 0 1 cc 53",
		},
		"call bluff, lose": {
			before: "4k3/8/8/4p3/3P4/8/8/4K3 w - - 0 1 cc 33",
			move:   "d4e5:00-",
			after:  "4k3/8/8/4P3/8/8/8/4K3 b - - 0 1 cc 42",
		},
		"promotion": {
			before: "4k3/1P6/8/8/8/8/8/4K3 w - - 0 1 cc 33",
			move:   "b7b8q",
			after:  "1Q2k3/8/8/8/8/8/8/4K3 b - - 0 1 cc 33",
		},
	}
	for name, config := range cases {
		before, err := ParseEpd(config.before)
		require.NoError(t, err, "EPD: %s  Name: %s", config.before, name)
		move, err := ParseUci(config.move)
		require.NoError(t, err, "Move: %s  Name: %s", config.move, name)
		after := before.ApplyMove(move)
		result := EncodeEpd(after)
		assert.Equal(t, config.after, result, "Case: %s", name)
	}
}
