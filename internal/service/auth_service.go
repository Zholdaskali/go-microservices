package service

import (
	"auth-service/internal/domain"
	"auth-service/internal/logger"
	"auth-service/internal/repository"
	"auth-service/internal/util/bcrypt"
	"auth-service/internal/util/jwt"
	pb "auth-service/pkg/api/service"
	"context"
	"errors"
)

var (
	ErrBadRequest         = errors.New("bad request")
	ErrUserNotFound       = errors.New("user not found")
	ErrPasswordBad        = errors.New("basd")
	ErrBadToken           = errors.New("asdasd")
	ErrInvalidCredentials = errors.New("")
)

type AuthService struct {
	UserRepo   repository.UserRepository
	jwtManager *jwt.Manager
	log        logger.Logger
	pb.UnimplementedAuthServiceServer
}

func NewAuthService(UserRepo repository.UserRepository, jwtManager *jwt.Manager, log logger.Logger) *AuthService {
	return &AuthService{
		UserRepo:   UserRepo,
		jwtManager: jwtManager,
		log:        log.With(logger.F("layer", "repository"), logger.F("component", "user_repository")),
	}
}

// ID           uuid.UUID `json:"id" db:"id"`
// UserName     string    `json:"user_name" db:"username"`
// Email        string    `json:"email" db:"email"`
// PasswordHash string    `json:"password_hash" db:"password_hash"`
// Create_at    time.Time `json:"create_at" db:"create_at"`
// Update_at    time.Time `json:"update_at" db:"update_at"`

func (s *AuthService) Register(ctx context.Context, registerRequest *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	username := registerRequest.UserName
	email := registerRequest.Email
	password := registerRequest.Password

	if username == "" || email == "" || password == "" {
		return nil, ErrBadRequest
	}

	passwordHash, err := bcrypt.Hash(password)

	if err != nil {
		return nil, err
	}

	user := &domain.User{
		UserName:     username,
		Email:        email,
		PasswordHash: passwordHash,
	}

	s.UserRepo.Create(ctx, user)

	return &pb.RegisterResponse{
		UserId: user.ID.String(),
	}, nil
}

func (s *AuthService) Login(ctx context.Context, loginRequest *pb.LoginRequest) (*pb.LoginResponse, error) {
	email := loginRequest.Email
	pass := loginRequest.Password

	if email == "" || pass == "" {
		return nil, ErrBadRequest
	}

	user, err := s.UserRepo.GetByEmail(ctx, loginRequest.Email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	isValid, err := bcrypt.Check(loginRequest.Password, user.PasswordHash)
	if err != nil {
		return nil, ErrPasswordBad
	}

	if !isValid {
		return nil, ErrInvalidCredentials
	}

	tokenPair, err := s.jwtManager.GenerateTokens(user.ID.String(), user.Email)
	if err != nil {
		return nil, ErrBadToken
	}

	// Возвращаем response
	return &pb.LoginResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
	}, nil

}

// TODO
// Create(ctx context.Context, user *domain.User) error
// GetByID(ctx context.Context, id string) (*domain.User, error)
// GetByEmail(ctx context.Context, email string) (*domain.User, error)
// Update(ctx context.Context, user *domain.User) error
// Delete(ctx context.Context, id string) error
