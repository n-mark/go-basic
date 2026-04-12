package chess

import (
	"strconv"
)

const (
	darkCell   = "\x1b[48;2;205;133;63m"
	lightCell  = "\x1b[48;2;245;222;179m"
	blackPiece = "\x1b[38;2;80;40;20m"
	whitePiece = "\x1b[38;2;255;255;255m"
	Reset      = "\x1b[0m"
)

var whiteFigures = []rune{'♖', '♘', '♗', '♕', '♔', '♗', '♘', '♖'}
var blackFigures = []rune{'♜', '♞', '♝', '♛', '♚', '♝', '♞', '♜'}

type cell struct {
	Figure    *Figure
	Color     string
	xIndex    int
	yIndex    int
	signature string
}

type ChessBoard struct {
	Size   int
	Layout [][]*cell
}

func NewChessBoard(size int) *ChessBoard {
	c := &ChessBoard{Size: size}
	c.Layout = initCells(size)
	return c
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
			c := &cell{Figure: figure,
				Color:     color,
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
	return ColToLetter(col) + strconv.Itoa(row+1)
}

func ColToLetter(col int) string {
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

func (b ChessBoard) MaxColLabelLen() int {
	max := 0
	for i := 0; i < b.Size; i++ {
		l := len(ColToLetter(i))
		if l > max {
			max = l
		}
	}
	return max
}