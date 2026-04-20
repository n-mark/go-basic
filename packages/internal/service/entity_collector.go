package service

import (
	"example.com/go-basic/packages/internal/interfaces"
	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
	"example.com/go-basic/packages/internal/repository"
)

type EntityCollector struct {
	EntityList []interfaces.Entity
	Repo       repository.Repo
}

func (e *EntityCollector) CreateRandomEntities() {
	p := player.Player{}
	cb := chess.ChessBoard{}
	f := chess.Figure{PieceColor: "black"}
	g := game.Game{}
	m := game.Move{}

	e.EntityList = append(e.EntityList, p, cb, f, g, m)

	for _, el := range e.EntityList {
		e.Repo.DefineAndAdd(el)
	}
}
