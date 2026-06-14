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

// NewGame godoc
// @Summary      Новая игра
// @Description  Создает новую шахматную партию между двумя игроками
// @Tags         game
// @Accept       json
// @Produce      json
// @Param        request body     createGame  true  "Параметры новой игры"
// @Success      200     {object} response
// @Failure      400     {object} Err
// @Router       /game [post]
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

// DisplayBoard godoc
// @Summary      Отображение доски
// @Description  Возвращает текущее состояние доски в JSON или HTML (зависит от Accept-заголовка)
// @Tags         game
// @Produce      json
// @Produce      html
// @Param        id   path     int  true  "ID игры"
// @Success      200  {object} object
// @Failure      400  {object} map[string]string
// @Router       /game/{id} [get]
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

// Move godoc
// @Summary      Сделать ход
// @Description  Выполняет ход указанного игрока в активной партии
// @Tags         game
// @Accept       json
// @Produce      json
// @Param        id      path     int   true  "ID игры"
// @Param        request body     move  true  "Параметры хода"
// @Success      200     {object} response
// @Failure      400     {object} Err
// @Router       /game/{id}/move [post]
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

// Surrender godoc
// @Summary      Сдаться
// @Description  Игрок сдается, партия завершается победой соперника
// @Tags         game
// @Accept       json
// @Produce      json
// @Param        id      path     int        true  "ID игры"
// @Param        request body     surrender  true  "Кто сдается"
// @Success      200     {object} response
// @Failure      400     {object} Err
// @Router       /game/{id}/surrender [post]
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

// StopGame godoc
// @Summary      Остановить игру
// @Description  Останавливает партию без объявления победителя
// @Tags         game
// @Produce      json
// @Param        id   path     int  true  "ID игры"
// @Success      200  {object} response
// @Failure      400  {object} Err
// @Router       /game/{id}/stop [post]
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

// AutoMove godoc
// @Summary      Автоход
// @Description  Запускает указанное количество автоматических ходов от лица игрока
// @Tags         game
// @Accept       json
// @Produce      json
// @Param        id      path     int      true  "ID игры"
// @Param        request body     autoMove true  "Параметры автоматических ходов"
// @Success      200     {object} response
// @Failure      400     {object} Err
// @Router       /game/{id}/automove [post]
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
