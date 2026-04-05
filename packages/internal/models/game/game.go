package game

import (
	"time"
	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
)

type Game struct {
	player1 *player.Player
	player2 *player.Player
	chessBoard *chess.ChessBoard
	gameStart time.Time
	gameEnd time.Time
	winner *player.Player
}

func NewGame(boardSize int, player1Name string, player2Name string) *Game {
	board := chess.NewChessBoard(boardSize)
	player1 := player.NewPlayer(player1Name, "white")
	player2 := player.NewPlayer(player2Name, "black")
	return &Game{
		player1: player1,
		player2: player2,
		chessBoard: board,
		gameStart: time.Now(),
	}
}

func (game *Game) StartGame() {
	game.gameStart = time.Now()
}

func (game *Game) StopGame() {
	game.gameEnd = time.Now()
}

func (game *Game) SetWinner(p *player.Player) {
	game.winner = p;
}

func (game Game) GetPlayer1() *player.Player {
	return game.player1
}

func (game Game) GetPlayer2() *player.Player {
	return game.player2
}

func (game Game) GetChessBoard() chess.ChessBoard {
	return *game.chessBoard
}