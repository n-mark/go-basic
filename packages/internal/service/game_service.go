package service

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/term"

	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
)

type MenuItem struct {
	Label  string
	Action func(*game.Game, *bufio.Reader) bool
}

func clearScreen() {
	fmt.Print("\x1b[2J\x1b[H")
}

func readKey() string {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		return input
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 10)
	n, _ := os.Stdin.Read(buf)
	return string(buf[:n])
}

func renderMenu(selected int, items []MenuItem) {
	fmt.Print("  ")
	for i, item := range items {
		if i == selected {
			fmt.Printf("\x1b[7m %s \x1b[0m", item.Label)
		} else {
			fmt.Printf(" %s ", item.Label)
		}
		if i < len(items)-1 {
			fmt.Print("  ")
		}
	}
	fmt.Println()
	fmt.Println("\x1b[90m(стрелки или W/S/A/D, Enter - подтвердить)\x1b[0m")
}

func showMenu(items []MenuItem) int {
	selected := 0

	for {
		renderMenu(selected, items)
		key := readKey()

		switch key {
		case "\x1b[A", "w", "W", "й", "Й":
			selected = (selected - 1 + len(items)) % len(items)
		case "\x1b[B", "s", "S", "ы", "Ы":
			selected = (selected + 1) % len(items)
		case "\x1b[C", "d", "D", "в", "В":
			selected = (selected + 1) % len(items)
		case "\x1b[D", "a", "A", "ф", "Ф":
			selected = (selected - 1 + len(items)) % len(items)
		case "\n", "\r":
			return selected
		case "1":
			return 0
		case "2":
			return 1
		case "3":
			return 2
		}

		fmt.Printf("\x1b[%dA", 2)
	}
}

func performAutoMove(chessGame *game.Game, renderFunc func(*game.Game, string, float64)) (string, string, bool) {
	fromMoves, toMoves := chessGame.GetValidMovesForCurrentPlayer()
	if len(fromMoves) == 0 {
		return "", "", false
	}

	// rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(fromMoves))
	fromSig := fromMoves[randomIndex]

	// board := chessGame.GetChessBoard()
	targets := toMoves
	if len(targets) == 0 {
		return "", "", false
	}

	var randomTarget string
	for {
		idx := rand.Intn(len(targets))
		if targets[idx] != fromSig {
			randomTarget = targets[idx]
			break
		}
		if len(targets) == 1 {
			return "", "", false
		}
	}

	delay := time.Duration(2+rand.Intn(3)) * time.Second

	// Рендерим с статусом "думаю"
	if renderFunc != nil {
		renderFunc(chessGame, "thinking", delay.Seconds())
	}

	time.Sleep(delay)

	return fromSig, randomTarget, true
}

func handleMove(chessGame *game.Game, reader *bufio.Reader) bool {
	fmt.Println("\n--- Ввод хода ---")
	fmt.Print("С какой клетки (например, E2): ")

	fromSig, _ := reader.ReadString('\n')
	fromSig = strings.TrimSpace(fromSig)

	fmt.Print("На какую клетку (например, E4): ")
	toSig, _ := reader.ReadString('\n')
	toSig = strings.TrimSpace(toSig)

	if fromSig == "q" || fromSig == "й" {
		fmt.Println("Игра завершена.")
		chessGame.StopGame()
		return false
	}

	valid, errMsg := chessGame.ValidateMove(fromSig, toSig)
	if !valid {
		fmt.Printf("Ошибка: %s\n", errMsg)
		fmt.Print("Нажмите Enter для продолжения...")
		reader.ReadString('\n')
		return true
	}

	chessGame.MakeMove(fromSig, toSig)
	return true
}

func handleSurrender(chessGame *game.Game) bool {
	fmt.Println("\nВы точно хотите сдаться? (Enter - да, q/й - отмена)")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" || input == "й" {
		return true
	}
	chessGame.Surrender()
	return false
}

