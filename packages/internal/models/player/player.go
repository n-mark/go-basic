package player

import (
	"encoding/json"

	"example.com/go-basic/packages/internal/models/chess"
)

type Player struct {
	name         string
	figuresColor string
	figuresTook  []chess.Figure
	moves        []Move
}

func NewPlayer(name string, figuresColor string) *Player {
	p := &Player{name: name, figuresColor: figuresColor, moves: make([]Move, 0)}
	return p
}

func NewWhitePlayer(name string) *Player {
	return NewPlayer(name, chess.ColorWhite)
}

func NewBlackPlayer(name string) *Player {
	return NewPlayer(name, chess.ColorBlack)
}

func (p *Player) SetFigureTook(f chess.Figure) {
	p.figuresTook = append(p.figuresTook, f)
}

func (p *Player) SetMove(m Move) {
	p.moves = append(p.moves, m)
}

func (p Player) GetMoves() []Move {
	return p.moves
}

func (p Player) GetLastNMoves(movesAmount int) []Move {
	if movesAmount >= len(p.moves) {
		return p.moves
	}
	start := len(p.moves) - movesAmount
	return p.moves[start:]
}

func (p Player) GetFiguresTook() []chess.Figure {
	return p.figuresTook
}

func (p Player) String() string {
	return p.name
}

func (p Player) GetFiguresColor() string {
	return p.figuresColor
}

func (p Player) GetName() string {
	return p.name
}

func (p Player) SerializeToJson() string {
	jsonData, _ := json.Marshal(p)
	return string(jsonData)
}
