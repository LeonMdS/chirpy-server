// Package auth responsible for user authentication
package auth

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/alexedwards/argon2id"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	correctPassword, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return correctPassword, nil
}

func MakeRefreshToken() (string, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", err
	}

	tokenString := hex.EncodeToString(tokenBytes)

	return tokenString, nil
}
