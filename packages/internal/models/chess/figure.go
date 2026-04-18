package chess

type Figure struct {
	Symbol     rune
	PieceColor string // ANSI color code for rendering
	GameColor  string // "white" or "black" for game logic
}
