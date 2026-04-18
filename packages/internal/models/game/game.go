package game

import (
	"strings"
	"time"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
)

type Game struct {
	player1       *player.Player
	player2       *player.Player
	chessBoard    *chess.ChessBoard
	isPlayer1Turn bool
	isOver        bool
	gameStart     time.Time
	gameEnd       time.Time
	winner        *player.Player
}

func NewGame(boardSize int, player1Name string, player2Name string) *Game {
	board := chess.NewChessBoard(boardSize)
	player1 := player.NewWhitePlayer(player1Name)
	player2 := player.NewBlackPlayer(player2Name)
	return &Game{
		player1:       player1,
		player2:       player2,
		chessBoard:    board,
		isPlayer1Turn: true,
		gameStart:     time.Now(),
	}
}

func (game *Game) StartGame() {
	game.gameStart = time.Now()
}

func (game *Game) StopGame() {
	game.gameEnd = time.Now()
}

func (game *Game) SetWinner(p *player.Player) {
	game.winner = p
}

func (game *Game) CurrentPlayer() *player.Player {
	if game.isPlayer1Turn {
		return game.player1
	}
	return game.player2
}

func (game *Game) IsOver() bool {
	return game.isOver
}

func (game *Game) GetWinner() *player.Player {
	return game.winner
}

func (game *Game) MakeMove(fromSig, toSig string) {
	fromKey := strings.ToUpper(fromSig)
	toKey := strings.ToUpper(toSig)

	toFigure := game.chessBoard.GetFigureBySignature(toKey)
	if toFigure != nil {
		game.CurrentPlayer().SetFigureTook(*toFigure)
	}

	game.chessBoard.MoveBySignature(fromKey, toKey)

	if game.checkKingCaptured() {
		game.isOver = true
		game.SetWinner(game.CurrentPlayer())
		game.StopGame()
	}

	game.isPlayer1Turn = !game.isPlayer1Turn
}

func (game *Game) ValidateMove(fromSig, toSig string) (bool, string) {
	fromKey := strings.ToUpper(fromSig)
	toKey := strings.ToUpper(toSig)

	if !game.chessBoard.CellExists(fromKey) {
		return false, "Клетка " + fromSig + " не существует"
	}

	if !game.chessBoard.CellExists(toKey) {
		return false, "Клетка " + toSig + " не существует"
	}

	if !game.chessBoard.HasFigure(fromKey) {
		return false, "На клетке " + fromSig + " нет фигуры"
	}

	figure := game.chessBoard.GetFigureBySignature(fromKey)
	currentColor := game.CurrentPlayer().GetFiguresColor()
	if figure.GameColor != currentColor {
		return false, "Эта фигура не принадлежит вам"
	}

	return true, ""
}

func (game *Game) checkKingCaptured() bool {
	return !game.chessBoard.HasWhiteKing() || !game.chessBoard.HasBlackKing()
}

func (game *Game) GetPlayer1() *player.Player {
	return game.player1
}

func (game *Game) GetPlayer2() *player.Player {
	return game.player2
}

func (game *Game) GetChessBoard() *chess.ChessBoard {
	return game.chessBoard
}

// GetValidMovesForCurrentPlayer возвращает все валидные ходы текущего игрока
func (game *Game) GetValidMovesForCurrentPlayer() []chess.MoveInfo {
	color := game.CurrentPlayer().GetFiguresColor()
	return game.chessBoard.GetAllValidMoves(color)
}

// Surrender текущий игрок сдаётся
func (game *Game) Surrender() {
	game.isOver = true
	otherPlayer := game.player2
	if game.isPlayer1Turn {
		otherPlayer = game.player2
	} else {
		otherPlayer = game.player1
	}
	game.SetWinner(otherPlayer)
	game.StopGame()
}
