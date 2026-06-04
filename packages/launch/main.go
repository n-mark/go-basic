package main

import (
	"example.com/go-basic/packages/internal/handlers"
)

func main() {
	s := handlers.NewCommonServer()
	s.Run()
}
