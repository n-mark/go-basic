package handlers

import (
	"example.com/go-basic/packages/internal/grpc"
	"example.com/go-basic/packages/internal/repository"
	"example.com/go-basic/packages/internal/service"
)

type CommonServer struct {
	webServer     *Server
	consoleServer *ConsoleServer
	grpc          *grpc.GrpcServer
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
	grpc := grpc.NewGrpcServer(entityService)

	return &CommonServer{webServer: ws, consoleServer: cs, grpc: grpc, repo: repo, entityService: entityService}
}

func (s *CommonServer) Run() {
	defer s.repo.CloseStorage()
	go s.consoleServer.RunConsoleListener()
	go grpc.RunGrpcListenerInParallel(s.grpc)
	s.webServer.RunServer()
}
