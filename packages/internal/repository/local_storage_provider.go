package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/game"
	"example.com/go-basic/packages/internal/models/player"
)

const (
	DIR = "storage"
	PLAYERS_FILENAME = "players.json"
	FIGURES_FILENAME = "figures.json"
	GAMES_FILENAME = "games.json"
	MOVES_FILENAME = "moves.json"
	CHESSBOARD_FILENAME = "chessBoard.json"
)

type LocalStorageProvider struct {
	playerFile         *os.File
	playerFileIsEmpty  bool
	figureFile         *os.File
	figureFileIsEmpty  bool
	gameFile           *os.File
	gameFileIsEmpty    bool
	moveFile           *os.File
	moveFilesIsEmpty   bool
	chessBoard         *os.File
	chessBoardsIsEmpty bool
}

func NewLocalStorageProvider() *LocalStorageProvider {
	playerFile, playerFileIsEmpty := initFile(PLAYERS_FILENAME)
	figureFile, figureFileIsEmpty := initFile(FIGURES_FILENAME)
	gameFile, gameFileIsEmpty := initFile(GAMES_FILENAME)
	playerMove, playerMoveIsEmpty := initFile(MOVES_FILENAME)
	chessBoard, chessBoardIsEmpty := initFile(CHESSBOARD_FILENAME)

	return &LocalStorageProvider{playerFile: playerFile,
		playerFileIsEmpty:  playerFileIsEmpty,
		figureFile:         figureFile,
		figureFileIsEmpty:  figureFileIsEmpty,
		gameFile:           gameFile,
		gameFileIsEmpty:    gameFileIsEmpty,
		moveFile:           playerMove,
		moveFilesIsEmpty:   playerMoveIsEmpty,
		chessBoard:         chessBoard,
		chessBoardsIsEmpty: chessBoardIsEmpty}
}

func (ls *LocalStorageProvider) Save(e Entity) {
	if ls == nil {
		return
	}

	if _, ok := e.(player.Player); ok {
		ls.appendFile(e, ls.playerFile, &ls.playerFileIsEmpty)
		return
	}
	if _, ok := e.(chess.Figure); ok {
		ls.appendFile(e, ls.figureFile, &ls.figureFileIsEmpty)
	}
	if _, ok := e.(game.Game); ok {
		ls.appendFile(e, ls.gameFile, &ls.gameFileIsEmpty)
	}
	if _, ok := e.(player.Move); ok {
		ls.appendFile(e, ls.moveFile, &ls.moveFilesIsEmpty)
	}
	if _, ok := e.(chess.ChessBoard); ok {
		ls.appendFile(e, ls.chessBoard, &ls.chessBoardsIsEmpty)
	}
}

func (ls *LocalStorageProvider) appendFile(e Entity, file *os.File, isEmpty *bool) {
	if ls == nil || file == nil || isEmpty == nil {
		return
	}

	jsonBytes, err := e.SerializeToJson()
	if err != nil {
		slog.Error("Ошибка сериализации структуры", "error", err)
		return
	}

	if *isEmpty {
		*isEmpty = false
	} else {
		removeLastBracket(file)
		if _, err := file.WriteString(","); err != nil {
			slog.Error("Ошибка записи разделителя", "file", file.Name(), "error", err)
			return
		}
	}

	if _, err := file.Write(jsonBytes); err != nil {
		slog.Error("Ошибка записи в файл", "file", file.Name(), "error", err)
	}
}

func (ls *LocalStorageProvider) Load(r *Repo) {
	if !ls.playerFileIsEmpty {
		load(PLAYERS_FILENAME, &r.Players)
	}
	if !ls.figureFileIsEmpty {
		load(FIGURES_FILENAME, &r.Figures)
	}
	if !ls.gameFileIsEmpty {
		load(GAMES_FILENAME, &r.Games)
	}
	if !ls.moveFilesIsEmpty {
		load(MOVES_FILENAME, &r.Moves)
	}
	if !ls.chessBoardsIsEmpty {
		load(CHESSBOARD_FILENAME, &r.ChessBoards)
	}
}

func load(fileName string, v any) {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s", DIR, fileName))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, v); err != nil {
		slog.Error("Ошибка десериализации", "filename", v, "error", err)
		return
	}
}

func (ls *LocalStorageProvider) Close() {
	if ls == nil {
		return
	}

	ls.closeFile(ls.playerFile)
	ls.closeFile(ls.figureFile)
	ls.closeFile(ls.gameFile)
	ls.closeFile(ls.moveFile)
	ls.closeFile(ls.chessBoard)
}

func (ls *LocalStorageProvider) closeFile(file *os.File) {
	if file == nil {
		return
	}

	if _, err := file.WriteString("]"); err != nil {
		slog.Error("Ошибка записи закрывающей скобки", "file", file.Name(), "error", err)
	}

	if err := file.Close(); err != nil {
		slog.Error("Ошибка закрытия файла", "file", file.Name(), "error", err)
	}
}

func initFile(fileName string) (*os.File, bool) {
	if err := os.MkdirAll(DIR, 0755); err != nil {
		// Ошибка будет только если нет прав или проблемы с диском
		panic(err)
	}

	file, err := os.OpenFile(fmt.Sprintf("%s/%s", DIR, fileName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		slog.Error("Ошибка открытия файла", "error", err)
		panic(err)
	}

	info, err := file.Stat()
	if err != nil {
		slog.Error("Ошибка получения информации о файле", "error", err)
		file.Close()
		panic(err)
	}

	size := info.Size()
	isEmpty := size == 0

	if isEmpty {
		if _, err := file.WriteString("["); err != nil {
			slog.Error("Ошибка записи начала массива", "error", err)
			file.Close()
			panic(err)
		}
	}

	return file, isEmpty
}

func removeLastBracket(file *os.File) {
	info, err := file.Stat()
	if err != nil {
		slog.Error("Ошибка получения информации о файле", "error", err)
		file.Close()
		panic(err)
	}

	size := info.Size()

	buf := make([]byte, 1)
	if _, err := file.ReadAt(buf, size-1); err != nil {
		slog.Error("Ошибка чтения последнего символа файла", "error", err)
		file.Close()
		panic(err)
	}

	if buf[0] == ']' {
		if err := file.Truncate(size - 1); err != nil {
			slog.Error("Ошибка усечения файла", "error", err)
			file.Close()
			panic(err)
		}
	}

	if _, err := file.Seek(0, io.SeekEnd); err != nil {
		slog.Error("Ошибка перемещения курсора в конец файла", "error", err)
		file.Close()
		panic(err)
	}
}
