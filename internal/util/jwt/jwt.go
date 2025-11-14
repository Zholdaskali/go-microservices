// pkg/jwt/jwt.go
package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// Claims - кастомные claims для нашего приложения
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// TokenPair - пара access и refresh токенов
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Config - конфигурация JWT
type Config struct {
	AccessTokenSecret  string        `json:"access_token_secret"`
	RefreshTokenSecret string        `json:"refresh_token_secret"`
	AccessTokenExpiry  time.Duration `json:"access_token_expiry"`  // например: 15 * time.Minute
	RefreshTokenExpiry time.Duration `json:"refresh_token_expiry"` // например: 7 * 24 * time.Hour
}

// Manager - менеджер JWT токенов
type Manager struct {
	config Config
}

// NewManager создает новый менеджер JWT
func NewManager(config Config) *Manager {
	return &Manager{
		config: config,
	}
}

// GenerateTokens создает пару access и refresh токенов
func (m *Manager) GenerateTokens(userID, email string) (*TokenPair, error) {
	// Генерация Access Token
	accessToken, err := m.generateAccessToken(userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Генерация Refresh Token
	refreshToken, err := m.generateRefreshToken(userID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateAccessToken создает access token
func (m *Manager) generateAccessToken(userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.AccessTokenSecret))
}

// generateRefreshToken создает refresh token
func (m *Manager) generateRefreshToken(userID, email string) (string, error) {
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.config.RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.config.RefreshTokenSecret))
}

// ValidateAccessToken проверяет access token
func (m *Manager) ValidateAccessToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.config.AccessTokenSecret)
}

// ValidateRefreshToken проверяет refresh token
func (m *Manager) ValidateRefreshToken(tokenString string) (*Claims, error) {
	return m.validateToken(tokenString, m.config.RefreshTokenSecret)
}

// validateToken общая функция валидации
func (m *Manager) validateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// RefreshTokens обновляет пару токенов
func (m *Manager) RefreshTokens(refreshToken string) (*TokenPair, error) {
	claims, err := m.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}

	return m.GenerateTokens(claims.UserID, claims.Email)
}
