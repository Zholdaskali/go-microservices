// cmd/server/main.go
package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/handler/grpchandler"
	"auth-service/internal/logger"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	"auth-service/internal/util/jwt"
	pb "auth-service/pkg/api/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	//* Загружаем конфигурацию
	cfg := сonfig(mustMode())

	//* 1. Инициализация логгера
	appLogger, err := logger.New(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	appLogger.Info("Starting auth service...")

	//* 2. Запуск миграции
	databaseURL := cfg.DatabaseURL
	appLogger.Info("Applying database migrations...")
	if err := database.MigrateUp(databaseURL); err != nil {
		appLogger.Fatal("Failed to apply migrations", logger.F("error", err))
	}
	appLogger.Info("Database migrations applied successfully")

	//* 3. Подключение к PostgreSQL
	db, err := connectDatabase(appLogger, cfg)
	if err != nil {
		appLogger.Fatal("Failed to connect to database", logger.F("error", err))
	}
	defer db.Close()
	appLogger.Info("Database connection established")

	//* 4. Инициализация JWT менеджера
	jwtManager := setupJWTManager(appLogger, cfg)

	//* 5. Создание репозиториев
	userRepo := postgres.NewUserRepository(db, appLogger)

	//* 6. Создание сервисов
	authService := service.NewAuthService(userRepo, jwtManager, appLogger)

	//* 7. Создание gRPC handlers
	authHandler := grpchandler.NewAuthHandler(authService, appLogger)

	//* 8. Запуск gRPC сервера
	grpcServer := startGRPCServer(authHandler, appLogger)
	reflection.Register(grpcServer)

	//* 9. Ожидание сигнала завершения
	waitForShutdown(grpcServer, db, appLogger)

}

// * connectDatabase подключается к PostgreSQL
func connectDatabase(log logger.Logger, cfg *config.Config) (*sqlx.DB, error) {

	log.Info("Информация попытка соединения с БД")

	connStr := cfg.DatabaseURL

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, err
	}

	//* Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	//* Проверяем подключение
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}

func сonfig(typeCfg string) *config.Config {
	// ! Реализовать срочно
	// TODO РЕАЛИЗОВАТЬ ПРИМЕР ДЛЯ ПРОД
	//* Определение по типу конфигурации
	var cfg *config.Config
	switch typeCfg {
	case "dev":
		cfg = config.LoadConfigDev()
	}
	return cfg

}

// setupJWTManager настраивает JWT менеджер
func setupJWTManager(log logger.Logger, cfg *config.Config) *jwt.Manager {

	accessSecret := cfg.JWTSecret
	refreshSecret := cfg.JWTRefreshSecret

	config := jwt.Config{
		AccessTokenSecret:  accessSecret,
		RefreshTokenSecret: refreshSecret,
		AccessTokenExpiry:  15 * time.Minute,
		RefreshTokenExpiry: 7 * 24 * time.Hour,
	}

	log.Info("JWT manager configured",
		logger.F("access_expiry", config.AccessTokenExpiry),
		logger.F("refresh_expiry", config.RefreshTokenExpiry),
	)

	return jwt.NewManager(config)

}

// startGRPCServer запускает gRPC сервер
func startGRPCServer(authHandler *grpchandler.AuthHandler, log logger.Logger) *grpc.Server {

	grpcServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(grpcServer, authHandler)

	port := getEnv("GRPC_PORT", "50051")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("Failed to listen", logger.F("error", err), logger.F("port", port))
	}

	go func() {
		log.Info("gRPC server started", logger.F("port", port))
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal("gRPC server failed", logger.F("error", err))
		}
	}()

	return grpcServer

}

// waitForShutdown ожидает сигнал завершения
func waitForShutdown(grpcServer *grpc.Server, db *sqlx.DB, log logger.Logger) {

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Info("Received signal, shutting down...", logger.F("signal", sig))

	grpcServer.GracefulStop()
	db.Close()

	log.Info("Service stopped gracefully")

}

// getEnv получает переменную окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue

}

func mustMode() string {

	mode := flag.String(
		"app-mode",
		"",
		"application launch mode for selecting settings in env",
	)

	flag.Parse()

	if *mode == "" { // ← ПРАВИЛЬНАЯ ПРОВЕРКА!
		log.Fatal("launch mode is not specified")
	}

	return *mode

}
