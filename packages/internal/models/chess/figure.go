package chess

import "encoding/json"

type Figure struct {
	ID         int    `json:"id"`
	Symbol     rune   `json:"symbol"`
	PieceColor string `json:"piece_color"`
	GameColor  string `json:"game_color"`
}

func (f Figure) GetID() int    { return f.ID }
func (f *Figure) SetID(id int) { f.ID = id }

func (f Figure) SerializeToJson() ([]byte, error) {
	return json.Marshal(f)
}
