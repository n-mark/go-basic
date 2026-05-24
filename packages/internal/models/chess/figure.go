package chess

import "encoding/json"

type Figure struct {
	Symbol     rune `json:"symbol"`
	PieceColor string `json:"piece_color"`
	GameColor  string `json:"game_color"`
}

func (f Figure) SerializeToJson() ([]byte, error) {
	return json.Marshal(f)
}