package game

import (
	"time"
	"example.com/go-basic/packages/internal/models/chess"
)

type Move struct {
	timeTook time.Duration
	positionFrom position
	positionTo position
	figure chess.Figure
}

type position struct {
	X int
	Y int
}
