package chess2

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEncodeEpd(t *testing.T) {
	game := GameFromArmies(ArmyNemesis, ArmyEmpowered)
	epd := EncodeEpd(game)
	expected := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1 ne 33"
	require.Equal(t, expected, epd)
}

func TestParseEpd(t *testing.T) {
	epd := "rnbqkbnr/pppp1ppp/8/4p3/3P4/8/PPP1PPPP/RNBQKBNR w KQkq - 0 2 cc 33"
	game, err := ParseEpd(epd)
	require.NoError(t, err)
	result := EncodeEpd(game)
	require.Equal(t, epd, result)
}
