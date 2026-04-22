package service

import (
	// "fmt"

	"fmt"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
	"example.com/go-basic/packages/internal/repository"
)

type EntityCollector struct {
	repo repository.Repo
}

func (e *EntityCollector) CreateRandomEntities() {
	p := player.Player{}
	cb := chess.ChessBoard{}
	f := chess.Figure{PieceColor: "black"}
	g := game.Game{}
	m := game.Move{}

	e.repo.DefineAndAdd(p)
	e.repo.DefineAndAdd(cb)
	e.repo.DefineAndAdd(f)
	e.repo.DefineAndAdd(g)
	e.repo.DefineAndAdd(m)
}

func NewCollector() *EntityCollector {
	ec := EntityCollector{
		repo: repository.Repo{
			Players:     make([]player.Player, 0),
			Figures:     make([]chess.Figure, 0),
			Games:       make([]game.Game, 0),
			Moves:       make([]game.Move, 0),
			ChessBoards: make([]chess.ChessBoard, 0),
		}}

	return &ec
}

func (e *EntityCollector) DisplayRepoContent() {
	fmt.Println("REPO CONTENT: ")
	fmt.Println(e.repo)
}