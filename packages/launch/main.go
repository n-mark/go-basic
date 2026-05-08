package main

import (
	"math/rand"

	"example.com/go-basic/packages/internal/service"
)

func main() {
	service.StartGame()

	ec := service.NewCollector()
	ec.Run(rand.Intn(3) + 3)
}
