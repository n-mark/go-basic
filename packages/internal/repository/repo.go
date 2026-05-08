package repository

import (
	"sync"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
)

type Entity interface {
	SerializeToJson() string
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
}

func New() *Repo {
	return &Repo{
		Players:     make([]player.Player, 0),
		Figures:     make([]chess.Figure, 0),
		Games:       make([]game.Game, 0),
		Moves:       make([]player.Move, 0),
		ChessBoards: make([]chess.ChessBoard, 0),
	}
}

func (r *Repo) DefineAndAdd(e Entity) {
	if x, ok := e.(player.Player); ok {
		r.muPlayers.Lock()
		r.Players = append(r.Players, x)
		r.muPlayers.Unlock()
	}
	if x, ok := e.(chess.Figure); ok {
		r.muFigures.Lock()
		r.Figures = append(r.Figures, x)
		r.muFigures.Unlock()
	}
	if x, ok := e.(game.Game); ok {
		r.muGames.Lock()
		r.Games = append(r.Games, x)
		r.muGames.Unlock()
	}
	if x, ok := e.(player.Move); ok {
		r.muMoves.Lock()
		r.Moves = append(r.Moves, x)
		r.muMoves.Unlock()
	}
	if x, ok := e.(chess.ChessBoard); ok {
		r.muBoards.Lock()
		r.ChessBoards = append(r.ChessBoards, x)
		r.muBoards.Unlock()
	}
}

func (r *Repo) ConsumeFromChan(data chan Entity) {
	for e := range data {
		r.DefineAndAdd(e)
	}
}