package repository

import (
	"context"
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

	movesQuery := "SELECT id, time_took, position_from, position_to, figure FROM figures WHERE 1=1"
	moves, err := scanRows[player.Move](ctx, p.db, movesQuery)
	if err != nil {
		slog.Error("pgstorage: load figures failed", "error", err)
	}
	r.Moves = moves
}

func (p *PostgresStorageProvider) SavePlayers(players []player.Player) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresStorageProvider) SaveFigures(figures []chess.Figure) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresStorageProvider) SaveMoves(moves []player.Move) error {
	//TODO implement me
	panic("implement me")
}

func (p *PostgresStorageProvider) Close() {
	//TODO implement me
	panic("implement me")
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
