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
