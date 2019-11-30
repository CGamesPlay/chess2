package chess2

import (
	"math/bits"
)

// GameState describes if the game is in progress, and the winner
type GameState int

const (
	// GameInProgress means the game is not yet finished.
	GameInProgress = GameState(iota)
	// GameOverWhite means the game is over and white won.
	GameOverWhite
	// GameOverBlack means the game is over and black won.
	GameOverBlack
	// GameOverDraw means the game is over and was a draw.
	GameOverDraw
)

var (
	castleWhiteKingside  = SquareFromName("H1").Mask()
	castleWhiteQueenside = SquareFromName("A1").Mask()
	castleBlackKingside  = SquareFromName("H8").Mask()
	castleBlackQueenside = SquareFromName("A8").Mask()
	castleKingside       = castleWhiteKingside | castleBlackKingside
	castleQueenside      = castleWhiteQueenside | castleBlackQueenside
	castles              = []uint64{castleWhiteKingside, castleWhiteQueenside, castleBlackKingside, castleBlackQueenside}
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

type moveExecution struct {
	epSquare       Square
	attackerStones int
	defenderStones int
	isCapture      bool
	duels          []Duel
	dryRun         bool
	err            error
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

// GameState returns the current state of the game.
func (g *Game) GameState() GameState {
	return g.gameState
}

func (g *Game) updateGameState() {
	// TODO - decide if the game is over
}

// IsInCheck determines if the given player is currently in check, regardless of
// if they are the player to move. If the game is over due to checkmate, this
// method will return true for the losing player.
func (g *Game) IsInCheck(color Color) bool {
	otherArmy := g.armies[ColorIdx(OtherColor(color))]
	kingMask := g.board.colorMask(color) & g.board.pieceMask(TypeKing)
	enemyMask := g.board.colorMask(OtherColor(color))
	if otherArmy == ArmyAnimals {
		// Calculate check from elephants separately
		//enemyMask &^= g.board.pieceMask(TypeRook)
		// TODO - the color is not in check if a rampage would kill one of the
		// other color's kings as well.
	}
	threatenedMask := g.fullAttackMask(enemyMask)
	return threatenedMask&kingMask != 0
}

// GenerateLegalMoves returns an array of all legal moves from the current board
// state.
func (g *Game) GenerateLegalMoves() []Move {
	var results []Move
	g.generateMoves(func(m Move) {
		if err := g.ValidateLegalMove(m); err == nil {
			results = append(results, m)
		}
	})
	return results
}

// As generateMovesFrom, but for every square on the game board.
func (g *Game) generateMoves(send func(Move)) {
	fromMask := g.board.colorMask(g.toMove)
	if g.kingTurn {
		fromMask &= g.board.pieceMask(TypeKing)
		send(MovePass)
	}
	eachSquareInMask(fromMask, func(from Square) {
		g.generateMovesFrom(from, send)
	})
}

// Generates a superset of pseudo-legal moves originating from the given square.
// The returned set includes all pseudo-legal moves from this game state, but
// may include illegal moves as well. Each move will only be sent once.
func (g *Game) generateMovesFrom(from Square, sendOne func(Move)) {
	piece, found := g.board.PieceAt(from)
	piece = piece.WithArmy(g.armies[ColorIdx(piece.Color())])
	if !found {
		return
	}

	lastRank := MaskRank[7*ColorIdx(piece.Color())]
	// send the move, but enumerate all possible promotions if there are any
	send := func(move Move) {
		if piece.Type() == TypePawn && move.To.Mask()&lastRank != 0 {
			for t := range basicTypeNames {
				move.Piece = NewPiece(t, ArmyNone, ColorWhite)
				sendOne(move)
			}
		} else {
			sendOne(move)
		}
	}

	// Cover all normal sliding attacks
	attackMask := g.attackMask(from)
	eachSquareInMask(attackMask, func(to Square) {
		send(Move{From: from, To: to})
	})

	if piece.Type() == TypePawn {
		if piece.Army() == ArmyNemesis {
			// Nemesis moves
			for y := -1; y <= 1; y++ {
				if from.Y()+y < 0 || from.Y()+y > 7 {
					continue
				}
				for x := -1; x <= 1; x++ {
					if from.X()+x < 0 || from.X()+x > 7 {
						continue
					}
					to := SquareFromCoords(from.X()+x, from.Y()+y)
					if to.Mask()&attackMask == 0 {
						send(Move{From: from, To: to})
					}
				}
			}
		} else {
			// Pawn advancement
			sign := ColorIdx(piece.Color())*2 - 1
			for n := 1; n <= 2; n++ {
				y := from.Y() + sign*n
				if y >= 0 && y < 8 {
					send(Move{From: from, To: SquareFromCoords(from.X(), y)})
				}
			}
		}
	} else if piece.Name() == PieceNameClassicKing {
		// Castling
		if from.X() == 4 {
			send(Move{From: from, To: SquareFromCoords(from.X()-2, from.Y())})
			send(Move{From: from, To: SquareFromCoords(from.X()+2, from.Y())})
		}
	} else if piece.Name() == PieceNameTwoKingsKing {
		send(Move{From: from, To: from})
	}
}

// ApplyMove clones the receiver, applies the given move to the clone, and
// returns the clone. It does not validate that the move is legal, and if an
// illegal move is made the resulting Game may not be in a valid state.
func (g *Game) ApplyMove(move Move) Game {
	return applyMoveTo(g, move)
}

func applyMoveTo(old *Game, move Move) Game {
	g := *old
	if move.IsDrop() {
		p := move.Piece.WithArmy(g.armies[ColorIdx(move.Piece.Color())])
		g.board.SetPieceAt(move.To, p)
		return g
	}

	// Advance turn
	movingPlayer := g.toMove
	epSquare := g.epSquare
	if !g.kingTurn {
		g.halfmoveClock++
		g.epSquare = InvalidSquare
	}
	if g.armies[ColorIdx(movingPlayer)] == ArmyTwoKings && !g.kingTurn {
		g.kingTurn = true
	} else {
		g.kingTurn = false
		g.toMove = OtherColor(movingPlayer)
	}
	if g.toMove == ColorWhite && !g.kingTurn {
		g.fullmoveNumber++
	}

	if move.IsPass() {
		return g
	}

	// Handle captures and duels
	p, _ := g.board.PieceAt(move.From)
	p = p.WithArmy(g.armies[ColorIdx(p.Color())])
	me := moveExecution{
		epSquare:       epSquare,
		attackerStones: g.stones[ColorIdx(movingPlayer)],
		defenderStones: g.stones[1-ColorIdx(movingPlayer)],
		duels:          move.Duels[:],
		dryRun:         false,
	}
	survived := g.handleAllCaptures(p, move, &me)

	// Update stones
	g.stones[ColorIdx(movingPlayer)] = me.attackerStones
	g.stones[1-ColorIdx(movingPlayer)] = me.defenderStones
	isZeroingMove := me.isCapture

	// Move the piece
	delta := int(move.To.Address) - int(move.From.Address)
	if survived {
		if p.Name() != PieceNameAnimalsBishop || !isZeroingMove {
			g.board.ClearPieceAt(move.From)
			if move.Piece != InvalidPiece {
				promoted := NewPiece(move.Piece.Type(), p.Army(), p.Color())
				g.board.SetPieceAt(move.To, promoted)
			} else {
				g.board.SetPieceAt(move.To, p)
			}
		}
	} else {
		g.board.ClearPieceAt(move.From)
	}
	if p.Type() == TypePawn && p.Army() != ArmyNemesis {
		isZeroingMove = true
		if SquareDistance(move.From, move.To) > 1 {
			g.epSquare = Square{Address: uint8(int(move.From.Address) + delta/2)}
		}
	} else if p.Name() == PieceNameClassicKing {
		// Clear castling rights on king move
		firstRank := MaskRank[7-7*ColorIdx(p.Color())]
		g.castlingRights &^= firstRank

		// Move rook when castling
		if delta == -2 {
			g.board.MovePiece(
				Square{Address: move.To.Address - 2},
				Square{Address: move.To.Address + 1},
			)
		} else if delta == 2 {
			g.board.MovePiece(
				Square{Address: move.To.Address + 1},
				Square{Address: move.To.Address - 1},
			)
		}
	}

	// Update the halfmove clock
	if isZeroingMove {
		g.halfmoveClock = 0
	}
	for _, mask := range castles {
		if move.From.Mask()&mask != 0 {
			g.castlingRights &^= mask
		}
	}

	return g
}

func (g *Game) handleAllCaptures(p Piece, move Move, me *moveExecution) bool {
	survived := true
	diff := int(move.To.Address) - int(move.From.Address)
	if p.Name() == PieceNameAnimalsRook {
		var delta, step int
		if diff <= -8 {
			step = -8
		} else if diff < 0 {
			step = -1
		} else if diff < 8 {
			step = 1
		} else {
			step = 8
		}
		for {
			delta += step
			if delta > 64 || delta < -64 {
				panic("invalid rampage")
			}
			target := Square{Address: uint8(int(move.From.Address) + delta)}
			survived = g.handleCapture(p, target, me)
			if !survived {
				break
			}
			if target == move.To {
				break
			}
		}
	} else if p.Name() == PieceNameTwoKingsKing && move.From == move.To {
		captureMask := dist1Mask[move.From.Address] & g.board.occupiedMask()
		if g.armies[1-ColorIdx(p.Color())] == ArmyReaper {
			captureMask &^= g.board.pieceMask(TypeRook) & g.board.colorMask(OtherColor(p.Color()))
		}
		eachSquareInMask(captureMask, func(target Square) {
			g.handleCapture(p, target, me)
		})
	} else {
		target := move.To
		if p.Type() == TypePawn && move.To == me.epSquare {
			// white = -1, black = 1
			sign := -1 + ColorIdx(p.Color())*2
			if diff == sign*7 || diff == sign*9 {
				target = Square{Address: uint8(int(move.To.Address) - sign*8)}
			}
		}
		survived = g.handleCapture(p, target, me)
	}
	return survived
}

func (g *Game) handleCapture(attacker Piece, target Square, me *moveExecution) bool {
	defender, isCapture := g.board.PieceAt(target)
	if !isCapture {
		return true
	}
	me.isCapture = me.isCapture || isCapture
	survived := true
	if !me.dryRun {
		g.board.ClearPieceAt(target)
	}
	if len(me.duels) > 0 {
		if d := me.duels[0]; d.IsStarted() {
			if (attacker.Type() == TypeKing || defender.Type() == TypeKing) && me.err == nil {
				me.err = NotDuelableError
			}
			if DuelingRank(attacker.Type()) < DuelingRank(defender.Type()) {
				if me.attackerStones > 0 {
					me.attackerStones--
				} else if me.err == nil {
					me.err = NotEnoughStonesError
				}
			}
			if d.Challenge() > me.defenderStones ||
				d.Response() > me.attackerStones && me.err == nil {
				me.err = NotEnoughStonesError
			}
			me.defenderStones -= d.Challenge()
			me.attackerStones -= d.Response()
			if d.Challenge() == 0 && d.Response() == 0 {
				if d.Gain() {
					if me.attackerStones < 6 {
						me.attackerStones++
					}
				} else {
					if me.defenderStones > 0 {
						me.defenderStones--
					}
				}
			}
			survived = d.Challenge() <= d.Response()
		}
		me.duels = me.duels[1:]
	}
	if defender.Color() != attacker.Color() && defender.Type() == TypePawn && me.attackerStones < 6 {
		me.attackerStones++
	}
	return survived
}

// ValidatePseudoLegalMove returns an error describing why the given move is not
// pseudo-legal.
//
// A move is "pseudo-legal" if it is one of the player to move's pieces; follows
// the distance and direction rules for the piece moved; all intermediate
// squares are empty (except for jump moves) or capturable (for the Elephant's
// rampage); and the target square is empty or contains a capturable piece.
// Additionally, a pass move is pseudo-legal during a king turn.
func (g *Game) ValidatePseudoLegalMove(move Move) error {
	// Basic checks
	if g.gameState != GameInProgress {
		return GameOverError
	} else if move.IsDrop() {
		return IllegalDropError
	} else if move.IsPass() {
		if !g.kingTurn {
			return IllegalPassError
		}
		return validateNoDuels(move, TooManyDuelsError)
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
			move.Piece.Type() == TypeKing ||
			move.Piece.Type() == TypeQueen && piece.Army() == ArmyTwoKings {
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
			return validateNoDuels(move, TooManyDuelsError)
		} else if forwardDiff == 16 && piece.Army() != ArmyNemesis {
			// targetRank = 4 for white, 3 for black
			targetRank := 4 - ColorIdx(piece.Color())
			middleOccupied := betweenMask[move.From.Address][move.To.Address]&g.board.occupiedMask() != 0
			if middleOccupied || move.To.Y() != targetRank {
				return UnreachableSquareError
			} else if targetOccupied {
				return IllegalCaptureError
			}
			return validateNoDuels(move, TooManyDuelsError)
		} else if forwardDiff == 7 || forwardDiff == 9 {
			if SquareDistance(move.From, move.To) > 1 {
				return UnreachableSquareError
			} else if move.To == g.epSquare {
				return g.ValidateDuels(move)
			}
		}
		if piece.Army() == ArmyNemesis && !targetOccupied {
			// Check for nemesis move
			mask := singleStepMask(move.From, g.board.pieceMask(TypeKing)&g.board.colorMask(OtherColor(piece.Color())))
			if move.To.Mask()&mask == 0 {
				return UnreachableSquareError
			} else if targetOccupied {
				return IllegalCaptureError
			}
			return validateNoDuels(move, TooManyDuelsError)
		}
		if !targetOccupied {
			return UnreachableSquareError
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
			checkMask := Square{Address: move.From.Address + uint8(diff/2)}.Mask()
			checkMask |= move.From.Mask()
			checkMask |= move.To.Mask()
			threats := g.board.colorMask(OtherColor(piece.Color()))
			threatenedMask := g.fullAttackMask(threats)
			if threatenedMask&checkMask != 0 {
				return IllegalCastleError
			}
			return validateNoDuels(move, NotDuelableError)
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
		return validateNoDuels(move, NotDuelableError)
	}

	// Check valid sliding attack
	if g.attackMask(move.From)&move.To.Mask() == 0 {
		return UnreachableSquareError
	}

	// Check captures
	// noncapturableMask is the mask of pieces that cannot be captured by the
	// moving piece. In the case of an elephant, noncapturableMask is the
	// mask of pieces that can stop a rampage.
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
		if army == ArmyReaper {
			// Cannot capture reaper rook
			noncapturableMask |= colorPieces & g.board.pieceMask(TypeRook)
		}
		if army == ArmyAnimals {
			// Cannot capture elephants more than 3 spaces away
			noncapturableMask |= colorPieces & g.board.pieceMask(TypeRook) & ^dist2Mask[move.From.Address]
		}
	}
	// visitedSquares is a mask of squares visited by the move. Generally these
	// need to be empty, except for the last one, for the move to be valid.
	visitedSquares := move.To.Mask()
	if piece.Name() != PieceNameReaperRook && piece.Name() != PieceNameReaperQueen {
		// These pieces attack each square they pass through
		visitedSquares |= betweenMask[move.From.Address][move.To.Address]
	}
	if piece.Name() != PieceNameAnimalsRook && visitedSquares&noncapturableMask != 0 {
		return IllegalCaptureError
	}

	if piece.Name() == PieceNameAnimalsRook {
		diff := int(move.To.Address) - int(move.From.Address)
		if SquareDistance(move.From, move.To) < 3 && visitedSquares&g.board.occupiedMask() != 0 {
			// Moving less than 3 spaces is only allowed if the elephant hits a
			// noncapturable piece or the edge of the board, or if all of the
			// spaces are empty.
			var step int
			var hitWall bool
			if diff <= -8 {
				step = -8
				hitWall = move.To.Y() == 0
			} else if diff < 0 {
				step = -1
				hitWall = move.To.X() == 0
			} else if diff < 8 {
				step = 1
				hitWall = move.To.X() == 7
			} else {
				step = 8
				hitWall = move.To.Y() == 7
			}
			if !hitWall {
				checkAddr := int(move.To.Address) + step
				if checkAddr >= 0 && checkAddr < 64 {
					wall := Square{Address: uint8(checkAddr)}
					if noncapturableMask&wall.Mask() == 0 {
						return IllegalRampageError
					}
				}
			}
		}
		if visitedSquares&g.board.colorMask(piece.Color())&g.board.pieceMask(TypeKing) != 0 {
			// Cannot capture own king
			return IllegalRampageError
		}
	}

	return g.ValidateDuels(move)
}

func validateNoDuels(move Move, err error) error {
	for _, d := range move.Duels {
		if d.IsStarted() {
			return err
		}
	}
	return nil
}

// ValidateLegalMove returns an error describing why the given move is not
// legal.
//
// - A move is "into check" if it leaves the board in a state where any of the
//   player's kings are threatened.
// - A move is "legal" if it is pseudo-legal and not into check.
func (g *Game) ValidateLegalMove(move Move) error {
	if err := g.ValidatePseudoLegalMove(move); err != nil {
		return err
	}
	result := g.ApplyMove(move)
	if result.IsInCheck(g.toMove) {
		return MoveIntoCheckError
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
		return orthAttackMask[from.Address][0] & dist3Mask[from.Address]
	default:
		panic("Invalid piece type")
	}
}

func (g *Game) fullAttackMask(from uint64) (result uint64) {
	eachSquareInMask(from, func(sq Square) {
		result |= g.attackMask(sq)
	})
	return
}

// ValidateDuels returns an error describing why the duels on given move are not
// legal.
func (g *Game) ValidateDuels(move Move) error {
	if move.IsDrop() || move.IsPass() {
		return validateNoDuels(move, TooManyDuelsError)
	}

	p, _ := g.board.PieceAt(move.From)
	p = p.WithArmy(g.armies[ColorIdx(p.Color())])
	me := moveExecution{
		epSquare:       g.epSquare,
		attackerStones: g.stones[ColorIdx(g.toMove)],
		defenderStones: g.stones[1-ColorIdx(g.toMove)],
		duels:          move.Duels[:],
		dryRun:         true,
	}
	g.handleAllCaptures(p, move, &me)
	if me.err != nil {
		return me.err
	}
	for _, d := range me.duels {
		if d.IsStarted() {
			return TooManyDuelsError
		}
	}
	return nil
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
	if color == ColorWhite {
		if isQueenside {
			return castleWhiteQueenside
		}
		return castleWhiteKingside
	}
	if isQueenside {
		return castleBlackQueenside
	}
	return castleBlackKingside
}
