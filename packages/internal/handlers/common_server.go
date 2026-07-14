package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"example.com/go-basic/packages/internal/grpc"
	"example.com/go-basic/packages/internal/history"
	"example.com/go-basic/packages/internal/repository"
	"example.com/go-basic/packages/internal/service"
)

type CommonServer struct {
	webServer     *Server
	consoleServer *ConsoleServer
	grpc          *grpc.GrpcServer
	repo          *repository.Repo
	entityService *service.EntityService
	historyLogger history.Logger
}

const (
	defaultMongoURI       = "mongodb://localhost:27017"
	defaultMongoDB        = "go_basic"
	defaultRedisAddr      = "localhost:6379"
	defaultRedisEntityTTL = 24 * time.Hour
	defaultRedisGlobalTTL = 7 * 24 * time.Hour
	defaultRedisLimit     = 1000
)

func NewCommonServer() *CommonServer {
	gameService := service.NewGameServiceNew()

	sp, hist, err := initStorage(context.Background())
	if err != nil {
		panic(fmt.Errorf("storage init: %w", err))
	}

	repo := repository.NewWithHistory(sp, hist)
	entityService := service.NewEntityService(repo)
	ws := InitServer(gameService, entityService)
	cs := NewConsoleServer(gameService)
	g := grpc.NewGrpcServer(entityService)

	return &CommonServer{
		webServer:     ws,
		consoleServer: cs,
		grpc:          g,
		repo:          repo,
		entityService: entityService,
		historyLogger: hist,
	}
}

func (s *CommonServer) Run() {
	defer s.repo.CloseStorage()
	if s.historyLogger != nil {
		defer func() {
			if err := s.historyLogger.Close(); err != nil {
				slog.Error("history logger close failed", "error", err)
			}
		}()
	}
	go s.consoleServer.RunConsoleListener()
	go grpc.RunGrpcListenerInParallel(s.grpc)
	s.webServer.RunServer()
}

func initStorage(ctx context.Context) (repository.StorageProvider, history.Logger, error) {
	storageType := os.Getenv("STORAGE_TYPE")
	if storageType == "" {
		storageType = "mongo"
	}

	switch storageType {
	case "local":
		slog.Info("using local storage")
		return repository.NewLocalStorageProvider(), nil, nil

	case "mongo":
		return initMongoAndRedis(ctx)

	default:
		return nil, nil, fmt.Errorf("unknown STORAGE_TYPE=%q (expected local|mongo)", storageType)
	}
}

func initMongoAndRedis(ctx context.Context) (repository.StorageProvider, history.Logger, error) {
	mongoURI := getenv("MONGO_URI", defaultMongoURI)
	mongoDB := getenv("MONGO_DB", defaultMongoDB)

	mongoProvider, err := repository.NewMongoStorageProvider(ctx, repository.MongoOptions{
		URI:      mongoURI,
		Database: mongoDB,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("mongo: %w", err)
	}
	slog.Info("mongo: connected", "uri", mongoURI, "db", mongoDB)

	redisLogger, err := history.NewRedisLogger(ctx, history.Options{
		Addr:         getenv("REDIS_ADDR", defaultRedisAddr),
		Password:     os.Getenv("REDIS_PASSWORD"),
		DB:           0,
		EntityTTL:    parseDuration("REDIS_ENTITY_TTL", defaultRedisEntityTTL),
		GlobalTTL:    parseDuration("REDIS_GLOBAL_TTL", defaultRedisGlobalTTL),
		HistoryLimit: defaultRedisLimit,
	})
	if err != nil {
		mongoProvider.Close()
		return nil, nil, fmt.Errorf("redis: %w", err)
	}
	slog.Info("redis: connected", "addr", getenv("REDIS_ADDR", defaultRedisAddr))

	return mongoProvider, redisLogger, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func parseDuration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		slog.Warn("can't parse duration, use default value",
			"key", key, "value", v, "fallback", fallback, "error", err)
		return fallback
	}
	return d
}
