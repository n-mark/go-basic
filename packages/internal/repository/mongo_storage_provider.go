package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"example.com/go-basic/packages/internal/models/chess"
	"example.com/go-basic/packages/internal/models/player"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	playersCollection = "players"
	figuresCollection = "figures"
	movesCollection   = "moves"
)

type MongoStorageProvider struct {
	client  *mongo.Client
	db      *mongo.Database
	players *mongo.Collection
	figures *mongo.Collection
	moves   *mongo.Collection
}

type MongoOptions struct {
	URI         string
	Database    string
	ConnectTime time.Duration
}

func NewMongoStorageProvider(ctx context.Context, opts MongoOptions) (*MongoStorageProvider, error) {
	if opts.URI == "" {
		opts.URI = "mongodb://localhost:27017"
	}
	if opts.Database == "" {
		opts.Database = "go_basic"
	}
	if opts.ConnectTime <= 0 {
		opts.ConnectTime = 10 * time.Second
	}

	connectCtx, cancel := context.WithTimeout(ctx, opts.ConnectTime)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(opts.URI))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}

	if err := client.Ping(connectCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	db := client.Database(opts.Database)
	return &MongoStorageProvider{
		client:  client,
		db:      db,
		players: db.Collection(playersCollection),
		figures: db.Collection(figuresCollection),
		moves:   db.Collection(movesCollection),
	}, nil
}

func (m *MongoStorageProvider) Database() *mongo.Database {
	return m.db
}

func (m *MongoStorageProvider) Load(r *Repo) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if players, err := m.loadPlayers(ctx); err != nil {
		slog.Error("mongo: load players failed", "error", err)
	} else {
		r.Players = players
	}

	if figures, err := m.loadFigures(ctx); err != nil {
		slog.Error("mongo: load figures failed", "error", err)
	} else {
		r.Figures = figures
	}

	if moves, err := m.loadMoves(ctx); err != nil {
		slog.Error("mongo: load moves failed", "error", err)
	} else {
		r.Moves = moves
	}
}

func (m *MongoStorageProvider) loadPlayers(ctx context.Context) ([]player.Player, error) {
	cursor, err := m.players.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	out := make([]player.Player, 0)
	if err := cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (m *MongoStorageProvider) loadFigures(ctx context.Context) ([]chess.Figure, error) {
	cursor, err := m.figures.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	out := make([]chess.Figure, 0)
	if err := cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (m *MongoStorageProvider) loadMoves(ctx context.Context) ([]player.Move, error) {
	cursor, err := m.moves.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	out := make([]player.Move, 0)
	if err := cursor.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (m *MongoStorageProvider) SavePlayers(data []player.Player) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := m.players.DeleteMany(ctx, bson.M{}); err != nil {
		return fmt.Errorf("mongo: delete players: %w", err)
	}
	if len(data) == 0 {
		return nil
	}

	docs := make([]any, 0, len(data))
	for _, p := range data {
		docs = append(docs, p)
	}
	if _, err := m.players.InsertMany(ctx, docs); err != nil {
		return fmt.Errorf("mongo: insert players: %w", err)
	}
	return nil
}

func (m *MongoStorageProvider) SaveFigures(data []chess.Figure) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := m.figures.DeleteMany(ctx, bson.M{}); err != nil {
		return fmt.Errorf("mongo: delete figures: %w", err)
	}
	if len(data) == 0 {
		return nil
	}

	docs := make([]any, 0, len(data))
	for _, f := range data {
		docs = append(docs, f)
	}
	if _, err := m.figures.InsertMany(ctx, docs); err != nil {
		return fmt.Errorf("mongo: insert figures: %w", err)
	}
	return nil
}

func (m *MongoStorageProvider) SaveMoves(data []player.Move) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := m.moves.DeleteMany(ctx, bson.M{}); err != nil {
		return fmt.Errorf("mongo: delete moves: %w", err)
	}
	if len(data) == 0 {
		return nil
	}

	docs := make([]any, 0, len(data))
	for _, mv := range data {
		docs = append(docs, mv)
	}
	if _, err := m.moves.InsertMany(ctx, docs); err != nil {
		return fmt.Errorf("mongo: insert moves: %w", err)
	}
	return nil
}

func (m *MongoStorageProvider) Close() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := m.client.Disconnect(ctx); err != nil {
		slog.Error("mongo: disconnect failed", "error", err)
	}
}