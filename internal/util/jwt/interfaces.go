package jwt

// TokenManager интерфейс для работы с JWT токенами
type TokenManager interface {
	GenerateTokens(userID, email string) (*TokenPair, error)
	ValidateAccessToken(token string) (*Claims, error)
	ValidateRefreshToken(token string) (*Claims, error)
	RefreshTokens(refreshToken string) (*TokenPair, error)
}

// Проверяем, что Manager реализует интерфейс
var _ TokenManager = (*Manager)(nil)
