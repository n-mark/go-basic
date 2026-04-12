package service

import (
	"fmt"
	"strconv"
	"strings"

	"example.com/go-basic/packages/internal/models/chess"
)

func renderChessBoard(b chess.ChessBoard) {
	sizeWidth := len(strconv.Itoa(b.Size))

	renderColumnHeader(b)

	for i := 0; i < b.Size; i++ {

		rowNum := i + 1

		fmt.Print(strings.Repeat(" ", sizeWidth-len(strconv.Itoa(rowNum))))
		fmt.Print(rowNum)

		for j := 0; j < b.Size; j++ {
			cell := b.Layout[i][j]

			if cell.Figure == nil {
				fmt.Print(cell.Color + " " + chess.Reset)
				continue
			}

			fmt.Print(
				cell.Color +
					cell.Figure.PieceColor +
					string(cell.Figure.Symbol) +
					chess.Reset,
			)
		}

		fmt.Print(rowNum)
		fmt.Println()
	}

	renderColumnHeader(b)
}


func renderColumnHeader(b chess.ChessBoard) {
		maxLen := b.MaxColLabelLen()

	// каждая строка = один уровень букв
	for level := 0; level < maxLen; level++ {

		fmt.Print(strings.Repeat(" ", len(strconv.Itoa(b.Size))))

		for col := 0; col < b.Size; col++ {
			label := chess.ColToLetter(col)

			// выравниваем справа (как Excel)
			padding := maxLen - len(label)

			if level < padding {
				fmt.Print(" ")
			} else {
				fmt.Print(string(label[level-padding]))
			}
		}
		fmt.Println()
	}
}