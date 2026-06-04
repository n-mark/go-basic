package handlers

import (
	"example.com/go-basic/packages/internal/repository"
	"example.com/go-basic/packages/internal/service"
)

type CommonServer struct {
	webServer     *Server
	consoleServer *ConsoleServer
	repo          *repository.Repo
	entityService *service.EntityService
}

func NewCommonServer() *CommonServer {
	gameService := service.NewGameServiceNew()
	storage := repository.NewLocalStorageProvider()
	repo := repository.New(storage)
	entityService := service.NewEntityService(repo)
	ws := InitServer(gameService, entityService)
	cs := NewConsoleServer(gameService)

	return &CommonServer{webServer: ws, consoleServer: cs, repo: repo, entityService: entityService}
}

func (s *CommonServer) Run() {
	defer s.repo.CloseStorage()
	go s.consoleServer.RunConsoleListener()
	s.webServer.RunServer()
}
