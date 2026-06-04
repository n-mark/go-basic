package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"example.com/go-basic/packages/internal/service"
	"github.com/gin-gonic/gin"
)

type GameHandler struct {
	gameService *service.GameServiceNew
	render      service.Render
}

type response struct {
	Success bool   `json:"success"`
	GameId  int64  `json:"game_id"`
	Message string `json:"message"`
}

type createGame struct {
	BoardSize     int    `json:"board_size"`
	PlayerOneName string `json:"player_one_name"`
	PlayerTwoName string `json:"player_two_name"`
}

type move struct {
	Player       string `json:"player"`
	PositionFrom string `json:"position_from"`
	PositionTo   string `json:"position_to"`
}

type surrender struct {
	Player string `json:"player"`
}

type autoMove struct {
	Player      string `json:"player"`
	MovesAmount int    `json:"moves_amount"`
}

func NewGameHandler(s *service.GameServiceNew) *GameHandler {
	return &GameHandler{gameService: s, render: &service.WebRender{}}
}

func (h *GameHandler) NewGame(ctx *gin.Context) {
	game := createGame{}

	err := ctx.ShouldBindJSON(&game)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}
	id := h.gameService.StartWebGame(game.BoardSize, game.PlayerOneName, game.PlayerTwoName)
	response := response{Success: true, GameId: id, Message: "Game created"}
	ctx.JSON(http.StatusOK, response)
}

func (h *GameHandler) DisplayBoard(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	accept := ctx.GetHeader("Accept")
	game := h.gameService.Games[id]

	switch {
	case strings.Contains(accept, "text/html"):
		ctx.Data(
			200,
			"text/html; charset=utf-8",
			[]byte(h.render.RenderLayout(game)),
		)

	default:
		ctx.JSON(500, game)
	}
}

func (h *GameHandler) Move(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	move := move{}

	err = ctx.ShouldBindJSON(&move)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}

	h.gameService.Move(id, move.PositionFrom, move.PositionTo)
	response := response{Success: true, GameId: id, Message: "move performed"}
	ctx.JSON(http.StatusOK, response)
}

func (h *GameHandler) Surrender(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	surrender := surrender{}

	err = ctx.ShouldBindJSON(&surrender)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}

	var resp response
	if ok, message := h.gameService.Surrender(id); ok {
		resp = response{Success: ok, GameId: id, Message: message}
	} else {
		resp = response{Success: ok, GameId: id, Message: message}
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *GameHandler) StopGame(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}

	var resp response
	if ok, message := h.gameService.Stop(id); ok {
		resp = response{Success: ok, GameId: id, Message: message}
	} else {
		resp = response{Success: ok, GameId: id, Message: message}
	}

	ctx.JSON(http.StatusOK, resp)
}

func (h *GameHandler) AutoMove(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid id"})
		return
	}
	move := autoMove{}

	err = ctx.ShouldBindJSON(&move)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}

	h.gameService.AutoMove(id, move.Player, move.MovesAmount)

	response := response{Success: true, GameId: id, Message: "automoves scheduled"}
	ctx.JSON(http.StatusOK, response)
}
