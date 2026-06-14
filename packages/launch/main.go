package main

import (
	"example.com/go-basic/packages/internal/handlers"
)

// @title           Chess Game API
// @version         1.0
// @description     API для управления шахматными партиями и сущностями (игроки, фигуры, ходы).
// @termsOfService  http://swagger.io/terms/
// @host      localhost:8080
// @BasePath  /api/v1
// @schemes   http
func main() {
	s := handlers.NewCommonServer()
	s.Run()
}
