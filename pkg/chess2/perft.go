package chess2

var promotions = []PieceType{TypeQueen, TypeRook, TypeBishop, TypeKnight}

// BruteforceMoveList calls the given function once for each possible move.
// Drop moves are not emitted, but passes are.
func BruteforceMoveList(send func(Move)) {
	for from := uint8(0); from < 64; from++ {
		for to := uint8(0); to < 64; to++ {
			move := Move{
				From: Square{Address: from},
				To:   Square{Address: to},
			}
			send(move)
			if move.To.Y() == 0 && move.From.Y() == 1 ||
				move.To.Y() == 7 && move.From.Y() == 6 {
				for _, promotion := range promotions {
					move.Piece = NewPiece(promotion, ArmyNone, ColorWhite)
					send(move)
				}
			}
		}
	}
	send(MovePass)
}

// Perft returns the number of valid sequences of moves of length depth from the
// given game. Challenges are never issued while counting moves.
func Perft(game Game, depth int) []uint64 {
	if depth <= 0 {
		return make([]uint64, 0)
	}
	results := make([]uint64, depth)
	doPerft(game, depth, results, func(g Game) []Move {
		return g.GenerateLegalMoves()
	})
	return results
}

// PerftBruteforce is similar to Perft, except that it generates the moves by
// trying every possible square combination rather than using the (much faster)
// move generator. It is useful for testing the move generator.
func PerftBruteforce(game Game, depth int) []uint64 {
	if depth <= 0 {
		return make([]uint64, 0)
	}
	results := make([]uint64, depth)
	doPerft(game, depth, results, func(g Game) []Move {
		moves := make([]Move, 0, 64)
		BruteforceMoveList(func(candidate Move) {
			if err := g.ValidateLegalMove(candidate); err == nil {
				moves = append(moves, candidate)
			}
		})
		return moves
	})
	return results
}

func doPerft(game Game, depth int, results []uint64, getMoves func(Game) []Move) {
	moves := getMoves(game)
	results[len(results)-depth] += uint64(len(moves))
	if depth == 1 {
		return
	}
	for _, move := range moves {
		child := game.ApplyMove(move)
		doPerft(child, depth-1, results, getMoves)
	}
}
