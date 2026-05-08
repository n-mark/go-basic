package player

import (
	"encoding/json"
	"math"
	"strconv"
	"time"
)

type Move struct {
	TimeTook     time.Duration
	PositionFrom string
	PositionTo   string
	Figure       rune
}

func (m Move) SerializeToJson() string {
	jsonData, _ := json.Marshal(m)
	return string(jsonData)
}

func (m Move) String() string {
	return string(m.Figure) + " " + m.PositionFrom + " -> " + m.PositionTo + " (" + formatDuration(m.TimeTook, 1) + " с)"
}

func formatDuration(d time.Duration, decimals int) string {
	totalSeconds := float64(d) / float64(time.Second)
	factor := math.Pow(10, float64(decimals))
	rounded := math.Round(totalSeconds*factor) / factor
	return strconv.FormatFloat(rounded, 'f', decimals, 64)
}