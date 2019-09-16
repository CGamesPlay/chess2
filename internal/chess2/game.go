package chess2

import (
	"math/bits"
)

// GameState describes if the game is in progress, and the winner
type GameState int

const (
	// GameInProgress means the game is not yet finished.
	GameInProgress = GameState(0)
	// GameOverWhite means the game is over and white won.
	GameOverWhite = GameState(1)
	// GameOverBlack means the game is over and black won.
	GameOverBlack = GameState(2)
	// GameOverDraw means the game is over and was a draw.
	GameOverDraw = GameState(3)
)

var (
	castleWhiteKingside  = SquareFromName("H1").Mask()
	castleWhiteQueenside = SquareFromName("A1").Mask()
	castleBlackKingside  = SquareFromName("H8").Mask()
	castleBlackQueenside = SquareFromName("A8").Mask()
	castleKingside       = castleWhiteKingside | castleBlackKingside
	castleQueenside      = castleWhiteQueenside | castleBlackQueenside
)

// buildSlidingAttackMask creates a bitmask for each of 64 squares of the
// squares reachable by adding each delta once.
func buildSlidingAttackMask(sqStart Square, mask uint64, deltas []int) uint64 {
	var results uint64
	for _, delta := range deltas {
		sq := sqStart
		for {
			prevSquare := sq
			if int(sq.Address)+delta < 0 || int(sq.Address)+delta >= 64 {
				break
			}
			sq.Address = uint8(int(sq.Address) + delta)
			if SquareDistance(prevSquare, sq) > 3 {
				// Dont wrap around edges
				break
			}
			results |= sq.Mask()
			if mask&sq.Mask() != 0 {
				break
			}
		}
	}
	return results
}

// buildSlidingAttackMaskAll calls buildSlidingAttack mask for each square on
// the board.
func buildSlidingAttackMaskAll(deltas []int) [64]uint64 {
	var results [64]uint64
	for startAddr := range results {
		sqStart := Square{Address: uint8(startAddr)}
		results[startAddr] = buildSlidingAttackMask(sqStart, MaskFull, deltas)
	}
	return results
}

// builkdAttackTable computes a mask used to calculate the set of threatened
// squares given a source square. The maskTable is the set of all squares
// threatened from the source square on an empty board, and attackTable maps the
// source square and a subset of the corresponding maskTable row to calculate
// the threatened squares on an occupied board.
func buildAttackTable(deltas []int) (maskTable [64]uint64, attackTable [64]map[uint64]uint64) {
	for startAddr := range maskTable {
		sqStart := Square{Address: uint8(startAddr)}
		mask := buildSlidingAttackMask(sqStart, MaskEmpty, deltas)
		maskTable[startAddr] = mask
		attackTable[startAddr] = make(map[uint64]uint64)
		eachBitSubset64(mask, func(subset uint64) {
			attackTable[startAddr][subset] = buildSlidingAttackMask(sqStart, subset, deltas)
		})
	}
	return
}

var (
	// Mask of all squares within distance 1.
	dist1Mask = buildSlidingAttackMaskAll([]int{-9, -8, -7, -1, 1, 7, 8, 9})
	// Mask of all squares within distance 2.
	dist2Mask = buildSlidingAttackMaskAll([]int{
		-18, -17, -16, -15, -14,
		-10, -9, -8, -7, -6,
		-2, -1, 1, 2,
		6, 7, 8, 9, 10,
		14, 15, 16, 17, 18,
	})
	// Mask of all squares within distance 3.
	dist3Mask = buildSlidingAttackMaskAll([]int{
		-27, -26, -25, -24, -23, -22, -21,
		-19, -18, -17, -16, -15, -14, -13,
		-11, -10, -9, -8, -7, -6, -5,
		-3, -2, -1, 1, 2, 3,
		5, 6, 7, 8, 9, 10, 11,
		13, 14, 15, 16, 17, 18, 19,
		21, 22, 23, 24, 25, 26, 27,
	})
	// Mask of all squares that are orthagonally adjacent.
	adjacentMask = buildSlidingAttackMaskAll([]int{-8, -1, 1, 8})
	// Mask of squares reachable with a knight move
	knightMask = buildSlidingAttackMaskAll([]int{17, 15, 10, 6, -17, -15, -10, -6})
	// Masks describing squares reachable diagonally/orthagonally from the given
	// square, and attackable given a particular occupancy mask.
	diagMask, diagAttackMask = buildAttackTable([]int{-9, -7, 7, 9})
	orthMask, orthAttackMask = buildAttackTable([]int{-8, -1, 1, 8})

	// Masks describing the files that have to be empty for a castle.
	queensideMask = uint64(0x0e0e0e0e0e0e0e0e)
	kingsideMask  = uint64(0x6060606060606060)
)

