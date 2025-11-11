package postgres

import (
	"auth-service/internal/domain"
	"auth-service/internal/logger"
	"auth-service/internal/repository"
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db  *sqlx.DB
	log logger.Logger
}

func NewUserRepository(db *sqlx.DB, log logger.Logger) repository.UserRepository {
	return &userRepository{
		db:  db,
		log: log.With(logger.F("layer", "repository"), logger.F("component", "user_repository")),
	}
}

// TODO реализация
//* Реализован
// Create(ctx context.Context, user *domain.User) error
// GetByID(ctx context.Context, id string) (*domain.User, error)
// GetByEmail(ctx context.Context, email string) (*domain.User, error)
// Update(ctx context.Context, user *domain.User) error
// Delete(ctx context.Context, id string) error

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {

	r.log.Debug("creating user",
		logger.F("user_id", user.ID),
		logger.F("email", user.Email),
	)

	query := `
		INSERT INTO t_users (id, username, email, password_hash, create_at, update_at) 
			VALUES ($1, $2, $3, $4, $5, $6)`
	fmt.Print(query)

	//* генерация нового айдишника для пользотеля
	user.ID = uuid.New()

	//* НАСТРОИВАЕМ ДАТУ
	now := time.Now()
	user.Create_at = now
	user.Update_at = now

	_, err := r.db.ExecContext(ctx, query,
		user.ID,
		user.UserName,
		user.Email,
		user.PasswordHash,
		user.Create_at,
		user.Update_at,
	)

	if err != nil {
		if isUniqueConstraintViolation(err) {
			return repository.ErrUserExists
		}
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func isUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "23505")
}

// TODO реализация
// * Реализован
func (r *userRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {

	r.log.Debug("GetByID user",
		logger.F("user_id", id),
	)

	query := `
		SELECT id, username, email, password_hash, create_at, update_at
		FROM t_users
		WHERE id = $1
	`

	var user domain.User

	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

// TODO реализация
// * Реализован
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {

	r.log.Debug("creating user",
		logger.F("email", email),
	)

	query := `
		SELECT id, username, email, password_hash, create_at, update_at
		FROM t_users
		WHERE email = $1
	`

	var user domain.User

	if err := r.db.GetContext(ctx, &user, query, email); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

// TODO реализация
// * Реализован
func (r *userRepository) Update(ctx context.Context, user *domain.User) error {
	r.log.Debug("Update user",
		logger.F("user_id", user.ID),
		logger.F("email", user.Email),
	)

	query := `
		UPDATE t_users 
        SET 
            username = $1,
            email = $2,
            update_at = $3
        WHERE id = $4
	`

	_, err := r.db.ExecContext(ctx, query,
		user.UserName,
		user.Email,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// TODO реализация
// * Реализован
func (r *userRepository) Delete(ctx context.Context, id string) error {

	r.log.Debug("creating user",
		logger.F("user_id", id),
	)

	query := `
        DELETE FROM t_users WHERE id = $1
    `

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}

	return nil
}
