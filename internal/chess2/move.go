package chess2

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	challengeMask = 0x03
	responseMask  = 0x0c
	gainMask      = 0x10
	startedMask   = 0x20
	completeMask  = 0x40
)

// A Duel is a packed representation of a duel outcome. A duel is initiated by
// the defender (the one who's piece was captured), and the attacker responds to
// the duel.
type Duel struct {
	repr uint8
}

// ParseDuel takes a duel string like "20+" and returns Duel object.
func ParseDuel(str string) (Duel, error) {
	var repr uint8
	if len(str) > 0 {
		challenge := uint8(str[0] - '0')
		if challenge < 0 || challenge > 2 {
			goto invalid
		}
		repr |= challenge | startedMask
	}
	if len(str) > 1 {
		response := uint8(str[1] - '0')
		if response < 0 || response > 2 {
			goto invalid
		}
		repr |= (response << 2) | completeMask
	}
	if len(str) > 2 {
		gain := str[2]
		if repr&responseMask != 0 || (gain != '+' && gain != '-') {
			goto invalid
		} else if gain == '+' {
			repr |= gainMask
		}
	}
	if len(str) > 3 {
		goto invalid
	}
	return Duel{repr: repr}, nil
invalid:
	return Duel{}, ParseError("Invalid duel UCI")
}

// NewDuel creates a complete Duel from scratch.
func NewDuel(challenge, response int, gain bool) Duel {
	if challenge < 0 || challenge > 2 || response < 0 || response > 2 {
		panic("Invalid bid")
	}
	repr := uint8(challenge)
	repr |= uint8(response) << 2
	if gain {
		repr |= gainMask
	}
	repr |= completeMask | startedMask
	return Duel{repr: repr}
}

// DuelWithChallenge creates an incomplete duel with the given challenge value.
func DuelWithChallenge(bid int) Duel {
	if bid < 0 || bid > 2 {
		panic("Invalid bid")
	}
	return Duel{repr: uint8(bid) | startedMask}
}

// DuelWithResponse creates a complete duel from the passed incomplete duel with
// the attacker's response included.
func DuelWithResponse(duel Duel, bid int, gain bool) Duel {
	if bid < 0 || bid > 2 {
		panic("Invalid bid")
	}
	// Clear out the old response if one happened to be in there
	duel.repr = uint8(duel.Challenge())
	duel.repr |= uint8(bid << 2)
	if gain {
		duel.repr |= gainMask
	}
	duel.repr |= completeMask | startedMask
	return duel
}

// Challenge returns the number of stones bid by the defender.
func (d Duel) Challenge() int {
	return int(d.repr & challengeMask)
}

// Response returns the number of stones bid by the attacker.
func (d Duel) Response() int {
	return int((d.repr & responseMask) >> 2)
}

// Gain indicates that the attacker has bid 0 and wishes to gain a stone, rather
// than having the opponent lose a stone.
func (d Duel) Gain() bool {
	return d.repr&gainMask != 0
}

// IsStarted indicates that the defender has initiated a duel. This is necessary
// to differentiate from the default value of a Duel.
func (d Duel) IsStarted() bool {
	return d.repr&startedMask != 0
}

// IsComplete indicates that the attacker has responded to the duel. This is
// necessary to differentiate from a 0 Response value.
func (d Duel) IsComplete() bool {
	return d.repr&completeMask != 0
}

func (d Duel) String() string {
	var sb strings.Builder
	if d.IsStarted() {
		sb.WriteRune('0' + rune(d.Challenge()))
	}
	if d.IsComplete() {
		sb.WriteRune('0' + rune(d.Response()))
		if d.Response() == 0 {
			gain := '-'
			if d.Gain() {
				gain = '+'
			}
			sb.WriteRune(gain)
		}
	}
	return sb.String()
}

