package service

import (
	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
	"example.com/go-basic/packages/internal/repository"
)

type EntityService struct {
	repo *repository.Repo
}

func NewEntityService(repo *repository.Repo) *EntityService {
	return &EntityService{repo: repo}
}

func (es *EntityService) CreatePlayer(p player.Player) (int, error) {
	return es.repo.AddPlayer(p)
}

func (es *EntityService) UpdatePlayer(id int, p player.Player) error {
	return es.repo.UpdatePlayer(id, p)
}

func (es *EntityService) DeletePlayer(id int) error {
	return es.repo.DeletePlayer(id)
}

func (es *EntityService) GetPlayer(id int) (player.Player, bool) {
	return es.repo.GetPlayer(id)
}

func (es *EntityService) ListPlayers() []player.Player {
	return es.repo.GetAllPlayers()
}

func (es *EntityService) CreateFigure(f chess.Figure) (int, error) {
	return es.repo.AddFigure(f)
}

func (es *EntityService) UpdateFigure(id int, f chess.Figure) error {
	return es.repo.UpdateFigure(id, f)
}

func (es *EntityService) DeleteFigure(id int) error {
	return es.repo.DeleteFigure(id)
}

func (es *EntityService) GetFigure(id int) (chess.Figure, bool) {
	return es.repo.GetFigure(id)
}

func (es *EntityService) ListFigures() []chess.Figure {
	return es.repo.GetAllFigures()
}

func (es *EntityService) CreateMove(m player.Move) (int, error) {
	return es.repo.AddMove(m)
}

func (es *EntityService) UpdateMove(id int, m player.Move) error {
	return es.repo.UpdateMove(id, m)
}

func (es *EntityService) DeleteMove(id int) error {
	return es.repo.DeleteMove(id)
}

func (es *EntityService) GetMove(id int) (player.Move, bool) {
	return es.repo.GetMove(id)
}

func (es *EntityService) ListMoves() []player.Move {
	return es.repo.GetAllMoves()
}
