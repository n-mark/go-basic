package player

import "example.com/go-basic/packages/internal/models/chess"

type Player struct {
	name         string
	figuresColor string
	figuresTook  []chess.Figure
}

func NewPlayer(name string, figuresColor string) *Player {
	p := &Player{name: name, figuresColor: figuresColor}
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
