package chess2

// Perft returns the number of valid sequences of moves of length depth from the
// given game. Challenges are never issued while counting moves.
func Perft(game Game, depth int) uint64 {
	panic("not implemented")
}

// PerftBruteforce is similar to Perft, except that it generates the moves by
// trying every possible square combination rather than using the (much faster)
// move generator. It is useful for testing the move generator.
func PerftBruteforce(game Game, depth int) uint64 {
	return doPerft(game, depth, func(g Game) []Move {
		moves := make([]Move, 0, 64)
		for from := uint8(0); from < 64; from++ {
			for to := uint8(0); to < 64; to++ {
				candidate := Move{
					From: Square{Address: from},
					To:   Square{Address: to},
				}
				if err := g.ValidateLegalMove(candidate); err == nil {
					moves = append(moves, candidate)
				}
			}
		}
		return moves
	})
}

func doPerft(game Game, depth int, getMoves func(Game) []Move) uint64 {
	if depth <= 0 {
		return 1
	}
	moves := getMoves(game)
	if depth == 1 {
		return uint64(len(moves))
	}
	var result uint64
	for _, move := range moves {
		child := game.ApplyMove(move)
		result += doPerft(child, depth-1, getMoves)
	}
	return result
}
