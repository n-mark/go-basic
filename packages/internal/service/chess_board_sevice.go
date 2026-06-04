package service

import (
	"strconv"
	"strings"

	"example.com/go-basic/packages/internal/models/chess"
)

func renderChessBoard(b *chess.ChessBoard) string {
	var sb strings.Builder
	sizeWidth := len(strconv.Itoa(b.Size))

	sb.WriteString(renderColumnHeader(b))

	for i := 0; i < b.Size; i++ {

		rowNum := b.Size - i

		sb.WriteString(strings.Repeat(" ", sizeWidth-len(strconv.Itoa(rowNum))))
		sb.WriteString(strconv.Itoa(rowNum))

		for j := 0; j < b.Size; j++ {
			cell := b.Layout[i][j]

			if cell.Figure == nil {
				sb.WriteString(cell.Color + " " + chess.Reset)
				continue
			}

			sb.WriteString(
				cell.Color +
					cell.Figure.PieceColor +
					string(cell.Figure.Symbol) +
					chess.Reset,
			)
		}

		sb.WriteString(strconv.Itoa(rowNum))
		sb.WriteString("\n")
	}

	sb.WriteString(renderColumnHeader(b))
	return sb.String()
}

func renderColumnHeader(b *chess.ChessBoard) string {
	var sb strings.Builder
	maxLen := b.MaxColLabelLen()

	// каждая строка = один уровень букв
	for level := 0; level < maxLen; level++ {

		sb.WriteString(strings.Repeat(" ", len(strconv.Itoa(b.Size))))

		for col := 0; col < b.Size; col++ {
			label := chess.ColToLetter(col)

			// выравниваем справа (как Excel)
			padding := maxLen - len(label)

			if level < padding {
				sb.WriteString(" ")
			} else {
				sb.WriteString(string(label[level-padding]))
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
