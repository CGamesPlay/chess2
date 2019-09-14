package chess2

import "fmt"
import "strings"

const (
	// FenEmpty is a FEN for an empty board
	FenEmpty = "8/8/8/8/8/8/8/8"
	// FenDefault is a FEN for a normal starting position
	FenDefault = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR"
)

var (
	fenToPieceType = map[rune]PieceType{
		'k': TypeKing,
		'q': TypeQueen,
		'b': TypeBishop,
		'n': TypeKnight,
		'r': TypeRook,
		'p': TypePawn,
	}
	pieceTypeToFen = map[PieceType]rune{
		TypeKing:   'k',
		TypeQueen:  'q',
		TypeBishop: 'b',
		TypeKnight: 'n',
		TypeRook:   'r',
		TypePawn:   'p',
	}
)

// ParseFenPiece takes the FEN for a specific piece and converts it to a Piece
// in the Classic army.
func ParseFenPiece(code rune) (Piece, error) {
	// Convert the code to a lowercase letter
	typeCode := code | 0x20
	// Detect if the code is uppercase, which means white
	isWhite := (code & 0x20) == 0
	if resultType, ok := fenToPieceType[typeCode]; ok {
		var color Color
		if isWhite {
			color = ColorWhite
		} else {
			color = ColorBlack
		}
		return NewPiece(resultType, ArmyClassic, color), nil
	}
	return Piece{}, EpdError(fmt.Sprintf("Invalid piece in FEN: '%s'", string(code)))
}

// EncodeFenPiece returns the FEN code for a specific piece
func EncodeFenPiece(p Piece) rune {
	pieceType := pieceTypeToFen[p.Type()]
	// White pieces are uppercase
	if p.Color() == ColorWhite {
		pieceType &= ^0x20
	}
	return pieceType
}

// ParseFen takes a fen string and returns the Board that is represented by it.
func ParseFen(fen string) (Board, error) {
	result := Board{}
	y := 0
	x := 0
	for _, op := range fen {
		switch {
		case op == '/':
			x = 0
			y++
			if y >= 8 {
				return Board{}, EpdError(fmt.Sprintf("FEN has too many ranks"))
			}
		case x >= 8:
			return Board{}, EpdError(fmt.Sprintf("Rank in FEN is too long"))
		case '1' <= op && op <= '8':
			x += int(op - '0')
		default:
			piece, err := ParseFenPiece(op)
			if err == nil {
				result.SetPieceAt(SquareFromCoords(x, y), piece)
			} else {
				return Board{}, err
			}
			x++
		}
	}
	return result, nil
}

// EncodeFen takes a board and returns the FEN that represents it.
func EncodeFen(board Board) string {
	var sb strings.Builder
	for y := 0; y < 8; y++ {
		if y != 0 {
			sb.WriteRune('/')
		}
		lastOccupied := -1
		for x := 0; x < 8; x++ {
			if piece, occupied := board.PieceAt(SquareFromCoords(x, y)); occupied {
				if lastOccupied != x-1 {
					sb.WriteRune(rune(x-lastOccupied-1) + '0')
				}
				sb.WriteRune(EncodeFenPiece(piece))
				lastOccupied = x
			}
		}
		if lastOccupied != 7 {
			sb.WriteRune(rune(7-lastOccupied) + '0')
		}
	}
	return sb.String()
}
