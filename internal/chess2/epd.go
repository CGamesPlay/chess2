package chess2

import (
	"fmt"
	"strings"
)

// EPD is: fen tomove castle epsquare hc fm armies stones operations

func colorSymbol(c Color) rune {
	if c == ColorWhite {
		return 'w'
	}
	return 'b'
}

func castlingRights(sb *strings.Builder, rights uint64) {
	if rights == 0 {
		sb.WriteRune('-')
	}
	if rights&castleWhiteKingside != 0 {
		sb.WriteRune('K')
	}
	if rights&castleWhiteQueenside != 0 {
		sb.WriteRune('Q')
	}
	if rights&castleBlackKingside != 0 {
		sb.WriteRune('k')
	}
	if rights&castleBlackQueenside != 0 {
		sb.WriteRune('q')
	}
}

var armySymbols = map[Army]rune{
	ArmyClassic:   'c',
	ArmyNemesis:   'n',
	ArmyEmpowered: 'e',
	ArmyReaper:    'r',
	ArmyTwoKings:  'k',
	ArmyAnimals:   'a',
}

// EncodeEpd returns the EPD of the given game object
func EncodeEpd(game Game) string {
	var sb strings.Builder
	sb.WriteString(EncodeFen(game.board))
	sb.WriteRune(' ')
	sb.WriteRune(colorSymbol(game.toMove))
	sb.WriteRune(' ')
	castlingRights(&sb, game.castlingRights)
	sb.WriteRune(' ')
	if game.epSquare != InvalidSquare {
		sb.WriteString(strings.ToLower(game.epSquare.String()))
	} else {
		sb.WriteRune('-')
	}
	sb.WriteString(fmt.Sprintf(
		" %d %d %c%c %d%d",
		game.halfmoveClock,
		game.fullmoveNumber+1,
		armySymbols[game.armies[0]],
		armySymbols[game.armies[1]],
		game.stones[0],
		game.stones[1],
	))
	return sb.String()
}

// ParseEpd parses an EPD string and returns a Game object.
func ParseEpd(epd string) Game {
	result := Game{}
	return result
}
