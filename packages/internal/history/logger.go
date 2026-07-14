package history

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

type Action string

const (
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

type ChangeRecord struct {
	EntityType string    `json:"entity_type"`
	EntityID   int       `json:"entity_id"`
	Action     Action    `json:"action"`
	OldValue   any       `json:"old_value,omitempty"`
	NewValue   any       `json:"new_value,omitempty"`
	Timestamp  time.Time `json:"timestamp"`
}

type Logger interface {
	Log(ctx context.Context, record ChangeRecord) error
	GetHistory(ctx context.Context, entityType string, entityID int) ([]ChangeRecord, error)
	GetAllHistory(ctx context.Context, entityType string) ([]ChangeRecord, error)
	Close() error
}

type RedisLogger struct {
	client       *redis.Client
	entityTTL    time.Duration
	globalTTL    time.Duration
	historyLimit int64
}

type Options struct {
	Addr         string
	Password     string
	DB           int
	EntityTTL    time.Duration
	GlobalTTL    time.Duration
	HistoryLimit int64
}


func NewRedisLogger(ctx context.Context, opts Options) (*RedisLogger, error) {
	if opts.Addr == "" {
		opts.Addr = "localhost:6379"
	}
	if opts.EntityTTL <= 0 {
		opts.EntityTTL = 24 * time.Hour
	}
	if opts.GlobalTTL <= 0 {
		opts.GlobalTTL = 7 * opts.EntityTTL
	}

	client := redis.NewClient(&redis.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &RedisLogger{
		client:       client,
		entityTTL:    opts.EntityTTL,
		globalTTL:    opts.GlobalTTL,
		historyLimit: opts.HistoryLimit,
	}, nil
}

func (l *RedisLogger) Client() *redis.Client {
	return l.client
}

func (l *RedisLogger) Log(ctx context.Context, record ChangeRecord) error {
	if record.Timestamp.IsZero() {
		record.Timestamp = time.Now()
	}

	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("history: marshal record: %w", err)
	}

	score := float64(record.Timestamp.UnixNano())
	globalKey := globalKey(record.EntityType)
	entityKey := entityKey(record.EntityType, record.EntityID)

	pipe := l.client.Pipeline()
	pipe.ZAdd(ctx, globalKey, redis.Z{Score: score, Member: data})
	pipe.ZAdd(ctx, entityKey, redis.Z{Score: score, Member: data})
	pipe.Expire(ctx, entityKey, l.entityTTL)
	pipe.Expire(ctx, globalKey, l.globalTTL)

	if l.historyLimit > 0 {
		pipe.ZRemRangeByRank(ctx, globalKey, 0, -(l.historyLimit + 1))
		pipe.ZRemRangeByRank(ctx, entityKey, 0, -(l.historyLimit + 1))
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("history: redis pipeline: %w", err)
	}

	slog.Debug("history: change recorded",
		"entity", record.EntityType,
		"id", record.EntityID,
		"action", record.Action,
	)
	return nil
}

func (l *RedisLogger) GetHistory(ctx context.Context, entityType string, entityID int) ([]ChangeRecord, error) {
	return l.readZSet(ctx, entityKey(entityType, entityID))
}

func (l *RedisLogger) GetAllHistory(ctx context.Context, entityType string) ([]ChangeRecord, error) {
	return l.readZSet(ctx, globalKey(entityType))
}

func (l *RedisLogger) readZSet(ctx context.Context, key string) ([]ChangeRecord, error) {
	raw, err := l.client.ZRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("history: zrange %q: %w", key, err)
	}
	records := make([]ChangeRecord, 0, len(raw))
	for _, m := range raw {
		var r ChangeRecord
		if err := json.Unmarshal([]byte(m), &r); err != nil {
			slog.Warn("history: can't deserialize record", "error", err)
			continue
		}
		records = append(records, r)
	}
	return records, nil
}

func (l *RedisLogger) Close() error {
	return l.client.Close()
}

func globalKey(entityType string) string {
	return "history:" + entityType
}

func entityKey(entityType string, id int) string {
	return fmt.Sprintf("history:%s:%d", entityType, id)
}