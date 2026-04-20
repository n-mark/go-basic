package repository

import (
	"example.com/go-basic/packages/internal/interfaces"
	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
)

type Repo struct {
	Players     []player.Player
	Figures     []chess.Figure
	Games       []game.Game
	Moves       []game.Move
	ChessBoards []chess.ChessBoard
}

func (r *Repo) DefineAndAdd(e interfaces.Entity) {
	if x, ok := e.(player.Player); ok {
		r.Players = append(r.Players, x)
	}
	if x, ok := e.(chess.Figure); ok {
		r.Figures = append(r.Figures, x)
	}
	if x, ok := e.(game.Game); ok {
		r.Games = append(r.Games, x)
	}
	if x, ok := e.(game.Move); ok {
		r.Moves = append(r.Moves, x)
	}
	if x, ok := e.(chess.ChessBoard); ok {
		r.ChessBoards = append(r.ChessBoards, x)
	}
}
