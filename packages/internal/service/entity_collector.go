package service

import (
	"context"
	"sync"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
	"example.com/go-basic/packages/internal/repository"
)

type EntityCollector struct {
	repo *repository.Repo
	data chan repository.Entity
}

func (e *EntityCollector) Run(ctx context.Context, times int) {
	var producerWg sync.WaitGroup
	var consumerWg sync.WaitGroup

	for range times {
		producerWg.Add(1)
		go func() {
			defer producerWg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			e.CreateRandomEntitiesWithChan()
		}()
	}

	for range times {
		consumerWg.Add(1)
		go func() {
			defer consumerWg.Done()
			select {
			case <-ctx.Done():
				return
			default:
			}
			e.repo.ConsumeFromChan(e.data)
		}()
	}

	go func() {
		producerWg.Wait()
		close(e.data)
	}()

	consumerWgDone := make(chan struct{})
	go func() {
		consumerWg.Wait()
		close(consumerWgDone)
	}()

	select {
	case <-ctx.Done():
	case <-consumerWgDone:
	}
}

func (e *EntityCollector) CreateRandomEntities() {
	p := player.Player{}
	cb := chess.ChessBoard{}
	f := chess.Figure{PieceColor: "black"}
	g := game.Game{}
	m := player.Move{}

	e.repo.DefineAndAdd(p)
	e.repo.DefineAndAdd(cb)
	e.repo.DefineAndAdd(f)
	e.repo.DefineAndAdd(g)
	e.repo.DefineAndAdd(m)
}

func (e *EntityCollector) CreateRandomEntitiesWithChan() {
	p := player.Player{}
	cb := chess.ChessBoard{}
	f := chess.Figure{PieceColor: "black"}
	g := game.Game{}
	m := player.Move{}

	e.data <- p
	e.data <- cb
	e.data <- f
	e.data <- g
	e.data <- m
}

func NewCollector() *EntityCollector {
	ec := EntityCollector{
		repo: repository.New(),
		data: make(chan repository.Entity),
	}

	return &ec
}