func buildBetweenMask() (results [64][64]uint64) {
	for a := 0; a < 64; a++ {
		maskA := uint64(1 << a)
		for b := 0; b < 64; b++ {
			maskB := uint64(1 << b)
			if diagMask[a]&maskB != 0 {
				results[a][b] = diagAttackMask[a][maskB] & diagAttackMask[b][maskA]
			} else if orthMask[a]&maskB != 0 {
				results[a][b] = orthAttackMask[a][maskB] & orthAttackMask[b][maskA]
			}
		}
	}
	return
}

// The set of squares that are between the two addresses, exclusively.
var betweenMask = buildBetweenMask()

func singleStepMask(from Square, toMask uint64) uint64 {
	mask := MaskEmpty
	eachSquareInMask(toMask, func(to Square) {
		dx, dy := to.X()-from.X(), to.Y()-from.Y()
		if dx < 0 {
			mask |= 1 << (from.Address - 1)
			if dy < 0 {
				mask |= 1 << (from.Address - 9)
			} else if dy > 0 {
				mask |= 1 << (from.Address + 7)
			}
		} else if dx > 0 {
			mask |= 1 << (from.Address + 1)
			if dy < 0 {
				mask |= 1 << (from.Address - 7)
			} else if dy > 0 {
				mask |= 1 << (from.Address + 9)
			}
		}
		if dy < 0 {
			mask |= 1 << (from.Address - 8)
		} else if dy > 0 {
			mask |= 1 << (from.Address + 8)
		}
	})
	return mask
}

// A Game fully describes a Chess 2 game.
type Game struct {
	board          Board
	castlingRights uint64
	armies         [2]Army
	stones         [2]int
	toMove         Color
	kingTurn       bool
	gameState      GameState
	halfmoveClock  int
	fullmoveNumber int
	epSquare       Square
}

// GameFromArmies initializes a new Game with the provided armies.
func GameFromArmies(white, black Army) Game {
	board, err := ParseFen(FenDefault)
	if err != nil {
		panic(err)
	}
	if white == ArmyTwoKings {
		board.ReplacePieces(ColorWhite, TypeQueen, TypeKing)
	}
	if black == ArmyTwoKings {
		board.ReplacePieces(ColorBlack, TypeQueen, TypeKing)
	}
	return Game{
		board:          board,
		castlingRights: castleKingside | castleQueenside,
		armies:         [2]Army{white, black},
		stones:         [2]int{3, 3},
		epSquare:       InvalidSquare,
	}
}

func (g *Game) updateGameState() {
	// TODO - decide if the game is over
}

// GeneratePseudoLegalMoves returns an array of all pseudo-legal moves from the
// current board state.
func (g *Game) GeneratePseudoLegalMoves() []Move {
	var results []Move
	g.generatePseudoLegalMoves(func(m Move) {
		results = append(results, m)
	})
	return results
}

func (g *Game) generatePseudoLegalMoves(send func(Move)) {

	fromMask := g.board.colorMask(g.toMove)
	eachSquareInMask(fromMask, func(from Square) {
		mask := g.attackMask(from)
		eachSquareInMask(mask, func(to Square) {
			send(Move{From: from, To: to})
		})
	})
}

