package player

import (
	"encoding/json"
	"math"
	"strconv"
	"time"
)

type Move struct {
	TimeTook     time.Duration `json:"time_took"`
	PositionFrom string `json:"position_from"`
	PositionTo   string `json:"position_to"`
	Figure       rune `json:"figure"`
}

func (m Move) SerializeToJson() ([]byte, error) {
	return json.Marshal(m)
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