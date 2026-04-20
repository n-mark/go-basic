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
	autoMoveCount map[string]int
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
		autoMoveCount: make(map[string]int),
	}
}

func (g *Game) StartGame() {
	g.gameStart = time.Now()
}

func (g *Game) StopGame() {
	g.gameEnd = time.Now()
}

func (g *Game) SetWinner(p *player.Player) {
	g.winner = p
}

func (g *Game) CurrentPlayer() *player.Player {
	if g.isPlayer1Turn {
		return g.player1
	}
	return g.player2
}

func (g *Game) CurrentPlayerName() string {
	return g.CurrentPlayer().GetName()
}

func (g *Game) IsPlayer1Turn() bool {
	return g.isPlayer1Turn
}

func (g *Game) ForceSetTurn(isPlayer1 bool) {
	g.isPlayer1Turn = isPlayer1
}

func (g *Game) IsOver() bool {
	return g.isOver
}

func (g *Game) GetWinner() *player.Player {
	return g.winner
}

func (g *Game) MakeMove(fromSig, toSig string) {
	fromKey := strings.ToUpper(fromSig)
	toKey := strings.ToUpper(toSig)

	toFigure := g.chessBoard.GetFigureBySignature(toKey)
	if toFigure != nil {
		g.CurrentPlayer().SetFigureTook(*toFigure)

		if '♔' == toFigure.Symbol || '♚' == toFigure.Symbol {
			g.isOver = true;
			g.SetWinner(g.CurrentPlayer())
			g.StopGame()
		}
	}

	g.chessBoard.MoveBySignature(fromKey, toKey)
	g.isPlayer1Turn = !g.isPlayer1Turn
}

func (g *Game) ValidateMove(fromSig, toSig string) (bool, string) {
	fromKey := strings.ToUpper(fromSig)
	toKey := strings.ToUpper(toSig)

	if !g.chessBoard.CellExists(fromKey) {
		return false, "Клетка " + fromSig + " не существует"
	}

	if !g.chessBoard.CellExists(toKey) {
		return false, "Клетка " + toSig + " не существует"
	}

	if !g.chessBoard.HasFigure(fromKey) {
		return false, "На клетке " + fromSig + " нет фигуры"
	}

	figure := g.chessBoard.GetFigureBySignature(fromKey)
	currentColor := g.CurrentPlayer().GetFiguresColor()
	if figure.GameColor != currentColor {
		return false, "Эта фигура не принадлежит вам"
	}

	return true, ""
}

func (g *Game) GetPlayer1() *player.Player {
	return g.player1
}

func (g *Game) GetPlayer2() *player.Player {
	return g.player2
}

func (g *Game) GetChessBoard() *chess.ChessBoard {
	return g.chessBoard
}

func (g *Game) GetValidMovesForCurrentPlayer() ([]string, []string) {
	color := g.CurrentPlayer().GetFiguresColor()
	return g.chessBoard.GetAllValidMoves(color)
}

func (g *Game) Surrender() {
	g.isOver = true
	otherPlayer := g.player2
	if g.isPlayer1Turn {
		otherPlayer = g.player2
	} else {
		otherPlayer = g.player1
	}
	g.SetWinner(otherPlayer)
	g.StopGame()
}

func (g *Game) SetAutoMoveCount(playerName string, count int) {
	g.autoMoveCount[playerName] = count
}

func (g *Game) GetAutoMoveCount(playerName string) int {
	return g.autoMoveCount[playerName]
}

func (g *Game) DecrementAutoMove(playerName string) bool {
	if g.autoMoveCount[playerName] > 0 {
		g.autoMoveCount[playerName]--
		return g.autoMoveCount[playerName] > 0
	}
	return false
}

func (g *Game) HasAutoMovePending(playerName string) bool {
	return g.autoMoveCount[playerName] > 0
}

func (g *Game) ClearAutoMoves() {
	g.autoMoveCount = make(map[string]int)
}
