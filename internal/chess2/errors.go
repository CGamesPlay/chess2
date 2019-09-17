package chess2

// IllegalMoveError represents the error that was encountered when attempting to
// make a move.
type IllegalMoveError int

const (
	// GameOverError is any move on a completed game.
	GameOverError = IllegalMoveError(iota + 1)
	// IllegalDropError is any drop move.
	IllegalDropError
	// IllegalPassError is a pass outside of a king-turn.
	IllegalPassError
	// NotMovablePieceError is a move attempting to move anything other than
	// your own piece.
	NotMovablePieceError
	// IllegalCastleError is any illegal castle.
	IllegalCastleError
	// IllegalKingTurnError is a move other than a king during a king-turn.
	IllegalKingTurnError
	// IllegalWhirlwindAttackError is a whirlwind attack from anything other
	// than a Warrior King on a king-turn
	IllegalWhirlwindAttackError
	// IllegalPromotionError is a promotion that doesn't involve a pawn moving
	// to the last rank.
	IllegalPromotionError
	// UnreachableSquareError is a move to an unreachable square. Reachability
	// is always checked before capture legality.
	UnreachableSquareError
	// IllegalCaptureError is a move that looks like a capture but isn't valid.
	IllegalCaptureError
	// IllegalRampageError is any Elephant move that captures and doesn't follow
	// the rampage rules.
	IllegalRampageError
	// TooManyDuelsError is a move with more duels than there are pieces
	// captures.
	TooManyDuelsError
	// NotEnoughStonesError is a move where the challenges or responses would
	// result in negative stone counts.
	NotEnoughStonesError
	// NotDuelableError is a move that attempts to duel with a king.
	NotDuelableError
)

func (code IllegalMoveError) Error() string {
	switch code {
	case GameOverError:
		return "game is over"
	case IllegalDropError:
		return "illegal drop"
	case IllegalPassError:
		return "illegal pass"
	case NotMovablePieceError:
		return "not from a movable piece"
	case IllegalCastleError:
		return "illegal castle"
	case IllegalKingTurnError:
		return "piece cannot move during king-turn"
	case IllegalWhirlwindAttackError:
		return "illegal whirlwind attack"
	case IllegalPromotionError:
		return "illegal promotion"
	case UnreachableSquareError:
		return "unreachable square"
	case IllegalCaptureError:
		return "illegal capture"
	case IllegalRampageError:
		return "illegal rampage"
	case TooManyDuelsError:
		return "too many duels"
	case NotEnoughStonesError:
		return "not enough stones"
	case NotDuelableError:
		return "cannot duel with kings"
	default:
		panic("invalid error code")
	}
}

// ParseError represents the error that was encountered when parsing a board or
// move.
type ParseError string

func (msg ParseError) Error() string {
	return string(msg)
}
