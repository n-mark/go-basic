package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
	"example.com/go-basic/packages/internal/service"
	"github.com/gin-gonic/gin"
)

func timeDurationFromInt(ns int64) time.Duration {
	return time.Duration(ns)
}

// ItemHandler реализует REST CRUD-операции по сущностям.
type ItemHandler struct {
	entityService *service.EntityService
}

func NewItemHandler(entityService *service.EntityService) *ItemHandler {
	return &ItemHandler{entityService: entityService}
}

// RegisterRoutes регистрирует маршруты CRUD для каждой сущности.
func (h *ItemHandler) RegisterRoutes(api *gin.RouterGroup) {
	// players
	api.POST("/player", h.CreatePlayer)
	api.PUT("/player/:id", h.UpdatePlayer)
	api.GET("/players", h.ListPlayers)
	api.GET("/player/:id", h.GetPlayer)
	api.DELETE("/player/:id", h.DeletePlayer)

	// figures
	api.POST("/figure", h.CreateFigure)
	api.PUT("/figure/:id", h.UpdateFigure)
	api.GET("/figures", h.ListFigures)
	api.GET("/figure/:id", h.GetFigure)
	api.DELETE("/figure/:id", h.DeleteFigure)

	// moves
	api.POST("/move", h.CreateMove)
	api.PUT("/move/:id", h.UpdateMove)
	api.GET("/moves", h.ListMoves)
	api.GET("/move/:id", h.GetMove)
	api.DELETE("/move/:id", h.DeleteMove)
}

// ----------------- Players -----------------

type playerDTO struct {
	Name         string         `json:"name" binding:"required"`
	FiguresColor string         `json:"figures_color" binding:"required"`
	FiguresTook  []chess.Figure `json:"figures_took"`
	Moves        []player.Move  `json:"moves"`
}

func (d playerDTO) toModel() player.Player {
	return player.Player{
		Name:         d.Name,
		FiguresColor: d.FiguresColor,
		FiguresTook:  d.FiguresTook,
		Moves:        d.Moves,
	}
}