func handleAutoMove(chessGame *game.Game) bool {
	fmt.Println("\n--- Автоход ---")
	fmt.Print("Сколько ходов сделать? ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	count, err := strconv.Atoi(input)
	if err != nil || count < 1 {
		fmt.Println("Введите число автоходов")
		return true
	}

	currentPlayer := chessGame.CurrentPlayerName()
	chessGame.SetAutoMoveCount(currentPlayer, count)

	// Сразу переходим к выполнению автоходов
	return false
}

func renderWithTitle(g *game.Game, status string, seconds float64) {
	player1 := g.GetPlayer1()
	player2 := g.GetPlayer2()

	renderPlayerHeader(g, player2, status, seconds)
	renderChessBoard(g.GetChessBoard())
	renderPlayerHeader(g, player1, status, seconds)

	fmt.Println()
}

func executeAutoMoves(chessGame *game.Game) {
	// Выполняем автоходы по очереди — каждый текущий игрок делает свой ход
	for !chessGame.IsOver() {
		clearScreen()
		currentPlayer := chessGame.CurrentPlayerName()

		// Если у текущего игрока нет автоходов — выходим
		if !chessGame.HasAutoMovePending(currentPlayer) {
			break
		}

		fromSig, toSig, ok := performAutoMove(chessGame, renderWithTitle)

		if !ok {
			fmt.Println("Нет доступных ходов!")
			chessGame.DecrementAutoMove(currentPlayer)
			break
		}

		valid, _ := chessGame.ValidateMove(fromSig, toSig)
		if !valid {
			continue
		}

		chessGame.MakeMove(fromSig, toSig)
		fmt.Printf("\x1b[36mАвтоход %s: %s -> %s\x1b[0m\n", currentPlayer, fromSig, toSig)

		chessGame.DecrementAutoMove(currentPlayer)

		clearScreen()
		renderWithTitle(chessGame, "", 0)
	}
}

func StartGame() {
	var sizeStr string
	var player1 string
	var player2 string
	fmt.Print("Введите размер доски: ")
	fmt.Scan(&sizeStr)
	fmt.Print("Введите имя игрока №1 (белые): ")
	fmt.Scan(&player1)
	fmt.Print("Введите имя игрока №2 (черные): ")
	fmt.Scan(&player2)

	size, _ := strconv.Atoi(sizeStr)

	chessGame := game.NewGame(size, player1, player2)
	chessGame.StartGame()

	reader := bufio.NewReader(os.Stdin)

	menuItems := []MenuItem{
		{"Сделать ход", func(g *game.Game, r *bufio.Reader) bool {
			return handleMove(g, r)
		}},
		{"Сдаться", func(g *game.Game, r *bufio.Reader) bool {
			return handleSurrender(g)
		}},
		{"Автоход", func(g *game.Game, r *bufio.Reader) bool {
			return handleAutoMove(g)
		}},
	}

	for !chessGame.IsOver() {
		clearScreen()
		renderWithTitle(chessGame, "", 0)

		selected := showMenu(menuItems)
		actionResult := menuItems[selected].Action(chessGame, reader)

		// Пункт "Автоход" (индекс 2) возвращает false чтобы сразу выполнить автоходы
		if selected == 2 && !actionResult {
			executeAutoMoves(chessGame)
			if chessGame.IsOver() {
				break
			}
			continue // показываем меню снова
		}

		if !actionResult {
			break
		}

		if chessGame.IsOver() {
			break
		}

		// Проверяем, есть ли у текущего игрока автоходы
		currentPlayer := chessGame.CurrentPlayerName()
		if chessGame.HasAutoMovePending(currentPlayer) {
			executeAutoMoves(chessGame)
			if chessGame.IsOver() {
				break
			}
		}
	}

	clearScreen()
	renderScreen(chessGame)

	if chessGame.GetWinner() != nil {
		fmt.Printf("\nИгра окончена! Победил %s!\n", chessGame.GetWinner())
	} else {
		fmt.Println("\nИгра окончена.")
	}
}

func renderScreen(g *game.Game) {
	renderWithTitle(g, "", 0)
}

func renderPlayerHeader(g *game.Game, p *player.Player, status string, seconds float64) {
	var color string
		if ("black" == p.GetFiguresColor()) {
			color = "черные"
		} else {
			color = "белые"
		}

	fmt.Println("═══════════════════════════════════════")
	if g.CurrentPlayerName() == p.GetName() {
		fmt.Print("> ")
	} else {
		fmt.Print("  ")
	}
	fmt.Printf("%s (%s)", p.GetName(), color)
	if g.GetAutoMoveCount(p.GetName()) > 0 {
		fmt.Printf(" Кол-во автоходов: %d", g.GetAutoMoveCount(p.GetName()))
		fmt.Print(" ")
	}
	if status == "thinking" && g.CurrentPlayerName() == p.GetName() {
		fmt.Printf(" - думаю %.0f сек...", seconds)
	}
	fmt.Println()
	fmt.Println("═══════════════════════════════════════")
}
