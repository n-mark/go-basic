package handlers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"example.com/go-basic/packages/internal/config"
	"example.com/go-basic/packages/internal/grpc"
	"example.com/go-basic/packages/internal/repository"
	"example.com/go-basic/packages/internal/service"

	"github.com/microbus-io/sequel"
	"github.com/microbus-io/sequel/sequelpgx"
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
	applyDbMigrations(ctx, pgCfg)
	pool, err := repository.ConnectPG(ctx, pgCfg.DSN())
	if err != nil {
		slog.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	pgStore := repository.NewPgStore(pool)

	return pgStore, nil
}

func applyDbMigrations(ctx context.Context, pgCfg config.PGConfig) {
	// 2. Устанавливаем соединение с БД через драйвер PostgreSQL
    db, err := sequel.Open("pgx", pgCfg.DSN())
    if err != nil {
        log.Fatalf("❌ Не удалось подключиться к PostgreSQL: %v", err)
    }
    defer db.Close()

    // Проверяем соединение
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    if err := db.Ping(); err != nil {
        log.Fatalf("❌ PostgreSQL не отвечает: %v", err)
    }
    log.Println("✅ Подключение к PostgreSQL установлено")

    // 3. Выполняем миграции
    log.Println("🔄 Запуск миграций...")
    
    // Важно: указываем уникальный идентификатор для отслеживания версии
    // Обычно это имя сервиса + версия миграций
    err = db.Migrate("my-service@v1", migrationFS)
    if err != nil {
        log.Fatalf("❌ Ошибка выполнения миграций: %v", err)
    }

    // Проверяем, какие миграции были применены
    applied, err := getAppliedMigrations(db)
    if err != nil {
        log.Printf("⚠️ Не удалось получить список примененных миграций: %v", err)
    } else {
        log.Printf("✅ Применено миграций: %d", len(applied))
        for _, m := range applied {
            log.Printf("   - %s (применена: %s)", m.Name, m.AppliedAt.Format(time.RFC3339))
        }
    }

    log.Println("✅ Все миграции успешно применены")

    // 4. Запускаем основное приложение
    log.Println("🚀 Запуск основного приложения...")
    // ... ваш код здесь
	panic("unimplemented")
}

// Структура для хранения информации о примененных миграциях
type AppliedMigration struct {
    Name      string
    AppliedAt time.Time
}

// Функция для получения списка примененных миграций
func getAppliedMigrations(db *sequel.DB) ([]AppliedMigration, error) {
    // sequel создает таблицу sequel_migrations для отслеживания
    rows, err := db.Query(`
        SELECT name, applied_at 
        FROM sequel_migrations 
        ORDER BY id
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var migrations []AppliedMigration
    for rows.Next() {
        var m AppliedMigration
        if err := rows.Scan(&m.Name, &m.AppliedAt); err != nil {
            return nil, err
        }
        migrations = append(migrations, m)
    }
    
    return migrations, rows.Err()
}