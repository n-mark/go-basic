package service

import (
	"fmt"
	"strconv"

	"example.com/go-basic/packages/internal/models/game"
)

func StartGame() {
	var sizeStr string
	var player1 string
	var player2 string
	fmt.Print("Введите размер доски: ")
	fmt.Scan(&sizeStr)
	fmt.Print("Введите имя игрока №1: ")
	fmt.Scan(&player1)
	fmt.Print("Введите имя игрока №2: ")
	fmt.Scan(&player2)

	size, _ := strconv.Atoi(sizeStr)

	game := game.NewGame(size, player1, player2)
	renderScreen(game)
}

func renderScreen(game *game.Game) {
	fmt.Println(game.GetPlayer1())
	renderChessBoard(game.GetChessBoard())
	fmt.Println(game.GetPlayer2())
}
