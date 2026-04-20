package chess

import "encoding/json"

type Figure struct {
	Symbol     rune
	PieceColor string // ANSI color code for rendering
	GameColor  string // "white" or "black" for game logic
}

func (f Figure) SerializeToJson() string {
	jsonData, _ := json.Marshal(f)
	return string(jsonData)
}