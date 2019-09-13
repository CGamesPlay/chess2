package chess2

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
	castleWhiteKingside  = SquareFromName("H1").mask()
	castleWhiteQueenside = SquareFromName("A1").mask()
	castleBlackKingside  = SquareFromName("H8").mask()
	castleBlackQueenside = SquareFromName("A8").mask()
	castleKingside       = castleWhiteKingside | castleBlackKingside
	castleQueenside      = castleWhiteQueenside | castleBlackQueenside
)

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
