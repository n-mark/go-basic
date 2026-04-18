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

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
)

// MenuItem представляет пункт меню
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
		// Если не получилось - используем обычный ввод
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
	fmt.Println("\x1b[90m(↑↓←→ или W/S/A/D для выбора, Enter для подтверждения)\x1b[0m")
}

func showMenu(items []MenuItem) int {
	selected := 0

	for {
		renderMenu(selected, items)
		key := readKey()

		switch key {
		case "\x1b[A", "w", "W", "й", "Й": // вверх
			selected = (selected - 1 + len(items)) % len(items)
		case "\x1b[B", "s", "S", "ы", "Ы": // вниз
			selected = (selected + 1) % len(items)
		case "\x1b[C", "d", "D", "в", "В": // вправо
			selected = (selected + 1) % len(items)
		case "\x1b[D", "a", "A", "ф", "Ф": // влево
			selected = (selected - 1 + len(items)) % len(items)
		case "\n", "\r": // Enter
			return selected
		case "1":
			return 0
		case "2":
			return 1
		case "3":
			return 2
		}

		// Перемещаем курсор вверх для перерисовки
		fmt.Printf("\x1b[%dA", 2)
	}
}

func performAutoMove(chessGame *game.Game) (string, string, bool) {
	moves := chessGame.GetValidMovesForCurrentPlayer()
	if len(moves) == 0 {
		return "", "", false
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(moves))
	fromSig := moves[randomIndex].From

	board := chessGame.GetChessBoard()
	targets := board.GetPossibleTargets(fromSig)
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
	fmt.Printf("\nДумаю... (%.0f сек)\n", delay.Seconds())
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
	fmt.Print("Сколько ходов сделать? (1-10): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	count, err := strconv.Atoi(input)
	if err != nil || count < 1 || count > 10 {
		fmt.Println("Введите число от 1 до 10")
		fmt.Print("Нажмите Enter для продолжения...")
		reader.ReadString('\n')
		return true
	}

	for i := 0; i < count && !chessGame.IsOver(); i++ {
		fromSig, toSig, ok := performAutoMove(chessGame)
		if !ok {
			fmt.Println("Нет доступных ходов!")
			break
		}
		valid, _ := chessGame.ValidateMove(fromSig, toSig)
		if !valid {
			i--
			continue
		}
		chessGame.MakeMove(fromSig, toSig)
		fmt.Printf("Автоход %d/%d: %s -> %s\n", i+1, count, fromSig, toSig)
	}

	fmt.Print("Нажмите Enter для продолжения...")
	reader.ReadString('\n')
	return true
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
		renderScreen(chessGame)

		currentPlayer := chessGame.CurrentPlayer()
		color := "белыми"
		if currentPlayer.GetFiguresColor() == chess.ColorBlack {
			color = "черными"
		}

		fmt.Printf("\n\x1b[1mХод игрока %s (%s)\x1b[0m\n", currentPlayer, color)
		selected := showMenu(menuItems)
		continueGame := menuItems[selected].Action(chessGame, reader)
		if !continueGame {
			break
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

func renderScreen(game *game.Game) {
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf(" Игрок %s vs %s \n", game.GetPlayer1(), game.GetPlayer2())
	fmt.Println("═══════════════════════════════════════")
	fmt.Println(game.GetPlayer1())
	renderChessBoard(game.GetChessBoard())
	fmt.Println(game.GetPlayer2())
}
