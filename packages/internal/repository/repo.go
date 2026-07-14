package repository

import (
	"context"
	"sync"

	"example.com/go-basic/packages/internal/history"
	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
)

type Entity interface {
	SerializeToJson() ([]byte, error)
}

type Repo struct {
	Players         []player.Player
	muPlayers       sync.Mutex
	playerIDCounter int
	Figures         []chess.Figure
	muFigures       sync.Mutex
	figureIDCounter int
	Games           []game.Game
	muGames         sync.Mutex
	Moves           []player.Move
	muMoves         sync.Mutex
	moveIDCounter   int
	ChessBoards     []chess.ChessBoard
	muBoards        sync.Mutex
	storage         StorageProvider
	history         history.Logger
}

func New(sp StorageProvider) *Repo {
	return NewWithHistory(sp, nil)
}

func NewWithHistory(sp StorageProvider, hl history.Logger) *Repo {
	repo := &Repo{
		Players:     make([]player.Player, 0),
		Figures:     make([]chess.Figure, 0),
		Games:       make([]game.Game, 0),
		Moves:       make([]player.Move, 0),
		ChessBoards: make([]chess.ChessBoard, 0),
		storage:     sp,
		history:     hl,
	}

	sp.Load(repo)
	repo.recalculateCounters()
	return repo
}

func (r *Repo) CloseStorage() {
	if r.storage != nil {
		r.storage.Close()
	}
}

func (r *Repo) recordHistory(entityType string, entityID int, action history.Action, oldVal, newVal any) {
	if r.history == nil {
		return
	}
	err := r.history.Log(context.Background(), history.ChangeRecord{
		EntityType: entityType,
		EntityID:   entityID,
		Action:     action,
		OldValue:   oldVal,
		NewValue:   newVal,
	})
	if err != nil {
		_ = err
	}
}
