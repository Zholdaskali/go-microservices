package app

import (
	"auth-service/internal/config"
	"auth-service/internal/database"
	"auth-service/internal/handler"
	"auth-service/internal/handler/grpchandler"
	"auth-service/internal/logger"
	"auth-service/internal/repository"
	"auth-service/internal/repository/postgres"
	"auth-service/internal/service"
	"auth-service/internal/util/jwt"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

// Dependencies контейнер зависимостей
type Dependencies struct {
	DB          *sqlx.DB
	JWTManager  jwt.TokenManager
	UserRepo    repository.UserRepository
	AuthService service.AuthService
	AuthHandler handler.AuthHandler
}

// NewDependencies создает все зависимости в правильном порядке
func NewDependencies(cfg *config.Config, log logger.Logger) (*Dependencies, error) {
	deps := &Dependencies{}

	// 1. База данных (с миграциями)
	if err := deps.initDatabase(cfg, log); err != nil {
		return nil, err
	}

	// 2. JWT менеджер
	deps.initJWTManager(cfg, log)

	// 3. Репозитории
	deps.initRepositories(log)

	// 4. Сервисы
	deps.initServices(log)

	// 5. Обработчики
	deps.initHandlers(log)

	return deps, nil
}

// initDatabase инициализирует БД и применяет миграции
func (d *Dependencies) initDatabase(cfg *config.Config, log logger.Logger) error {
	log.Info("Connecting to database...")

	// Создаем подключение
	db, err := sqlx.Connect("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Error("Fatal error connect ")
		return err
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверяем подключение
	if err := db.Ping(); err != nil {
		return err
	}

	log.Info("Database connection established")

	// Применяем миграции ЧЕРЕЗ СУЩЕСТВУЮЩЕЕ ПОДКЛЮЧЕНИЕ
	log.Info("Applying database migrations...")
	if err := database.MigrateUp(cfg.DatabaseURL); err != nil {
		fmt.Println("Ошибка миграции: ", err)
		return err
	}
	log.Info("Database migrations applied successfully")

	d.DB = db
	return nil
}

// initJWTManager инициализирует JWT менеджер
func (d *Dependencies) initJWTManager(cfg *config.Config, log logger.Logger) {
	jwtCfg := jwt.Config{
		AccessTokenSecret:  cfg.JWTSecret,
		RefreshTokenSecret: cfg.JWTRefreshSecret,
		AccessTokenExpiry:  cfg.AccessTokenExpiry,
		RefreshTokenExpiry: cfg.RefreshTokenExpiry,
	}

	d.JWTManager = jwt.NewManager(jwtCfg)

	log.Info("JWT manager configured",
		logger.F("access_expiry", jwtCfg.AccessTokenExpiry),
		logger.F("refresh_expiry", jwtCfg.RefreshTokenExpiry),
	)
}

// initRepositories инициализирует репозитории
func (d *Dependencies) initRepositories(log logger.Logger) {
	d.UserRepo = postgres.NewUserRepository(d.DB, log)
	log.Info("User repository initialized")
}

// initServices инициализирует сервисы
func (d *Dependencies) initServices(log logger.Logger) {
	d.AuthService = service.NewAuthService(d.UserRepo, d.JWTManager, log)
	log.Info("Auth service initialized")
}

// initHandlers инициализирует обработчики
func (d *Dependencies) initHandlers(log logger.Logger) {
	d.AuthHandler = grpchandler.NewAuthHandler(d.AuthService, log)
	log.Info("Auth handler initialized")
}

// Close закрывает все зависимости
func (d *Dependencies) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
}
