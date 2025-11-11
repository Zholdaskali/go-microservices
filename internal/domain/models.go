package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	UserName     string    `json:"user_name" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"password_hash" db:"password_hash"`
	Create_at    time.Time `json:"create_at" db:"create_at"`
	Update_at    time.Time `json:"update_at" db:"update_at"`
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id" db:"id"`
	AuthId    uuid.UUID `json:"auth_id" db:"auth_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreateAt  time.Time `json:"create_at" db:"create_at"`
}
