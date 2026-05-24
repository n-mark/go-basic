package repository

import (
	"sync"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
)

type Entity interface {
	SerializeToJson() ([]byte, error)
}

type Repo struct {
	Players     []player.Player
	muPlayers   sync.Mutex
	Figures     []chess.Figure
	muFigures   sync.Mutex
	Games       []game.Game
	muGames     sync.Mutex
	Moves       []player.Move
	muMoves     sync.Mutex
	ChessBoards []chess.ChessBoard
	muBoards    sync.Mutex
	storage     StorageProvider
}

func New(sp StorageProvider) *Repo {
    repo := &Repo{
		Players:     make([]player.Player, 0),
		Figures:     make([]chess.Figure, 0),
		Games:       make([]game.Game, 0),
		Moves:       make([]player.Move, 0),
		ChessBoards: make([]chess.ChessBoard, 0),
		storage:     sp,
	}

	sp.Load(repo)
	return repo
}

func (r *Repo) DefineAndAdd(e Entity) {
	if x, ok := e.(player.Player); ok {
		r.muPlayers.Lock()
		r.Players = append(r.Players, x)
		r.storage.Save(x)
		r.muPlayers.Unlock()
	}
	if x, ok := e.(chess.Figure); ok {
		r.muFigures.Lock()
		r.Figures = append(r.Figures, x)
		r.storage.Save(x)
		r.muFigures.Unlock()
	}
	if x, ok := e.(game.Game); ok {
		r.muGames.Lock()
		r.Games = append(r.Games, x)
		r.storage.Save(x)
		r.muGames.Unlock()
	}
	if x, ok := e.(player.Move); ok {
		r.muMoves.Lock()
		r.Moves = append(r.Moves, x)
		r.storage.Save(x)
		r.muMoves.Unlock()
	}
	if x, ok := e.(chess.ChessBoard); ok {
		r.muBoards.Lock()
		r.ChessBoards = append(r.ChessBoards, x)
		r.storage.Save(x)
		r.muBoards.Unlock()
	}
}

func (r *Repo) ConsumeFromChan(data chan Entity) {
	for e := range data {
		r.DefineAndAdd(e)
	}
}

func (r *Repo) CloseStorage() {
	r.storage.Close()
}