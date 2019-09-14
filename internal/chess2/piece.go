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
)

const (
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
	pieceNames = map[uint8]string{
		uint8(ArmyClassic) | uint8(TypeKing):     "Classic King",
		uint8(ArmyNemesis) | uint8(TypeQueen):    "Nemesis",
		uint8(ArmyNemesis) | uint8(TypePawn):     "Nemesis Pawn",
		uint8(ArmyEmpowered) | uint8(TypeQueen):  "Abdicated Queen",
		uint8(ArmyEmpowered) | uint8(TypeBishop): "Empowered Bishop",
		uint8(ArmyEmpowered) | uint8(TypeKnight): "Empowered Knight",
		uint8(ArmyEmpowered) | uint8(TypeRook):   "Empowered Rook",
		uint8(ArmyReaper) | uint8(TypeQueen):     "Reaper",
		uint8(ArmyReaper) | uint8(TypeRook):      "Ghost",
		uint8(ArmyTwoKings) | uint8(TypeKing):    "Warrior King",
		uint8(ArmyAnimals) | uint8(TypeQueen):    "Jungle Queen",
		uint8(ArmyAnimals) | uint8(TypeBishop):   "Tiger",
		uint8(ArmyAnimals) | uint8(TypeKnight):   "Wild Horse",
		uint8(ArmyAnimals) | uint8(TypeRook):     "Elephant",
	}
)

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

// NewPiece returns a piece with the given properties
func NewPiece(pieceType PieceType, army Army, color Color) Piece {
	return Piece{repr: uint8(pieceType) | uint8(army) | uint8(color)}
}

// Type returns the piece type of the receiver.
func (p Piece) Type() PieceType {
	return PieceType(p.repr & typeMask)
}

// Army returns the piece army of the receiver.
func (p Piece) Army() Army {
	return Army(p.repr & armyMask)
}

// Color returns the piece color of the receiver.
func (p Piece) Color() Color {
	return Color(p.repr & colorMask)
}

// Name of the piece based on army and type
func (p Piece) Name() string {
	if name, found := pieceNames[p.repr&^colorMask]; found {
		return name
	}
	return fmt.Sprintf("%v %v", p.Army(), p.Type())
}

func (p Piece) String() string {
	return fmt.Sprintf("%v %s", p.Color(), p.Name())
}
