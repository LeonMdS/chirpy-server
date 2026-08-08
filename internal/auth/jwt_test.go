package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestJWTFunctionality(t *testing.T) {
	secret := "secret"
	userID := uuid.New()

	tokenString, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	resultingID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatal(err)
	}

	if resultingID != userID {
		t.Fatal("userID does not match")
	}
}

func TestExpiredToken(t *testing.T) {
	type testCase struct {
		userID      uuid.UUID
		testName    string
		secret      string
		expiresIn   time.Duration
		expectedErr error
	}

	testList := []testCase{}
	defaultSecret := "secret"
	defaultExpiresIn := time.Hour

	testList = append(testList, testCase{
		userID:      uuid.New(),
		testName:    "Expired token",
		secret:      defaultSecret,
		expiresIn:   -time.Hour,
		expectedErr: jwt.ErrTokenExpired,
	})

	testList = append(testList, testCase{
		userID:      uuid.New(),
		testName:    "Wrong secret",
		secret:      "wrong",
		expiresIn:   defaultExpiresIn,
		expectedErr: jwt.ErrTokenSignatureInvalid,
	})

	for _, test := range testList {
		t.Run(test.testName, func(t *testing.T) {
			tokenString, err := MakeJWT(test.userID, test.secret, test.expiresIn)
			if err != nil {
				t.Fatal(err)
			}

			resultingID, err := ValidateJWT(tokenString, defaultSecret)
			if !errors.Is(err, test.expectedErr) {
				t.Fatalf("unexpected error: %v", err)
			}

			if test.expectedErr != nil {
				return
			}

			if resultingID != test.userID {
				t.Fatal("userID does not match")
			}
		})
	}
}
