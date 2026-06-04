package repository

import (
	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
)

type StorageProvider interface {
	Load(r *Repo)
	SavePlayers([]player.Player) error
	SaveFigures([]chess.Figure) error
	SaveMoves([]player.Move) error
	Close()
}
