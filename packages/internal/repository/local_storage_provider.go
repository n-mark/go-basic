package repository

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
)

const (
	storageDir      = "storage"
	playersFileName = "players.json"
	figuresFileName = "figures.json"
	movesFileName   = "moves.json"
)

type LocalStorageProvider struct {
	dir string
}

func NewLocalStorageProvider() *LocalStorageProvider {
	if err := os.MkdirAll(storageDir, 0755); err != nil {
		panic(fmt.Errorf("не удалось создать директорию %s: %w", storageDir, err))
	}

	ls := &LocalStorageProvider{dir: storageDir}
	ls.ensureFile(playersFileName)
	ls.ensureFile(figuresFileName)
	ls.ensureFile(movesFileName)
	return ls
}

func (ls *LocalStorageProvider) ensureFile(name string) {
	path := filepath.Join(ls.dir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.WriteFile(path, []byte("[]"), 0644); err != nil {
			panic(fmt.Errorf("не удалось создать файл %s: %w", path, err))
		}
	}
}

func (ls *LocalStorageProvider) Load(r *Repo) {
	readJSON(filepath.Join(ls.dir, playersFileName), &r.Players)
	readJSON(filepath.Join(ls.dir, figuresFileName), &r.Figures)
	readJSON(filepath.Join(ls.dir, movesFileName), &r.Moves)
}

func readJSON(path string, v any) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("не удалось прочитать %s: %w", path, err))
	}
	if len(data) == 0 {
		return
	}
	if err := json.Unmarshal(data, v); err != nil {
		slog.Error("ошибка десериализации", "file", path, "error", err)
	}
}

func (ls *LocalStorageProvider) SavePlayers(data []player.Player) error {
	return ls.writeJSON(playersFileName, data)
}

func (ls *LocalStorageProvider) SaveFigures(data []chess.Figure) error {
	return ls.writeJSON(figuresFileName, data)
}

func (ls *LocalStorageProvider) SaveMoves(data []player.Move) error {
	return ls.writeJSON(movesFileName, data)
}

func (ls *LocalStorageProvider) writeJSON(fileName string, data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	path := filepath.Join(ls.dir, fileName)
	return os.WriteFile(path, bytes, 0644)
}

func (ls *LocalStorageProvider) Close() {}
