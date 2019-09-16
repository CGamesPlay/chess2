package chess2

import (
	"fmt"
)

type (
	// PieceType is King, Queen, etc.
	PieceType int
	// Army is Classic, Nemesis, etc
	Army int
	// Color is White or Black
	Color int
	// PieceName is a combination of PieceType and Army
	PieceName int
)

const (
	// TypeNone is not a real piece
	TypeNone = PieceType(0x00)
	// TypeKing means a king
	TypeKing = PieceType(0x01)
	// TypeQueen means a queen
	TypeQueen = PieceType(0x02)
	// TypeBishop means a bishop
	TypeBishop = PieceType(0x03)
	// TypeKnight means a knight
	TypeKnight = PieceType(0x04)
	// TypeRook means a rook
	TypeRook = PieceType(0x05)
	// TypePawn means a pawn
	TypePawn = PieceType(0x06)
	// ArmyNone means the army is not known from the piece alone
	ArmyNone = Army(0x00)
	// ArmyClassic means classic army
	ArmyClassic = Army(0x10)
	// ArmyNemesis means nemesis army
	ArmyNemesis = Army(0x20)
	// ArmyEmpowered means empowered army
	ArmyEmpowered = Army(0x30)
	// ArmyReaper means reaper army
	ArmyReaper = Army(0x40)
	// ArmyTwoKings means two kings army
	ArmyTwoKings = Army(0x50)
	// ArmyAnimals means animals army
	ArmyAnimals = Army(0x60)
	// ColorWhite means white
	ColorWhite = Color(0x00)
	// ColorBlack means black
	ColorBlack = Color(0x80)
	// PieceNameNone is not a real piece
	PieceNameNone = PieceName(0x00)
	// PieceNameBasicKing means a king
	PieceNameBasicKing = PieceName(0x01)
	// PieceNameBasicQueen means a queen
	PieceNameBasicQueen = PieceName(0x02)
	// PieceNameBasicBishop means a bishop
	PieceNameBasicBishop = PieceName(0x03)
	// PieceNameBasicKnight means a knight
	PieceNameBasicKnight = PieceName(0x04)
	// PieceNameBasicRook means a rook
	PieceNameBasicRook = PieceName(0x05)
	// PieceNameBasicPawn means a pawn
	PieceNameBasicPawn = PieceName(0x06)
	// PieceNameClassicKing means Classic King
	PieceNameClassicKing = PieceName(int(ArmyClassic) | int(TypeKing))
	// PieceNameNemesisQueen means Nemesis
	PieceNameNemesisQueen = PieceName(int(ArmyNemesis) | int(TypeQueen))
	// PieceNameNemesisPawn means Nemesis Pawn
	PieceNameNemesisPawn = PieceName(int(ArmyNemesis) | int(TypePawn))
	// PieceNameEmpoweredQueen means Abdicated Queen
	PieceNameEmpoweredQueen = PieceName(int(ArmyEmpowered) | int(TypeQueen))
	// PieceNameEmpoweredBishop means Empowered Bishop
	PieceNameEmpoweredBishop = PieceName(int(ArmyEmpowered) | int(TypeBishop))
	// PieceNameEmpoweredKnight means Empowered Knight
	PieceNameEmpoweredKnight = PieceName(int(ArmyEmpowered) | int(TypeKnight))
	// PieceNameEmpoweredRook means Empowered Rook
	PieceNameEmpoweredRook = PieceName(int(ArmyEmpowered) | int(TypeRook))
	// PieceNameReaperQueen means Reaper
	PieceNameReaperQueen = PieceName(int(ArmyReaper) | int(TypeQueen))
	// PieceNameReaperRook means Ghost
	PieceNameReaperRook = PieceName(int(ArmyReaper) | int(TypeRook))
	// PieceNameTwoKingsKing means Warrior King
	PieceNameTwoKingsKing = PieceName(int(ArmyTwoKings) | int(TypeKing))
	// PieceNameAnimalsQueen means Jungle Queen
	PieceNameAnimalsQueen = PieceName(int(ArmyAnimals) | int(TypeQueen))
	// PieceNameAnimalsBishop means Tiger
	PieceNameAnimalsBishop = PieceName(int(ArmyAnimals) | int(TypeBishop))
	// PieceNameAnimalsKnight means Wild Horse
	PieceNameAnimalsKnight = PieceName(int(ArmyAnimals) | int(TypeKnight))
	// PieceNameAnimalsRook means Elephant
	PieceNameAnimalsRook = PieceName(int(ArmyAnimals) | int(TypeRook))
)