// CreatePlayer godoc
// @Summary      Создать игрока
// @Description  Создает нового игрока в хранилище
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        player body     playerDTO true  "Данные игрока"
// @Success      201    {object} map[string]int
// @Failure      400    {object} Err
// @Failure      500    {object} Err
// @Router       /player [post]
func (h *ItemHandler) CreatePlayer(ctx *gin.Context) {
	var dto playerDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}
	id, err := h.entityService.CreatePlayer(dto.toModel())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Err{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdatePlayer godoc
// @Summary      Обновить игрока
// @Description  Обновляет данные игрока по ID
// @Tags         players
// @Accept       json
// @Produce      json
// @Param        id     path     int        true  "ID игрока"
// @Param        player body     playerDTO  true  "Обновленные данные"
// @Success      200    {object} map[string]int
// @Failure      400    {object} Err
// @Failure      404    {object} Err
// @Router       /player/{id} [put]
func (h *ItemHandler) UpdatePlayer(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	var dto playerDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}
	if err := h.entityService.UpdatePlayer(id, dto.toModel()); err != nil {
		ctx.JSON(http.StatusNotFound, Err{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

// ListPlayers godoc
// @Summary      Список игроков
// @Description  Возвращает всех сохраненных игроков
// @Tags         players
// @Produce      json
// @Success      200  {array}  player.Player
// @Router       /players [get]
func (h *ItemHandler) ListPlayers(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.entityService.ListPlayers())
}

// GetPlayer godoc
// @Summary      Получить игрока
// @Description  Возвращает игрока по ID
// @Tags         players
// @Produce      json
// @Param        id   path     int  true  "ID игрока"
// @Success      200  {object} player.Player
// @Failure      400  {object} Err
// @Failure      404  {object} Err
// @Router       /player/{id} [get]
func (h *ItemHandler) GetPlayer(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	p, found := h.entityService.GetPlayer(id)
	if !found {
		ctx.JSON(http.StatusNotFound, Err{Err: "player не найден"})
		return
	}
	ctx.JSON(http.StatusOK, p)
}

// DeletePlayer godoc
// @Summary      Удалить игрока
// @Description  Удаляет игрока по ID
// @Tags         players
// @Produce      json
// @Param        id   path     int  true  "ID игрока"
// @Success      200  {object} map[string]int
// @Failure      400  {object} Err
// @Failure      404  {object} Err
// @Router       /player/{id} [delete]
func (h *ItemHandler) DeletePlayer(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	if err := h.entityService.DeletePlayer(id); err != nil {
		ctx.JSON(http.StatusNotFound, Err{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

// ----------------- Figures -----------------

type figureDTO struct {
	Symbol     string `json:"symbol" binding:"required"`
	PieceColor string `json:"piece_color" binding:"required"`
	GameColor  string `json:"game_color" binding:"required"`
}

func (d figureDTO) toModel() chess.Figure {
	var sym rune
	for _, r := range d.Symbol {
		sym = r
		break
	}
	return chess.Figure{
		Symbol:     sym,
		PieceColor: d.PieceColor,
		GameColor:  d.GameColor,
	}
}

// CreateFigure godoc
// @Summary      Создать фигуру
// @Description  Создает новую фигуру
// @Tags         figures
// @Accept       json
// @Produce      json
// @Param        figure body     figureDTO true  "Данные фигуры"
// @Success      201    {object} map[string]int
// @Failure      400    {object} Err
// @Failure      500    {object} Err
// @Router       /figure [post]
func (h *ItemHandler) CreateFigure(ctx *gin.Context) {
	var dto figureDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}
	id, err := h.entityService.CreateFigure(dto.toModel())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Err{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateFigure godoc
// @Summary      Обновить фигуру
// @Description  Обновляет фигуру по ID
// @Tags         figures
// @Accept       json
// @Produce      json
// @Param        id     path     int        true  "ID фигуры"
// @Param        figure body     figureDTO  true  "Новые данные"
// @Success      200    {object} map[string]int
// @Failure      400    {object} Err
// @Failure      404    {object} Err
// @Router       /figure/{id} [put]
func (h *ItemHandler) UpdateFigure(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	var dto figureDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}
	if err := h.entityService.UpdateFigure(id, dto.toModel()); err != nil {
		ctx.JSON(http.StatusNotFound, Err{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

// ListFigures godoc
// @Summary      Список фигур
// @Description  Возвращает все сохраненные фигуры
// @Tags         figures
// @Produce      json
// @Success      200  {array}  chess.Figure
// @Router       /figures [get]
func (h *ItemHandler) ListFigures(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.entityService.ListFigures())
}

// GetFigure godoc
// @Summary      Получить фигуру
// @Description  Возвращает фигуру по ID
// @Tags         figures
// @Produce      json
// @Param        id   path     int  true  "ID фигуры"
// @Success      200  {object} chess.Figure
// @Failure      400  {object} Err
// @Failure      404  {object} Err
// @Router       /figure/{id} [get]
func (h *ItemHandler) GetFigure(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	f, found := h.entityService.GetFigure(id)
	if !found {
		ctx.JSON(http.StatusNotFound, Err{Err: "figure не найдена"})
		return
	}
	ctx.JSON(http.StatusOK, f)
}

// DeleteFigure godoc
// @Summary      Удалить фигуру
// @Description  Удаляет фигуру по ID
// @Tags         figures
// @Produce      json
// @Param        id   path     int  true  "ID фигуры"
// @Success      200  {object} map[string]int
// @Failure      400  {object} Err
// @Failure      404  {object} Err
// @Router       /figure/{id} [delete]
func (h *ItemHandler) DeleteFigure(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	if err := h.entityService.DeleteFigure(id); err != nil {
		ctx.JSON(http.StatusNotFound, Err{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

// ----------------- Moves -----------------

type moveDTO struct {
	TimeTook     int64  `json:"time_took"`
	PositionFrom string `json:"position_from" binding:"required"`
	PositionTo   string `json:"position_to" binding:"required"`
	Figure       string `json:"figure" binding:"required"`
}

func (d moveDTO) toModel() player.Move {
	var sym rune
	for _, r := range d.Figure {
		sym = r
		break
	}
	return player.Move{
		TimeTook:     timeDurationFromInt(d.TimeTook),
		PositionFrom: strings.ToUpper(d.PositionFrom),
		PositionTo:   strings.ToUpper(d.PositionTo),
		Figure:       sym,
	}
}

// CreateMove godoc
// @Summary      Создать ход
// @Description  Создает новый ход в хранилище
// @Tags         moves
// @Accept       json
// @Produce      json
// @Param        move body     moveDTO true  "Данные хода"
// @Success      201  {object} map[string]int
// @Failure      400  {object} Err
// @Failure      500  {object} Err
// @Router       /move [post]
func (h *ItemHandler) CreateMove(ctx *gin.Context) {
	var dto moveDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}
	id, err := h.entityService.CreateMove(dto.toModel())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, Err{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"id": id})
}

// UpdateMove godoc
// @Summary      Обновить ход
// @Description  Обновляет ход по ID
// @Tags         moves
// @Accept       json
// @Produce      json
// @Param        id   path     int     true  "ID хода"
// @Param        move body     moveDTO true  "Новые данные хода"
// @Success      200  {object} map[string]int
// @Failure      400  {object} Err
// @Failure      404  {object} Err
// @Router       /move/{id} [put]
func (h *ItemHandler) UpdateMove(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	var dto moveDTO
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: err.Error()})
		return
	}
	if err := h.entityService.UpdateMove(id, dto.toModel()); err != nil {
		ctx.JSON(http.StatusNotFound, Err{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

// ListMoves godoc
// @Summary      Список ходов
// @Description  Возвращает все сохраненные ходы
// @Tags         moves
// @Produce      json
// @Success      200  {array}  player.Move
// @Router       /moves [get]
func (h *ItemHandler) ListMoves(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, h.entityService.ListMoves())
}

// GetMove godoc
// @Summary      Получить ход
// @Description  Возвращает ход по ID
// @Tags         moves
// @Produce      json
// @Param        id   path     int  true  "ID хода"
// @Success      200  {object} player.Move
// @Failure      400  {object} Err
// @Failure      404  {object} Err
// @Router       /move/{id} [get]
func (h *ItemHandler) GetMove(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	m, found := h.entityService.GetMove(id)
	if !found {
		ctx.JSON(http.StatusNotFound, Err{Err: "move не найден"})
		return
	}
	ctx.JSON(http.StatusOK, m)
}

// DeleteMove godoc
// @Summary      Удалить ход
// @Description  Удаляет ход по ID
// @Tags         moves
// @Produce      json
// @Param        id   path     int  true  "ID хода"
// @Success      200  {object} map[string]int
// @Failure      400  {object} Err
// @Failure      404  {object} Err
// @Router       /move/{id} [delete]
func (h *ItemHandler) DeleteMove(ctx *gin.Context) {
	id, ok := parseID(ctx)
	if !ok {
		return
	}
	if err := h.entityService.DeleteMove(id); err != nil {
		ctx.JSON(http.StatusNotFound, Err{Err: err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"id": id})
}

// ----------------- helpers -----------------

func parseID(ctx *gin.Context) (int, bool) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, Err{Err: "невалидный id"})
		return 0, false
	}
	return id, true
}
