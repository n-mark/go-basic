package handlers

import (
	// "net/http"

	"example.com/go-basic/packages/internal/service"
	"github.com/gin-gonic/gin"
)

type Server struct {
	gameHandler *GameHandler
	itemHandler *ItemHandler
}

func (s *Server) RunServer() {
	router := gin.Default()

	api := router.Group("/api")
	v1 := api.Group("/v1")

	game := v1.Group("/game")
	game.POST("", s.gameHandler.NewGame)
	game.POST("/:id/move", s.gameHandler.Move)
	game.POST("/:id/surrender", s.gameHandler.Surrender)
	game.POST("/:id/stop", s.gameHandler.StopGame)
	game.POST("/:id/automove", s.gameHandler.AutoMove)
	game.GET("/:id", s.gameHandler.DisplayBoard)

	router.Run(":8080")
}

func InitServer(s *service.GameServiceNew) *Server {
	return &Server{gameHandler: NewGameHandler(s)}
}
