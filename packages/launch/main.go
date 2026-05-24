package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"

	"example.com/go-basic/packages/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signalChan
		fmt.Println("\nПолучен сигнал завершения, останавливаю приложение...")
		cancel()
	}()

	service.StartGame(ctx)

	ec := service.NewCollector()
	ec.Run(ctx, (rand.Intn(3)+3))
}
