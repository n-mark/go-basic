package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"example.com/go-basic/packages/internal/config"
	"example.com/go-basic/packages/internal/grpc"
	"example.com/go-basic/packages/internal/migrations"
	"example.com/go-basic/packages/internal/repository"
	"example.com/go-basic/packages/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type CommonServer struct {
	webServer     *Server
	consoleServer *ConsoleServer
	grpc          *grpc.GrpcServer
	repo          *repository.Repo
	entityService *service.EntityService
}

func NewCommonServer() *CommonServer {
	gameService := service.NewGameServiceNew()

	sp, err := initPostgresStorageProvider(context.Background())
	if err != nil {
		panic(fmt.Errorf("storage init: %w", err))
	}

	repo := repository.New(sp)
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
	}
}

func (s *CommonServer) Run() {
	defer s.repo.CloseStorage()
	go s.consoleServer.RunConsoleListener()
	go grpc.RunGrpcListenerInParallel(s.grpc)
	s.webServer.RunServer()
}

func initPostgresStorageProvider(ctx context.Context) (repository.StorageProvider, error) {
	pgCfg := config.GetPGConfig()
	if err := applyDbMigrations(ctx, pgCfg); err != nil {
		slog.Error("failed to apply db migrations", "err", err)
		os.Exit(1)
	}
	pool, err := repository.ConnectPG(ctx, pgCfg.DSN())
	if err != nil {
		slog.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	pgStore := repository.NewPgStore(pool)

	return pgStore, nil
}

// applyDbMigrations прогоняет все SQL-миграции из встроенного в бинарник embed.FS
// через goose при старте сервиса.
func applyDbMigrations(ctx context.Context, pgCfg config.PGConfig) error {
	db, err := sql.Open("pgx", pgCfg.DSN())
	if err != nil {
		return fmt.Errorf("open postgres: %w", err)
	}
	defer db.Close()

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("ping postgres: %w", err)
	}
	log.Println("✅ Подключение к PostgreSQL установлено")

	goose.SetBaseFS(migrations.FS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	log.Println("🔄 Запуск миграций...")
	if err := goose.UpContext(ctx, db, "."); err != nil {
		return fmt.Errorf("goose up: %w", err)
	}

	log.Println("✅ Все миграции успешно применены")
	return nil
}
