package handlers

import "example.com/go-basic/packages/internal/service"


type CommonServer struct {
	webServer *Server
	consoleServer *ConsoleServer
}

func NewCommonServer() *CommonServer {
	gameService := service.NewGameServiceNew()
	ws := InitServer(gameService)
	cs := NewConsoleServer(gameService)

	return &CommonServer{webServer: ws, consoleServer: cs}
}

func (s *CommonServer) Run() {
	go s.consoleServer.RunConsoleListener()
	s.webServer.RunServer()
}
