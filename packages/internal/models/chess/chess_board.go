package chess

import (
	"strconv"
	"strings"
)

const (
	darkCell   = "\x1b[48;2;205;133;63m"
	lightCell  = "\x1b[48;2;245;222;179m"
	blackPiece = "\x1b[38;2;80;40;20m"
	whitePiece = "\x1b[38;2;255;255;255m"
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
	c.Layout = initCells(size)
	c.Cells = make(map[string]*cell, size*size)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			c.Cells[c.Layout[i][j].signature] = c.Layout[i][j]
		}
	}
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

func (b *ChessBoard) HasWhiteKing() bool {
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.Layout[i][j].Figure != nil && b.Layout[i][j].Figure.Symbol == '♔' {
				return true
			}
		}
	}
	return false
}

func (b *ChessBoard) HasBlackKing() bool {
	for i := 0; i < b.Size; i++ {
		for j := 0; j < b.Size; j++ {
			if b.Layout[i][j].Figure != nil && b.Layout[i][j].Figure.Symbol == '♚' {
				return true
			}
		}
	}
	return false
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
func (b *ChessBoard) GetAllValidMoves(color string) []MoveInfo {
	var moves []MoveInfo
	for sig, cell := range b.Cells {
		if cell.Figure != nil && cell.Figure.GameColor == color {
			moves = append(moves, MoveInfo{From: sig})
		}
	}
	return moves
}

// GetPossibleTargets возвращает все возможные целевые клетки для фигуры
func (b *ChessBoard) GetPossibleTargets(fromSig string) []string {
	fromKey := strings.ToUpper(fromSig)
	cell := b.Cells[fromKey]
	if cell == nil || cell.Figure == nil {
		return nil
	}
	
	var targets []string
	for sig := range b.Cells {
		if sig != fromKey {
			targets = append(targets, sig)
		}
	}
	return targets
}

type MoveInfo struct {
	From string
	To   string
}
