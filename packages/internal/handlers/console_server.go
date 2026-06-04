package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strconv"

	"example.com/go-basic/packages/internal/service"
)

type ConsoleServer struct {
	svc *service.GameServiceNew
	render service.Render
}

// ServerMessage - сообщение от сервера клиенту.
// Type определяет, что клиент должен сделать:
//   - "menu"   : показать меню (Title, Items) и прислать обратно Selection
//   - "prompt" : показать prompt (Title) и прислать обратно Input
//   - "board"  : отрендерить доску (Board), без ответа
//   - "info"   : показать Message, без ответа
//   - "bye"    : показать Message и завершить соединение
type ServerMessage struct {
	Type    string   `json:"type"`
	Title   string   `json:"title,omitempty"`
	Items   []string `json:"items,omitempty"`
	Board   string   `json:"board,omitempty"`
	Message string   `json:"message,omitempty"`
}

// ClientMessage - ответ клиента серверу.
type ClientMessage struct {
	Index int    `json:"index,omitempty"`
	Input string `json:"input,omitempty"`
}

func NewConsoleServer(svc *service.GameServiceNew) *ConsoleServer {
	return &ConsoleServer{svc: svc, render: &service.ConsoleRender{}}
}

type clientSession struct {
	conn    net.Conn
	encoder *json.Encoder
	decoder *json.Decoder
}

func (s *clientSession) send(msg ServerMessage) error {
	return s.encoder.Encode(msg)
}

func (s *clientSession) recv() (ClientMessage, error) {
	var m ClientMessage
	err := s.decoder.Decode(&m)
	return m, err
}

func (s *clientSession) showMenu(title string, items []string) (int, error) {
	if err := s.send(ServerMessage{Type: "menu", Title: title, Items: items}); err != nil {
		return 0, err
	}
	m, err := s.recv()
	if err != nil {
		return 0, err
	}
	return m.Index, nil
}

func (s *clientSession) prompt(title string) (string, error) {
	if err := s.send(ServerMessage{Type: "prompt", Title: title}); err != nil {
		return "", err
	}
	m, err := s.recv()
	if err != nil {
		return "", err
	}
	return m.Input, nil
}

func (s *clientSession) info(msg string) error {
	return s.send(ServerMessage{Type: "info", Message: msg})
}

func (s *clientSession) renderBoard(board string) error {
	return s.send(ServerMessage{Type: "board", Board: board})
}

func (cs *ConsoleServer) handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	session := &clientSession{
		conn:    conn,
		encoder: json.NewEncoder(conn),
		decoder: json.NewDecoder(reader),
	}

	for {
		idx, err := session.showMenu("Главное меню", []string{
			"Новая игра",
			"Подключиться к существующей игре",
			"Выход",
		})
		if err != nil {
			return
		}

		switch idx {
		case 0:
			if err := cs.handleNewGame(session); err != nil {
				return
			}
		case 1:
			if err := cs.handleJoinGame(session); err != nil {
				return
			}
		case 2:
			session.send(ServerMessage{Type: "bye", Message: "До свидания!"})
			return
		}
	}
}

func (cs *ConsoleServer) handleNewGame(s *clientSession) error {
	sizeStr, err := s.prompt("Введите размер доски")
	if err != nil {
		return err
	}
	size, convErr := strconv.Atoi(sizeStr)
	if convErr != nil || size <= 0 {
		return s.info("Некорректный размер доски")
	}

	player1, err := s.prompt("Имя игрока №1 (белые)")
	if err != nil {
		return err
	}
	player2, err := s.prompt("Имя игрока №2 (черные)")
	if err != nil {
		return err
	}

	id := cs.svc.StartWebGame(size, player1, player2)
	if err := s.info(fmt.Sprintf("Создана игра ID=%d", id)); err != nil {
		return err
	}

	return cs.gameLoop(s, id)
}

func (cs *ConsoleServer) handleJoinGame(s *clientSession) error {
	idStr, err := s.prompt("Введите ID игры")
	if err != nil {
		return err
	}
	id, convErr := strconv.ParseInt(idStr, 10, 64)
	if convErr != nil {
		return s.info("Некорректный ID")
	}

	if _, ok := cs.svc.Games[id]; !ok {
		return s.info(fmt.Sprintf("Игра с ID=%d не найдена", id))
	}

	if err := s.info(fmt.Sprintf("Подключение к игре ID=%d", id)); err != nil {
		return err
	}

	return cs.gameLoop(s, id)
}

func (cs *ConsoleServer) gameLoop(s *clientSession, gameId int64) error {
	g, ok := cs.svc.Games[gameId]
	if !ok {
		return s.info("Игра не найдена")
	}

	for {
		if err := s.renderBoard(cs.render.RenderLayout(g)); err != nil {
			return err
		}

		if g.IsOver() {
			msg := "Игра окончена."
			if w := g.GetWinner(); w != nil {
				msg = fmt.Sprintf("Игра окончена! Победил %s", w.GetName())
			}
			if err := s.info(msg); err != nil {
				return err
			}
			return nil
		}

		idx, err := s.showMenu(
			fmt.Sprintf("Игра ID=%d. Ход: %s", gameId, g.CurrentPlayerName()),
			[]string{"Сделать ход", "Автоход", "Сдаться", "Выйти в главное меню"},
		)
		if err != nil {
			return err
		}

		switch idx {
		case 0:
			if err := cs.doMove(s, gameId); err != nil {
				return err
			}
		case 1:
			if err := cs.doAutoMove(s, gameId); err != nil {
				return err
			}
		case 2:
			cs.svc.Surrender(gameId)
			if err := s.info("Вы сдались."); err != nil {
				return err
			}
		case 3:
			return nil
		}
	}
}

func (cs *ConsoleServer) doMove(s *clientSession, gameId int64) error {
	g := cs.svc.Games[gameId]

	from, err := s.prompt("С какой клетки (например, E2)")
	if err != nil {
		return err
	}
	to, err := s.prompt("На какую клетку (например, E4)")
	if err != nil {
		return err
	}

	valid, errMsg := g.ValidateMove(from, to)
	if !valid {
		return s.info("Ошибка: " + errMsg)
	}

	cs.svc.Move(gameId, from, to)
	return nil
}

func (cs *ConsoleServer) doAutoMove(s *clientSession, gameId int64) error {
	g := cs.svc.Games[gameId]

	countStr, err := s.prompt("Сколько автоходов сделать?")
	if err != nil {
		return err
	}
	count, convErr := strconv.Atoi(countStr)
	if convErr != nil || count < 1 {
		return s.info("Введите корректное число автоходов")
	}

	currentPlayer := g.CurrentPlayerName()
	cs.svc.AutoMove(gameId, currentPlayer, count)
	return s.info(fmt.Sprintf("Запланировано %d автоходов для %s", count, currentPlayer))
}

func (cs *ConsoleServer) RunConsoleListener() {
	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}

	fmt.Println("Console server started on :8081")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go cs.handleClient(conn)
	}
}
