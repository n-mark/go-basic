package main

import (
	"fmt"
	"strconv"
	"strings"
)

const (
    darkCell   = "\x1b[48;2;205;133;63m"
    lightCell  = "\x1b[48;2;245;222;179m"
    blackPiece = "\x1b[38;2;80;40;20m"
    whitePiece = "\x1b[38;2;255;255;255m"
    reset      = "\x1b[0m"
)

var whiteFigures = []rune{'♖', '♘', '♗', '♕', '♔', '♗', '♘', '♖'}
var blackFigures = []rune{'♜', '♞', '♝', '♛', '♚', '♝', '♞', '♜'}

func getFigure(figures []rune, col int) rune {
	return figures[col%len(figures)]
}

func getStartingPiece(row, col, size int) rune {
	switch row {
	case 0:
		return getFigure(blackFigures, col)
	case 1:
		return '♟'
	case size - 2:
		return '♙'
	case size - 1:
		return getFigure(whiteFigures, col)
	}
	return 0
}

func main() {
	var sizeStr string
	var player1 string
	var player2 string
	fmt.Print("Введите размер доски: ")
	fmt.Scan(&sizeStr)
	size, _ := strconv.Atoi(sizeStr)
	sizeSize := len(sizeStr)
	fmt.Print("Введите имя игрока №1: ")
	fmt.Scan(&player1)
	fmt.Print("Введите имя игрока №2: ")
	fmt.Scan(&player2)

	fmt.Println(player1)

	fmt.Print(strings.Repeat(" ", sizeSize))
	for i := 0; i < size; i++ {
		fmt.Print(string(rune('A' + i%26)))
	}

	fmt.Println()

	for i := 0; i < size; i++ {
		var rowNum = i + 1
		fmt.Print(strings.Repeat(" ", sizeSize-len(strconv.Itoa(rowNum))))
		fmt.Print(rowNum)
		for j := 0; j < size; j++ {
			cell := getStartingPiece(i, j, size)
			bg := lightCell
			if (i+j)%2 == 0 {
				bg = darkCell
			}
			if cell == 0 {
				fmt.Print(bg + " " + reset)
			} else if i <= 1 {
				fmt.Print(bg + blackPiece + string(cell) + reset)
			} else {
				fmt.Print(bg + whitePiece + string(cell) + reset)
			}
		}
		fmt.Print(rowNum)
		fmt.Println()
	}

	fmt.Print(strings.Repeat(" ", sizeSize))
	for i := 0; i < size; i++ {
		fmt.Print(string(rune('A' + i%26)))
	}

	fmt.Println()

	fmt.Println(player2)
}