const (
	typeMask  = 0x0f
	armyMask  = 0x70
	colorMask = 0x80
)

// Piece represents a single piece, including army, color, and type.
type Piece struct {
	repr uint8
}

var (
	basicTypeNames = map[PieceType]string{
		TypeNone:   "nothing",
		TypeKing:   "king",
		TypeQueen:  "queen",
		TypeBishop: "bishop",
		TypeKnight: "knight",
		TypeRook:   "rook",
		TypePawn:   "pawn",
	}
	armyNames = map[Army]string{
		ArmyNone:      "basic",
		ArmyClassic:   "Classic",
		ArmyNemesis:   "Nemesis",
		ArmyEmpowered: "Empowered",
		ArmyReaper:    "Reaper",
		ArmyTwoKings:  "Two Kings",
		ArmyAnimals:   "Animals",
	}
	colorNames = map[Color]string{
		ColorWhite: "white",
		ColorBlack: "black",
	}
	pieceNames = map[PieceName]string{
		PieceNameClassicKing:     "Classic King",
		PieceNameNemesisQueen:    "Nemesis",
		PieceNameNemesisPawn:     "Nemesis Pawn",
		PieceNameEmpoweredQueen:  "Abdicated Queen",
		PieceNameEmpoweredBishop: "Empowered Bishop",
		PieceNameEmpoweredKnight: "Empowered Knight",
		PieceNameEmpoweredRook:   "Empowered Rook",
		PieceNameReaperQueen:     "Reaper",
		PieceNameReaperRook:      "Ghost",
		PieceNameTwoKingsKing:    "Warrior King",
		PieceNameAnimalsQueen:    "Jungle Queen",
		PieceNameAnimalsBishop:   "Tiger",
		PieceNameAnimalsKnight:   "Wild Horse",
		PieceNameAnimalsRook:     "Elephant",
	}
)

// InvalidPiece is the default value for Piece. It represents "no piece".
var InvalidPiece = Piece{}

// ColorIdx returns 0 for white, 1 for black
func ColorIdx(color Color) int {
	if color == ColorWhite {
		return 0
	}
	return 1
}

// OtherColor returns black for white and vice versa
func OtherColor(color Color) Color {
	if color == ColorWhite {
		return ColorBlack
	}
	return ColorWhite
}

func (t PieceType) String() string {
	if name, found := basicTypeNames[t]; found {
		return name
	}
	return string(int(t))
}

func (a Army) String() string {
	if name, found := armyNames[a]; found {
		return name
	}
	return string(int(a))
}

func (c Color) String() string {
	if name, found := colorNames[c]; found {
		return name
	}
	return string(int(c))
}

func (p PieceName) String() string {
	if name, found := pieceNames[p]; found {
		return name
	}
	return PieceType(uint8(p) & typeMask).String()
}

// NewPiece returns a piece with the given properties
func NewPiece(pieceType PieceType, army Army, color Color) Piece {
	return Piece{repr: uint8(pieceType) | uint8(army) | uint8(color)}
}

// WithArmy returns a copied Piece with the army set to the given value.
func (p Piece) WithArmy(army Army) Piece {
	return NewPiece(p.Type(), army, p.Color())
}

// Type returns the piece type of the receiver.
func (p Piece) Type() PieceType {
	return PieceType(p.repr & typeMask)
}

// Army returns the piece army of the receiver.
func (p Piece) Army() Army {
	return Army(p.repr & armyMask)
}

// Name returns one of the PieceName* constants. It's a combination of the Army
// and Type, but pieces which aren't special are converted to ArmyBasic.
func (p Piece) Name() PieceName {
	if _, found := pieceNames[PieceName(p.repr&^colorMask)]; found {
		return PieceName(p.repr &^ colorMask)
	}
	return PieceName(p.repr & typeMask)
}

// Color returns the piece color of the receiver.
func (p Piece) Color() Color {
	return Color(p.repr & colorMask)
}

func (p Piece) String() string {
	return fmt.Sprintf("%v %s", p.Color(), p.Name().String())
}
