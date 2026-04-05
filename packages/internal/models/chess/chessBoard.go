package chess

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

type cell struct {
	figure    *Figure
	color     string
	xIndex    int
	yIndex    int
	signature string
}

type ChessBoard struct {
	size   int
	layout [][]*cell
}

func NewChessBoard(size int) *ChessBoard {
	c := &ChessBoard{size: size}
	c.layout = initCells(size)
	return c
}

func (b ChessBoard) RenderChessBoard() {

	sizeWidth := len(strconv.Itoa(b.size))

	b.renderColumnHeader()

	for i := 0; i < b.size; i++ {

		rowNum := i + 1

		fmt.Print(strings.Repeat(" ", sizeWidth-len(strconv.Itoa(rowNum))))
		fmt.Print(rowNum)

		for j := 0; j < b.size; j++ {
			cell := b.layout[i][j]

			if cell.figure == nil {
				fmt.Print(cell.color + " " + reset)
				continue
			}

			fmt.Print(
				cell.color +
					cell.figure.PieceColor +
					string(cell.figure.Symbol) +
					reset,
			)
		}

		fmt.Print(rowNum)
		fmt.Println()
	}

	b.renderColumnHeader()
}

func initCells(size int) [][]*cell {
	layout := make([][]*cell, size)

	for i := 0; i < size; i++ {
		row := make([]*cell, size)
		for j := 0; j < size; j++ {
			figure := getStartingPiece(i, j, size)
			color := lightCell

			if (i+j)%2 == 0 {
				color = darkCell
			}
			c := &cell{figure: figure,
				color:     color,
				xIndex:    i,
				yIndex:    j,
				signature: getSignature(i, j),
			}
			row[j] = c
		}
		layout[i] = row
	}

	return layout
}

func getSignature(row int, col int) string {
	return colToLetter(col) + strconv.Itoa(row+1)
}

func colToLetter(col int) string {
	result := ""
	c := col + 1 // 1-based
	for c > 0 {
		c--
		result = string(rune('A'+c%26)) + result
		c /= 26
	}
	return result
}

func getFigure(figures []rune, col int) rune {
	return figures[col%len(figures)]
}

func getStartingPiece(row, col, size int) *Figure {
	switch row {
	case 0:
		return &Figure{Symbol: getFigure(blackFigures, col), PieceColor: blackPiece}
	case 1:
		return &Figure{Symbol: '♟', PieceColor: blackPiece}
	case size - 2:
		return &Figure{Symbol: '♙', PieceColor: whitePiece}
	case size - 1:
		return &Figure{Symbol: getFigure(whiteFigures, col), PieceColor: whitePiece}
	}
	return nil
}

func (b ChessBoard) renderColumnHeader() {
	maxLen := b.maxColLabelLen()

	// каждая строка = один уровень букв
	for level := 0; level < maxLen; level++ {

		fmt.Print(strings.Repeat(" ", len(strconv.Itoa(b.size))))

		for col := 0; col < b.size; col++ {
			label := colToLetter(col)

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

func (b ChessBoard) maxColLabelLen() int {
	max := 0
	for i := 0; i < b.size; i++ {
		l := len(colToLetter(i))
		if l > max {
			max = l
		}
	}
	return max
}