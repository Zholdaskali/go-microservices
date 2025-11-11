package bcrypt

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrHashFailed       = errors.New("failed to hash password")
	ErrInvalidHash      = errors.New("invalid hash format")
	ErrWrongPassword    = errors.New("wrong password")
	ErrPasswordTooShort = errors.New("password must be at least 6 characters")
)

var bcryptCost = bcrypt.DefaultCost

// TODO ИСПРАВИТЬ ОБРАБОТКУ ОШИБОК
// * Реализован
func Hash(password string) (string, error) {

	if password == "" {
		return "", ErrHashFailed
	}

	if len(password) < 6 {
		return "", ErrPasswordTooShort
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)

	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrInvalidHash, err)
	}
	return string(hash), nil
}

// TODO ЧЕК РЕАЛИЗОВАТЬ ПРОВЕРКУ ПАРОЛЯ
// * Реализован
func Check(password, hash string) (bool, error) {
	if password == "" || hash == "" {
		return false, ErrInvalidHash
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, ErrWrongPassword
		case errors.Is(err, bcrypt.ErrHashTooShort):
			return false, ErrInvalidHash
		default:
			return false, err
		}
	}

	return true, nil

}
