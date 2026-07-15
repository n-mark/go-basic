package repository

import (
	"context"
	"time"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStorageProvider struct {
	db *pgxpool.Pool
}

func NewPgStore(db *pgxpool.Pool) *PostgresStorageProvider {
	return &PostgresStorageProvider{db: db}
}

func (p *PostgresStorageProvider) Load(r *Repo) {
	//TODO implement me
	panic("implement me")
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
