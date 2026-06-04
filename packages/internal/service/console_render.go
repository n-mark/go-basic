package service

import (
	"strings"

	"example.com/go-basic/packages/internal/models/game"
)

type ConsoleRender struct{}

func (r *ConsoleRender) RenderLayout(g *game.Game) string {
	var sb strings.Builder
	player1 := g.GetPlayer1()
	player2 := g.GetPlayer2()

	sb.WriteString(renderPlayerHeader(g, player2, ""))
	sb.WriteString(renderChessBoard(g.GetChessBoard()))
	sb.WriteString(renderPlayerHeader(g, player1, ""))

	return sb.String()
}