package service

import (
	"fmt"
	"strings"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
)

type WebRender struct{}

func (r *WebRender) RenderLayout(g *game.Game) string {
	board := g.GetChessBoard()
	size := board.Size

	var sb strings.Builder

	sb.WriteString(`<!doctype html>
<html>
<head>
<meta charset="utf-8" />
<title>Chess Viewer</title>
<style>
body {
    margin: 0;
    background: #e6e6e6;
    font-family: Arial, sans-serif;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
}
.wrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
}
.player {
    font-size: 18px;
    margin: 10px 0;
    color: #333;
}
.coords {
    width: 35vw;
    display: grid;
    grid-template-columns: repeat(`)
	sb.WriteString(fmt.Sprintf("%d", size))
	sb.WriteString(`, 1fr);
    text-align: center;
    font-weight: bold;
    color: #333;
    margin: 5px 0;
}
.middle {
    display: flex;
    flex-direction: row;
    align-items: stretch;
}
.side {
    display: grid;
    grid-template-rows: repeat(`)
	sb.WriteString(fmt.Sprintf("%d", size))
	sb.WriteString(`, 1fr);
    width: 30px;
    text-align: center;
    font-weight: bold;
    color: #333;
    align-items: center;
}
.board {
    display: grid;
    grid-template-columns: repeat(`)
	sb.WriteString(fmt.Sprintf("%d", size))
	sb.WriteString(`, 1fr);
    grid-template-rows: repeat(`)
	sb.WriteString(fmt.Sprintf("%d", size))
	sb.WriteString(`, 1fr);
    width: 35vw;
    height: 35vw;
    border: 3px solid #333;
    box-shadow: 0 10px 30px rgba(0,0,0,0.2);
}
.cell {
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 2.2vw;
    line-height: 1;
}
.white { background: #f0d9b5; }
.black { background: #b58863; }
.player.active {
    color: #c0392b;
    font-weight: bold;
}
.player.active::before {
    content: "> ";
}
.player.active::after {
    content: " <";
}
</style>
</head>
<body>
<div class="wrapper">
`)

	// верхний игрок (чёрные)
	p2Class := "player"
	if !g.IsPlayer1Turn() {
		p2Class = "player active"
	}
	sb.WriteString(fmt.Sprintf(`    <div class="%s">`, p2Class))
	sb.WriteString(htmlEscape(g.GetPlayer2().GetName()))
	sb.WriteString(" (черные)</div>\n")

	// верхние буквы
	sb.WriteString(`    <div class="coords">` + "\n")
	for j := 0; j < size; j++ {
		sb.WriteString("        <div>")
		sb.WriteString(chess.ColToLetter(j))
		sb.WriteString("</div>\n")
	}
	sb.WriteString("    </div>\n")

	// центр
	sb.WriteString(`    <div class="middle">` + "\n")

	// левые цифры
	sb.WriteString(`        <div class="side">` + "\n")
	for i := 0; i < size; i++ {
		sb.WriteString(fmt.Sprintf("            <div>%d</div>\n", size-i))
	}
	sb.WriteString("        </div>\n")

	// доска
	sb.WriteString(`        <div class="board">` + "\n")
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			cellColor := "white"
			if (i+j)%2 == 0 {
				cellColor = "black"
			}
			sig := chess.ColToLetter(j) + fmt.Sprintf("%d", size-i)
			symbol := ""
			if board.HasFigure(sig) {
				fig := board.GetFigureBySignature(sig)
				symbol = string(fig.Symbol)
			}
			sb.WriteString(fmt.Sprintf("            <div class=\"cell %s\">%s</div>\n", cellColor, symbol))
		}
	}
	sb.WriteString("        </div>\n")

	// правые цифры
	sb.WriteString(`        <div class="side">` + "\n")
	for i := 0; i < size; i++ {
		sb.WriteString(fmt.Sprintf("            <div>%d</div>\n", size-i))
	}
	sb.WriteString("        </div>\n")

	sb.WriteString("    </div>\n")

	// нижние буквы
	sb.WriteString(`    <div class="coords">` + "\n")
	for j := 0; j < size; j++ {
		sb.WriteString("        <div>")
		sb.WriteString(chess.ColToLetter(j))
		sb.WriteString("</div>\n")
	}
	sb.WriteString("    </div>\n")

	// нижний игрок (белые)
	p1Class := "player"
	if g.IsPlayer1Turn() {
		p1Class = "player active"
	}
	sb.WriteString(fmt.Sprintf(`    <div class="%s">`, p1Class))
	sb.WriteString(htmlEscape(g.GetPlayer1().GetName()))
	sb.WriteString(" (белые)</div>\n")

	sb.WriteString("</div>\n</body>\n</html>\n")

	return sb.String()
}

func htmlEscape(s string) string {
	r := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		`"`, "&quot;",
		"'", "&#39;",
	)
	return r.Replace(s)
}