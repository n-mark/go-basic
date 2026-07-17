package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorageProvider struct {
	db *pgxpool.Pool
}

func NewPgStore(db *pgxpool.Pool) *PostgresStorageProvider {
	return &PostgresStorageProvider{db: db}
}

func (p *PostgresStorageProvider) Load(r *Repo) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	playersQuery := "SELECT id, name, figures_color FROM players WHERE 1=1"
	players, err := scanRows[player.Player](ctx, p.db, playersQuery)
	if err != nil {
		slog.Error("pgstorage: load players failed", "error", err)
	}
	r.Players = players

	figuresQuery := "SELECT id, symbol, piece_color FROM figures WHERE 1=1"
	figures, err := scanRows[chess.Figure](ctx, p.db, figuresQuery)
	if err != nil {
		slog.Error("pgstorage: load figures failed", "error", err)
	}
	r.Figures = figures

	movesQuery := "SELECT id, time_took, position_from, position_to, figure FROM moves WHERE 1=1"
	moves, err := scanRows[player.Move](ctx, p.db, movesQuery)
	if err != nil {
		slog.Error("pgstorage: load moves failed", "error", err)
	}
	r.Moves = moves
}

// SavePlayers вставляет/обновляет всех игроков в одной транзакции (upsert).
func (p *PostgresStorageProvider) SavePlayers(players []player.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
		INSERT INTO players (id, name, figures_color)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET name = EXCLUDED.name,
		    figures_color = EXCLUDED.figures_color`
	for _, pl := range players {
		if _, err := tx.Exec(ctx, q, pl.ID, pl.Name, pl.FiguresColor); err != nil {
			return fmt.Errorf("upsert player %d: %w", pl.ID, err)
		}
	}

	return tx.Commit(ctx)
}

// SaveFigures вставляет/обновляет все фигуры в одной транзакции (upsert).
func (p *PostgresStorageProvider) SaveFigures(figures []chess.Figure) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
		INSERT INTO figures (id, symbol, piece_color)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET symbol = EXCLUDED.symbol,
		    piece_color = EXCLUDED.piece_color`
	for _, f := range figures {
		if _, err := tx.Exec(ctx, q, f.ID, f.Symbol, f.PieceColor); err != nil {
			return fmt.Errorf("upsert figure %d: %w", f.ID, err)
		}
	}

	return tx.Commit(ctx)
}

// SaveMoves вставляет/обновляет все ходы в одной транзакции (upsert).
func (p *PostgresStorageProvider) SaveMoves(moves []player.Move) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	tx, err := p.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	const q = `
		INSERT INTO moves (id, time_took, position_from, position_to, figure)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO UPDATE
		SET time_took     = EXCLUDED.time_took,
		    position_from = EXCLUDED.position_from,
		    position_to   = EXCLUDED.position_to,
		    figure        = EXCLUDED.figure`
	for _, m := range moves {
		if _, err := tx.Exec(ctx, q,
			m.ID,
			m.TimeTook.Nanoseconds(),
			m.PositionFrom,
			m.PositionTo,
			m.Figure,
		); err != nil {
			return fmt.Errorf("upsert move %d: %w", m.ID, err)
		}
	}

	return tx.Commit(ctx)
}

func (p *PostgresStorageProvider) Close() {
	p.db.Close()
}

func ConnectPG(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	dialCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	pool, err := pgxpool.New(dialCtx, dsn)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(dialCtx); err != nil {
		pool.Close()
		return nil, err
	}
	return pool, nil
}

func scanRows[T any](ctx context.Context, db *pgxpool.Pool, query string) ([]T, error) {
	var result []T
	err := pgxscan.Select(ctx, db, &result, query)
	if err != nil {
		return nil, err
	}
	return result, nil
}