// ValidatePseudoLegalMove returns an error describing why the given move is not
// pseudo-legal.
//
// A move is "pseudo-legal" if it follows the distance and direction rules for
// the piece moved; all intermediate squares are empty (except for jump moves)
// or capturable (for the Elephant's rampage); and the target square is empty or
// contains a capturable piece. Additionally:
//  - A move that captures one of the player's own kings is not pseudo-legal.
//  - A pass move is pseudo-legal during a king turn.
func (g *Game) ValidatePseudoLegalMove(move Move) error {
	// Basic checks
	if g.gameState != GameInProgress {
		return GameOverError
	} else if move.IsDrop() {
		return IllegalDropError
	} else if move.IsPass() && !g.kingTurn {
		return IllegalPassError
	}

	// Check turn
	piece, found := g.board.PieceAt(move.From)
	piece = piece.WithArmy(g.armies[ColorIdx(piece.Color())])
	if !found || piece.Color() != g.toMove {
		return NotMovablePieceError
	} else if g.kingTurn && piece.Type() != TypeKing {
		return IllegalKingTurnError
	}

	// Check promotions
	lastRank := MaskRank[7*ColorIdx(piece.Color())]
	if move.Piece != InvalidPiece {
		if piece.Type() != TypePawn ||
			move.To.Mask()&lastRank == 0 ||
			move.Piece.Type() == TypePawn ||
			move.Piece.Type() == TypeKing {
			return IllegalPromotionError
		}
	} else if piece.Type() == TypePawn && move.To.Mask()&lastRank != 0 {
		// Pawns must promote at last rank
		return IllegalPromotionError
	}

	// Check noncapturing pawn moves
	if piece.Type() == TypePawn {
		diff := int(move.To.Address) - int(move.From.Address)
		sign := ColorIdx(piece.Color())*2 - 1
		forwardDiff := sign * diff
		_, targetOccupied := g.board.PieceAt(move.To)
		if forwardDiff < 0 && piece.Army() != ArmyNemesis {
			return UnreachableSquareError
		} else if forwardDiff == 8 {
			if targetOccupied {
				return IllegalCaptureError
			}
			return nil
		} else if forwardDiff == 16 && piece.Army() != ArmyNemesis {
			// targetRank = 4 for white, 3 for black
			targetRank := 4 - ColorIdx(piece.Color())
			middleOccupied := betweenMask[move.From.Address][move.To.Address]&g.board.occupiedMask() != 0
			if middleOccupied || move.To.Y() != targetRank {
				return UnreachableSquareError
			} else if targetOccupied {
				return IllegalCaptureError
			}
			return nil
		} else if forwardDiff == 7 || forwardDiff == 9 {
			if SquareDistance(move.From, move.To) > 1 {
				return UnreachableSquareError
			} else if move.To == g.epSquare {
				return nil
			}
		} else if piece.Army() == ArmyNemesis {
			// Check for nemesis move
			mask := singleStepMask(move.From, g.board.pieceMask(TypeKing)&g.board.colorMask(OtherColor(piece.Color())))
			if move.To.Mask()&mask == 0 {
				return UnreachableSquareError
			} else if targetOccupied {
				return IllegalCaptureError
			}
			return nil
		}
		// Otherwise normal sliding attack rules apply
	}

	// Check castling
	if piece.Type() == TypeKing {
		diff := int(move.To.Address) - int(move.From.Address)
		if diff == 2 || diff == -2 {
			if piece.Name() != PieceNameClassicKing ||
				g.castlingRights&requiredCastlingRight(piece.Color(), diff < 0) == 0 {
				return IllegalCastleError
			}
			firstRank := MaskRank[7-7*ColorIdx(piece.Color())]
			if diff < 0 {
				firstRank &= queensideMask
			} else {
				firstRank &= kingsideMask
			}
			if g.board.occupiedMask()&firstRank != 0 {
				return IllegalCastleError
			}
			return nil
		}
		// Otherwise normal sliding attack rules apply
	}

	// Check whirlwind attack
	if move.From == move.To {
		if piece.Name() != PieceNameTwoKingsKing || !g.kingTurn {
			return IllegalWhirlwindAttackError
		}
		attackedKings := g.board.pieceMask(TypeKing) & g.board.colorMask(piece.Color()) & dist1Mask[move.From.Address]
		if attackedKings != 0 {
			return IllegalCaptureError
		}
		return nil
	}

	// Check valid sliding attack
	if g.attackMask(move.From)&move.To.Mask() == 0 {
		return UnreachableSquareError
	}

	// Check captures
	noncapturableMask := MaskEmpty
	if piece.Name() == PieceNameAnimalsKnight {
		// Cannot capture own king
		noncapturableMask |= g.board.colorMask(piece.Color()) & g.board.pieceMask(TypeKing)
	} else if piece.Name() != PieceNameAnimalsRook {
		// Cannot capture own pieces
		noncapturableMask |= g.board.colorMask(piece.Color())
	}
	for colorIdx := 0; colorIdx < 2; colorIdx++ {
		army := g.armies[colorIdx]
		colorPieces := g.board.colors[colorIdx]
		if piece.Type() != TypeKing && army == ArmyNemesis {
			// Cannot capture nemesis queen
			noncapturableMask |= colorPieces & g.board.pieceMask(TypeQueen)
		}
		if army == ArmyNemesis {
			// Cannot capture nemesis rook
			noncapturableMask |= colorPieces & g.board.pieceMask(TypeRook)
		}
		if army == ArmyAnimals {
			// Cannot capture elephants more than 3 spaces away
			noncapturableMask |= colorPieces & g.board.pieceMask(TypeRook) & ^dist2Mask[move.From.Address]
		}
	}
	attemptedCapturesMask := betweenMask[move.From.Address][move.To.Address] | move.To.Mask()
	if attemptedCapturesMask&noncapturableMask != 0 {
		return IllegalCaptureError
	}

	return nil
}

