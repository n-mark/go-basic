package main

import (
	"math/rand"
	"time"

	"example.com/go-basic/packages/internal/interfaces"
	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
	"example.com/go-basic/packages/internal/repository"
	"example.com/go-basic/packages/internal/service"
)

func main() {
	service.StartGame()

	ec := service.EntityCollector{
		EntityList: make([]interfaces.Entity, 0),
		Repo: repository.Repo{
			Players:     make([]player.Player, 0),
			Figures:     make([]chess.Figure, 0),
			Games:       make([]game.Game, 0),
			Moves:       make([]game.Move, 0),
			ChessBoards: make([]chess.ChessBoard, 0),
		}}

	times := rand.Intn(3) + 3 // случайное число от 3 до 5

	for i := 0; i < times; i++ {
		ec.CreateRandomEntities()
		// fmt.Println(ec.Repo) uncomment this to check
		time.Sleep(10 * time.Second)
	}
}
