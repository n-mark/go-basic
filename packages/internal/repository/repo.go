package repository

import (
	"strings"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
)

type Entity interface {
	SerializeToJson() string
}

type Repo struct {
	total       []Entity
	Players     []player.Player
	Figures     []chess.Figure
	Games       []game.Game
	Moves       []game.Move
	ChessBoards []chess.ChessBoard
}

func (r *Repo) DefineAndAdd(e Entity) {
	if x, ok := e.(player.Player); ok {
		r.Players = append(r.Players, x)
		r.total = append(r.total, x)
	}
	if x, ok := e.(chess.Figure); ok {
		r.Figures = append(r.Figures, x)
		r.total = append(r.total, x)
	}
	if x, ok := e.(game.Game); ok {
		r.Games = append(r.Games, x)
		r.total = append(r.total, x)
	}
	if x, ok := e.(game.Move); ok {
		r.Moves = append(r.Moves, x)
		r.total = append(r.total, x)
	}
	if x, ok := e.(chess.ChessBoard); ok {
		r.ChessBoards = append(r.ChessBoards, x)
		r.total = append(r.total, x)
	}
}

func (r Repo) String() string {
	var s strings.Builder

	for _, e := range r.total {
		s.WriteString(e.SerializeToJson())
		s.WriteString("\n")
	}

	return s.String()
}
