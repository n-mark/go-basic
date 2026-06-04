package service

import "example.com/go-basic/packages/internal/models/game"

type Stringifier interface {
	StringifyGameLayout(g *game.Game, status string) string
}