// attackMask returns the mask of threatened squares from the given square.
// A square which is reachable but not threatened is not included in this mask
// (e.g. pawns advancing).
func (g *Game) attackMask(from Square) uint64 {
	piece, found := g.board.PieceAt(from)
	if !found {
		// No piece, nothing threatened
		return 0
	}
	piece = piece.WithArmy(g.armies[ColorIdx(piece.Color())])
	switch piece.Name() {
	case PieceNameClassicKing, PieceNameBasicKing, PieceNameEmpoweredQueen, PieceNameTwoKingsKing:
		return dist1Mask[from.Address]
	case PieceNameBasicQueen:
		diag := diagMask[from.Address] & g.board.occupiedMask()
		orth := orthMask[from.Address] & g.board.occupiedMask()
		return diagAttackMask[from.Address][diag] | orthAttackMask[from.Address][orth]
	case PieceNameBasicBishop:
		diag := diagMask[from.Address] & g.board.occupiedMask()
		return diagAttackMask[from.Address][diag]
	case PieceNameBasicKnight, PieceNameAnimalsKnight:
		return knightMask[from.Address]
	case PieceNameBasicRook:
		orth := orthMask[from.Address] & g.board.occupiedMask()
		return orthAttackMask[from.Address][orth]
	case PieceNameBasicPawn, PieceNameNemesisPawn:
		sign := ColorIdx(piece.Color())*2 - 1
		mask := uint64(0)
		if from.X() > 0 {
			mask |= 1 << (int(from.Address) - 1 + 8*sign)
		}
		if from.X() < 7 {
			mask |= 1 << (int(from.Address) + 1 + 8*sign)
		}
		return mask
	case PieceNameNemesisQueen:
		diag := diagMask[from.Address] & g.board.occupiedMask()
		orth := orthMask[from.Address] & g.board.occupiedMask()
		// Nemesis queen can only threaten kings
		threat := ^g.board.occupiedMask() | g.board.pieceMask(TypeKing)
		return (diagAttackMask[from.Address][diag] | orthAttackMask[from.Address][orth]) & threat
	case PieceNameEmpoweredBishop, PieceNameEmpoweredKnight, PieceNameEmpoweredRook:
		// Same-colored adjacent pieces, including the from square.
		ownAdjacent := (adjacentMask[from.Address] | from.Mask()) & g.board.colorMask(piece.Color())
		mask := uint64(0)
		if ownAdjacent&g.board.pieceMask(TypeRook) != 0 {
			orth := orthMask[from.Address] & g.board.occupiedMask()
			mask |= orthAttackMask[from.Address][orth]
		}
		if ownAdjacent&g.board.pieceMask(TypeBishop) != 0 {
			diag := diagMask[from.Address] & g.board.occupiedMask()
			mask |= diagAttackMask[from.Address][diag]
		}
		if ownAdjacent&g.board.pieceMask(TypeKnight) != 0 {
			mask |= knightMask[from.Address]
		}
		return mask
	case PieceNameReaperQueen:
		candidates := ^MaskRank[7*ColorIdx(piece.Color())]
		kings := g.board.pieceMask(TypeKing) & g.board.colorMask(OtherColor(piece.Color()))
		return candidates &^ kings
	case PieceNameReaperRook:
		return ^g.board.occupiedMask()
	case PieceNameAnimalsQueen:
		orth := orthMask[from.Address] & g.board.occupiedMask()
		return orthAttackMask[from.Address][orth] | knightMask[from.Address]
	case PieceNameAnimalsBishop:
		diag := diagMask[from.Address] & g.board.occupiedMask()
		return diagAttackMask[from.Address][diag] & dist2Mask[from.Address]
	case PieceNameAnimalsRook:
		orth := orthMask[from.Address] & g.board.occupiedMask()
		return orthAttackMask[from.Address][orth] & dist3Mask[from.Address]
	default:
		panic("Invalid piece type")
	}
}

// Given a mask, call the function for each Square set in the mask.
func eachSquareInMask(mask uint64, f func(Square)) {
	for mask != 0 {
		lsb := uint8(bits.TrailingZeros64(mask))
		f(Square{Address: lsb})
		mask &^= 1 << lsb
	}
}

func requiredCastlingRight(color Color, isQueenside bool) uint64 {
	switch color {
	case ColorWhite:
		if isQueenside {
			return castleWhiteQueenside
		}
		return castleWhiteKingside
	case ColorBlack:
		if isQueenside {
			return castleBlackQueenside
		}
		return castleBlackKingside
	}
	return 0
}
