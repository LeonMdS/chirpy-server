// Package auth responsible for user authentication
package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
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

func GetBearerToken(headers http.Header) (string, error) {
	authString := headers.Get("Authorization")
	if authString == "" {
		return "", errors.New("authorization header not found")
	}

	wordsList := strings.Fields(authString)

	if len(wordsList) != 2 || strings.ToLower(wordsList[0]) != "bearer" {
		return "", errors.New("invalid authorization header format")
	}

	return wordsList[1], nil
}

func GetUIDFromHeaderToken(header http.Header, secretKey string) (uuid.UUID, error) {
	token, err := GetBearerToken(header)
	if err != nil {
		return uuid.Nil, err
	}

	foundID, err := ValidateJWT(token, secretKey)
	if err != nil {
		return uuid.Nil, err
	}

	return foundID, nil
}

func GetAPIKeyFromHeader(header http.Header) (string, error) {
	authString := header.Get("Authorization")
	if authString == "" {
		return "", errors.New("authorization header not found")
	}

	wordsList := strings.Fields(authString)

	if len(wordsList) != 2 || wordsList[0] != "ApiKey" {
		return "", errors.New("invalid authorization header format")
	}

	return wordsList[1], nil
}
