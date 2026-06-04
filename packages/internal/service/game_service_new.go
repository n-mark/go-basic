package service

import (
	"math/rand"
	"time"

	"example.com/go-basic/packages/internal/models/game"
)

type GameServiceNew struct {
	idCounter int64
	Games     map[int64]*game.Game
}

func NewGameServiceNew() *GameServiceNew {
	return &GameServiceNew{Games: make(map[int64]*game.Game)}
}

func (gs *GameServiceNew) StartWebGame(size int, player1 string, player2 string) int64 {
	gs.idCounter++
	chessGame := game.NewGame(size, player1, player2)
	gs.Games[gs.idCounter] = chessGame
	chessGame.StartGame()
	return gs.idCounter
}

func (gs *GameServiceNew) GetGameById(gameId int64) (*game.Game, string) {
	if game, ok := gs.Games[gameId]; ok {
		return game, "success"
	} else {
		return nil, "game does not exist"
	}
}

func (gs *GameServiceNew) Move(gameId int64, positionFrom string, positionTo string) (bool, string) {
	if game, ok := gs.Games[gameId]; ok {
		game.MakeMove(positionFrom, positionTo)
		gs.executeAutoMoves(game)
		return ok, "success"
	} else {
		return ok, "game does not exist"
	}
}

func (gs *GameServiceNew) AutoMove(gameId int64, player string, movesAmount int) (bool, string) {
	if game, ok := gs.Games[gameId]; ok {
		game.SetAutoMoveCount(player, movesAmount)
		go gs.executeAutoMoves(game)
		return ok, "success"
	} else {
		return ok, "game does not exist"
	}
}

func (gs *GameServiceNew) Surrender(gameId int64) (bool, string) {
	if game, ok := gs.Games[gameId]; ok {
		game.Surrender()
		return ok, "surrended. game over"
	} else {
		return ok, "game does not exist"
	}
}

func (gs *GameServiceNew) Stop(gameId int64) (bool, string) {
	if game, ok := gs.Games[gameId]; ok {
		game.StopGame()
		return ok, "game stopped"
	} else {
		return ok, "game does not exist"
	}
}

func (gs *GameServiceNew) executeAutoMoves(chessGame *game.Game) {
	for !chessGame.IsOver() {
		currentPlayer := chessGame.CurrentPlayerName()

		if !chessGame.HasAutoMovePending(currentPlayer) {
			break
		}

		fromSig, toSig, ok := gs.performAutoMove(chessGame)

		if !ok {
			chessGame.DecrementAutoMove(currentPlayer)
			break
		}

		valid, _ := chessGame.ValidateMove(fromSig, toSig)
		if !valid {
			continue
		}

		chessGame.MakeMove(fromSig, toSig)
		chessGame.DecrementAutoMove(currentPlayer)
	}
}

func (gs *GameServiceNew) performAutoMove(chessGame *game.Game) (string, string, bool) {
	fromMoves, toMoves := chessGame.GetValidMovesForCurrentPlayer()
	if len(fromMoves) == 0 {
		return "", "", false
	}

	randomIndex := rand.Intn(len(fromMoves))
	fromSig := fromMoves[randomIndex]
	targets := toMoves

	if len(targets) == 0 {
		return "", "", false
	}

	var randomTarget string
	for {
		idx := rand.Intn(len(targets))
		if targets[idx] != fromSig {
			randomTarget = targets[idx]
			break
		}
		if len(targets) == 1 {
			return "", "", false
		}
	}

	delay := time.Duration(1+rand.Intn(6)) * time.Second
	time.After(delay)
	return fromSig, randomTarget, true
}
