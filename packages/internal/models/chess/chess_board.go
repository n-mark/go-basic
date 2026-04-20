package chess

import (
	"strconv"
	"strings"
)

const (
	// Используем 256-цветовую палитру для лучшей совместимости
	darkCell   = "\x1b[48;5;130m" // коричневый
	lightCell  = "\x1b[48;5;223m" // пшеничный
	blackPiece = "\x1b[38;5;52m"  // тёмно-коричневый для текста
	whitePiece = "\x1b[38;5;231m" // белый для текста
	Reset      = "\x1b[0m"

	ColorWhite = "white"
	ColorBlack = "black"
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
	Cells  map[string]*cell
}

func NewChessBoard(size int) *ChessBoard {
	c := &ChessBoard{Size: size}
	c.Layout, c.Cells = initCells(size)

	return c
}

func initCells(size int) ([][]*cell, map[string]*cell) {
	layout := make([][]*cell, size)
	cellsMap := make(map[string]*cell, size*size)

	for i := 0; i < size; i++ {
		row := make([]*cell, size)
		for j := 0; j < size; j++ {
			figure := getStartingPiece(i, j, size)
			color := lightCell

			if (i+j)%2 == 0 {
				color = darkCell
			}

			cellSignature := getSignature(i, j)

			c := &cell{Figure: figure,
				Color:     color,
				xIndex:    i,
				yIndex:    j,
				signature: cellSignature,
			}

			cellsMap[cellSignature] = c
			row[j] = c
		}
		layout[i] = row
	}

	return layout, cellsMap
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
		return &Figure{Symbol: getFigure(blackFigures, col), PieceColor: blackPiece, GameColor: ColorBlack}
	case 1:
		return &Figure{Symbol: '♟', PieceColor: blackPiece, GameColor: ColorBlack}
	case size - 2:
		return &Figure{Symbol: '♙', PieceColor: whitePiece, GameColor: ColorWhite}
	case size - 1:
		return &Figure{Symbol: getFigure(whiteFigures, col), PieceColor: whitePiece, GameColor: ColorWhite}
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

func (b *ChessBoard) MoveBySignature(fromSig, toSig string) {
	fromKey := strings.ToUpper(fromSig)
	toKey := strings.ToUpper(toSig)
	b.Cells[toKey].Figure = b.Cells[fromKey].Figure
	b.Cells[fromKey].Figure = nil
}

func (b *ChessBoard) GetFigureBySignature(sig string) *Figure {
	return b.Cells[strings.ToUpper(sig)].Figure
}

func (b *ChessBoard) CellExists(sig string) bool {
	_, exists := b.Cells[strings.ToUpper(sig)]
	return exists
}

func (b *ChessBoard) HasFigure(sig string) bool {
	cell := b.Cells[strings.ToUpper(sig)]
	return cell != nil && cell.Figure != nil
}

// GetAllValidMoves возвращает все возможные ходы для указанного цвета
func (b *ChessBoard) GetAllValidMoves(color string) ([]string, []string) {
	var fromMoves []string
	var toMoves []string
	for sig, cell := range b.Cells {
		if cell.Figure != nil && cell.Figure.GameColor == color {
			fromMoves = append(fromMoves, sig)
		} else if (cell.Figure != nil && cell.Figure.GameColor != color) || cell.Figure == nil {
			toMoves = append(toMoves, sig)
		}
	}
	return fromMoves, toMoves
}
