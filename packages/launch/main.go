package main

import (
	"math/rand"
	"time"

	"example.com/go-basic/packages/internal/service"
)

func main() {
	service.StartGame()

	ec := service.NewCollector()
	times := rand.Intn(3) + 3 // случайное число от 3 до 5
	for i := 0; i < times; i++ {
		ec.CreateRandomEntities()
		time.Sleep(1 * time.Second)
	}

	// uncomment this to check
	// ec.DisplayRepoContent()
}
