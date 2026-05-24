package player

import (
	"encoding/json"

	"example.com/go-basic/packages/internal/models/chess"
)

type Player struct {
	Name         string	`json:"name"`
	FiguresColor string `json:"figures_color"`
	FiguresTook  []chess.Figure `json:"figures_took"`
	Moves        []Move `json:"moves"`
}

func NewPlayer(name string, figuresColor string) *Player {
	p := &Player{Name: name, FiguresColor: figuresColor, Moves: make([]Move, 0)}
	return p
}

func NewWhitePlayer(name string) *Player {
	return NewPlayer(name, chess.ColorWhite)
}

func NewBlackPlayer(name string) *Player {
	return NewPlayer(name, chess.ColorBlack)
}

func (p *Player) SetFigureTook(f chess.Figure) {
	p.FiguresTook = append(p.FiguresTook, f)
}

func (p *Player) SetMove(m Move) {
	p.Moves = append(p.Moves, m)
}

func (p Player) GetMoves() []Move {
	return p.Moves
}

func (p Player) GetLastNMoves(movesAmount int) []Move {
	if movesAmount >= len(p.Moves) {
		return p.Moves
	}
	start := len(p.Moves) - movesAmount
	return p.Moves[start:]
}

func (p Player) GetFiguresTook() []chess.Figure {
	return p.FiguresTook
}

func (p Player) String() string {
	return p.Name
}

func (p Player) GetFiguresColor() string {
	return p.FiguresColor
}

func (p Player) GetName() string {
	return p.Name
}

func (p Player) SerializeToJson() ([]byte, error) {
	return json.Marshal(p)
}
