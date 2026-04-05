package player

import "example.com/go-basic/packages/internal/models/chess"

type Player struct {
	name string
	figuresColor string
	figuresTook []chess.Figure
}

func NewPlayer(name string, figuresColor string) *Player {
	p := &Player{name: name, figuresColor: figuresColor}
	return p
}

func (p *Player) SetFigureTook(f chess.Figure) {
	p.figuresTook = append(p.figuresTook, f);
}

func (p Player) GetFiguresTook() []chess.Figure {
	return p.figuresTook;
}

func (p Player) String() string {
	return p.name
}