// Partial regexps that we combine to form a giant matcher for UCI.
const (
	patSkippedDuel    = ":"
	patIncompleteDuel = ":[0-2]"
	patBluffCall      = ":[0-2]0[-+]"
	patNonbluff       = ":[0-2][1-2]"
	patCompleteDuel   = "(" + patBluffCall + "|" + patNonbluff + ")"
	patDuels          = "(" + patSkippedDuel + "|" + patCompleteDuel + "){0,2}(" + patIncompleteDuel + "|" + patCompleteDuel + ")?"
	patSquare         = "[a-hA-H][1-8]"
	patPiece          = "[kqbnrpKQBNRP]"
	patNormalMove     = "(" + patSquare + patSquare + patPiece + "?)" + patDuels
	patDropMove       = patPiece + "@" + patSquare
)

var (
	reNormalMove = regexp.MustCompile("^" + patNormalMove + "$")
	reDropMove   = regexp.MustCompile("^" + patDropMove + "$")
)

// Move represents a move, including the duels that resulted from it.
type Move struct {
	From, To Square
	Piece    Piece
	Duels    [3]Duel
}

// MovePass represents a pass move.
var MovePass Move = Move{From: InvalidSquare, To: InvalidSquare}

// ParseUci takes a UCI string and returns the Move it represents.
func ParseUci(uci string) (Move, error) {
	if uci == "0000" {
		return MovePass, nil
	} else if reNormalMove.MatchString(uci) {
		move := Move{}
		move.From = SquareFromName(uci[0:2])
		move.To = SquareFromName(uci[2:4])
		if len(uci) > 4 {
			duelStart := 5
			if uci[4] != ':' {
				piece, err := ParseFenPiece(rune(uci[4]))
				if err != nil {
					panic("Regexp matched invalid promotion")
				}
				move.Piece = piece
				duelStart++
			}
			duelNumber := 0
			for duelStart < len(uci) {
				duelEnd := strings.IndexRune(uci[duelStart:], ':')
				var s string
				if duelEnd != -1 {
					s = uci[duelStart : duelEnd+duelStart]
					duelStart += duelEnd + 1
				} else {
					s = uci[duelStart:]
					duelStart = len(uci)
				}
				duel, err := ParseDuel(s)
				if err != nil {
					panic("Regexp matched invalid duel")
				}
				move.Duels[duelNumber] = duel
				duelNumber++
				if duelNumber > 2 && duelStart != -1 {
					// This shouldn't happen because the regexp validates the
					// number of duels.
					panic("too many duels")
				}
			}
		}
		if move.From != InvalidSquare && move.To != InvalidSquare && move.Piece.Type() != TypePawn {
			return move, nil
		}
	} else if reDropMove.MatchString(uci) {
		move := Move{}
		move.From = InvalidSquare
		move.To = SquareFromName(uci[2:4])
		piece, err := ParseFenPiece(rune(uci[0]))
		move.Piece = piece
		if err == nil && move.To != InvalidSquare {
			return move, nil
		}
	}
	return Move{}, ParseError("Invalid UCI")
}

// IsPass returns true if the move is a pass move.
func (m Move) IsPass() bool {
	return m.From == InvalidSquare && m.To == InvalidSquare
}

// IsDrop returns true if the move is a drop move.
func (m Move) IsDrop() bool {
	return m.From == InvalidSquare && m.To != InvalidSquare
}

func (m Move) String() string {
	switch {
	case m.IsPass():
		return "0000"
	case m.IsDrop():
		return fmt.Sprintf("%c@%s", EncodeFenPiece(m.Piece), m.To)
	default:
		var sb strings.Builder
		sb.WriteString(m.From.String())
		sb.WriteString(m.To.String())
		if m.Piece != InvalidPiece {
			sb.WriteRune(EncodeFenPiece(m.Piece))
		}
		numDuels := 0
		for i, d := range m.Duels {
			if d.IsStarted() {
				numDuels = i + 1
			}
		}
		for i := 0; i < numDuels; i++ {
			sb.WriteRune(':')
			sb.WriteString(m.Duels[i].String())
		}
		return sb.String()
	}
}
