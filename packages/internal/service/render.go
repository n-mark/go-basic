package service

import "example.com/go-basic/packages/internal/models/game"

type Render interface {
	RenderLayout(game *game.Game) string
}