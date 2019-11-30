package chess2

import (
	"fmt"
	"strings"
)

func toMoveSymbol(c Color, kingTurn bool) rune {
	if c == ColorWhite {
		if kingTurn {
			return 'K'
		}
		return 'w'
	}
	if kingTurn {
		return 'k'
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

var (
	armyToSymbol = map[Army]rune{
		ArmyClassic:   'c',
		ArmyNemesis:   'n',
		ArmyEmpowered: 'e',
		ArmyReaper:    'r',
		ArmyTwoKings:  'k',
		ArmyAnimals:   'a',
	}

	symbolToArmy = map[rune]Army{
		'c': ArmyClassic,
		'n': ArmyNemesis,
		'e': ArmyEmpowered,
		'r': ArmyReaper,
		'k': ArmyTwoKings,
		'a': ArmyAnimals,
	}

	castlingCodes = map[rune]uint64{
		'K': castleWhiteKingside,
		'Q': castleWhiteQueenside,
		'k': castleBlackKingside,
		'q': castleBlackQueenside,
	}
)

// FindArmySymbol takes the EPD symbol for an army and converts it to the proper
// Army value.
func FindArmySymbol(symbol rune) (Army, bool) {
	army, found := symbolToArmy[symbol]
	return army, found
}

// EncodeEpd returns the EPD of the given game object
func EncodeEpd(game Game) string {
	var sb strings.Builder
	sb.WriteString(EncodeFen(game.board))
	sb.WriteRune(' ')
	sb.WriteRune(toMoveSymbol(game.toMove, game.kingTurn))
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
		armyToSymbol[game.armies[0]],
		armyToSymbol[game.armies[1]],
		game.stones[0],
		game.stones[1],
	))
	return sb.String()
}

// ParseEpd parses an EPD string and returns a Game object.
func ParseEpd(epd string) (Game, error) {
	game := Game{}
	// Split the EPD into components
	var fenStr, castleStr, epStr string
	var toMoveRune rune
	var armyRunes, stoneRunes [2]rune
	// EPD is: fen tomove castle epsquare hc fm armies stones operations
	num, err := fmt.Sscanf(
		epd, "%s %c %s %s %d %d %c%c %c%c",
		&fenStr, &toMoveRune, &castleStr, &epStr,
		&game.halfmoveClock, &game.fullmoveNumber,
		&armyRunes[0], &armyRunes[1],
		&stoneRunes[0], &stoneRunes[1],
	)
	if err != nil || num != 10 {
		return Game{}, ParseError("EPD invalid")
	}
	board, err := ParseFen(fenStr)
	if err != nil {
		return Game{}, err
	}
	game.board = board

	switch toMoveRune {
	case 'w':
		game.toMove = ColorWhite
	case 'b':
		game.toMove = ColorBlack
	case 'K':
		game.toMove = ColorWhite
		game.kingTurn = true
	case 'k':
		game.toMove = ColorBlack
		game.kingTurn = true
	default:
		return Game{}, ParseError("EPD has invalid to-move")
	}

	for symbol, value := range castlingCodes {
		if strings.IndexRune(castleStr, symbol) != -1 {
			game.castlingRights |= value
		}
	}

	if epStr == "-" {
		game.epSquare = InvalidSquare
	} else {
		game.epSquare = SquareFromName(epStr)
		if game.epSquare == InvalidSquare || game.epSquare.Y() != 2 && game.epSquare.Y() != 5 {
			return Game{}, ParseError("EPD has invalid en passant square")
		}
	}

	// Adjust to be 0-based
	game.fullmoveNumber--

	for i, symbol := range armyRunes {
		army, found := FindArmySymbol(symbol)
		if !found {
			return Game{}, ParseError("EPD has invalid armies")
		}
		game.armies[i] = army
	}

	for i, stoneRune := range stoneRunes {
		if stoneRune < '0' || stoneRune > '6' {
			return Game{}, ParseError("EPD has invalid stones")
		}
		game.stones[i] = int(stoneRune - '0')
	}

	// Finally, some sanity checking
	if game.kingTurn {
		if game.armies[ColorIdx(game.toMove)] != ArmyTwoKings {
			return Game{}, ParseError("King turn for army other than two kings")
		}
	}
	game.updateGameState()

	return game, nil
}
