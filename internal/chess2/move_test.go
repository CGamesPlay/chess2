package chess2

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMoveUci(t *testing.T) {
	cases := []string{
		"d2d4",
		"0000",
		"p@b3",
		"e7a7:10+",
		"a4d4:21",
		"a8d8::10+",
		"a7a8Q:22",
	}
	for _, uci := range cases {
		move, err := ParseUci(uci)
		require.NoError(t, err, "UCI: %s", uci)
		assert.Equal(t, uci, move.String())
	}
}
