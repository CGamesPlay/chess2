package chess2

import (
	"math/rand"
	"testing"
)

// Generates random, sometimes invalid boards
func randomBoard() Board {
	filled := rand.Uint64()
	white := rand.Uint64() & filled
	black := filled &^ white
	pawns := rand.Uint64() & filled &^ maskRank[0] &^ maskRank[7]
	rooks := rand.Uint64() & filled &^ pawns
	bishops := rand.Uint64() & filled &^ pawns &^ rooks
	knights := rand.Uint64() & filled &^ pawns &^ rooks &^ bishops
	kings := rand.Uint64() & filled &^ pawns &^ rooks &^ bishops &^ knights
	queens := filled &^ pawns &^ rooks &^ bishops &^ knights &^ kings
	return Board{
		pieces: [6]uint64{kings, queens, bishops, knights, rooks, pawns},
		colors: [2]uint64{white, black},
	}
}

// Generates a random, sometimes invalid game
func randomGame() Game {
	game := Game{
		flags:  VariantChess2,
		board:  randomBoard(),
		armies: [2]Army{Army(rand.Uint64() & armyMask), Army(rand.Uint64() & armyMask)},
		stones: [2]int{rand.Intn(7), rand.Intn(7)},
	}
	return game
}

func BenchmarkPieceAt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		board := randomBoard()
		for i := uint8(0); i < 64; i++ {
			board.PieceAt(Square{Address: i})
		}
	}
}

func BenchmarkAttackMask(b *testing.B) {
	for n := 0; n < b.N; n++ {
		game := randomGame()
		for i := uint8(0); i < 64; i++ {
			game.attackMask(Square{Address: i})
		}
	}
}
