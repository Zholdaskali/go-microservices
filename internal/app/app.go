package app

import (
	"auth-service/internal/config"
	"auth-service/internal/logger"
	grpcserver "auth-service/internal/server/grpc" // единый алиас
	"os"
	"os/signal"
	"syscall"
)

// App - основная структура приложения
type App struct {
	config     *config.Config
	logger     logger.Logger
	grpcServer *grpcserver.Server // используем алиас
}

// New создает новое приложение
func New(mode string) (*App, error) {
	app := &App{}

	// Инициализация в правильном порядке
	if err := app.initConfig(mode); err != nil {
		return nil, err
	}

	if err := app.initLogger(); err != nil {
		return nil, err
	}

	if err := app.initDependencies(); err != nil {
		return nil, err
	}

	return app, nil
}

// initConfig загружает конфигурацию
func (a *App) initConfig(mode string) error {
	switch mode {
	case "dev":
		a.config = config.LoadConfigDev()
	case "prod":
		a.config = config.LoadConfigProd()
	case "test":
		a.config = config.LoadConfigTest()
	default:
		a.config = config.LoadConfigDev()
	}
	return nil
}

// initLogger инициализирует логгер
func (a *App) initLogger() error {
	logger, err := logger.New(a.config.LogLevel)
	if err != nil {
		return err
	}
	a.logger = logger
	return nil
}

// initDependencies инициализирует все зависимости
func (a *App) initDependencies() error {
	// Создаем контейнер зависимостей
	deps, err := NewDependencies(a.config, a.logger)
	if err != nil {
		return err
	}

	// Создаем gRPC сервер и сохраняем в структуру App
	a.grpcServer = grpcserver.NewServer(deps.AuthHandler, a.logger)

	return nil
}

// Run запускает приложение
func (a *App) Run() error {
	a.logger.Info("Starting auth service...")

	// Запускаем сервер
	if err := a.grpcServer.Start(a.config.GRPCPort); err != nil {
		return err
	}
	// reflection.Register(a.grpcServer)
	a.waitForShutdown()
	return nil
}

// waitForShutdown ожидает сигнал завершения
func (a *App) waitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	a.logger.Info("Received signal, shutting down...", logger.F("signal", sig))

	a.grpcServer.Stop()
	a.logger.Info("Service stopped gracefully")
}
