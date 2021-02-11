package game

type IBoard interface {
	PlaceMark(i, j int, mark Mark) error
	RotateLeft() IBoard
	GetLines() <-chan Line
	String() string
}

type IGame interface {
	GetPossibleMoves() <-chan Coordinates
	GetScore() float64
	SwapPlayers() IGame
	MakeMove(c Coordinates) error
}
