package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LeonMdS/chirpy-server/internal/auth"
	"github.com/LeonMdS/chirpy-server/internal/database"
	"github.com/google/uuid"
)

type loginParams struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type loginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	IsRed        bool      `json:"is_chirpy_red"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
}

const (
	accessTokenExpiry  = time.Hour
	refreshTokenExpiry = time.Hour * 24 * 60
)

func (cfg *APIConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	requestParams := loginParams{}
	if err := json.NewDecoder(r.Body).Decode(&requestParams); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding login request", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), requestParams.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	authorized, err := auth.CheckPasswordHash(requestParams.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if !authorized {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	newToken, err := auth.MakeJWT(user.ID, cfg.secretKey, accessTokenExpiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating token", err)
		return
	}
	newRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	_, err = cfg.db.AddRefreshToken(r.Context(), database.AddRefreshTokenParams{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(refreshTokenExpiry),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error saving refresh token", err)
		return
	}

	response := loginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		IsRed:        user.IsChirpyRed,
		Token:        newToken,
		RefreshToken: newRefreshToken,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *APIConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find refresh token in header", err)
		return
	}

	foundToken, err := cfg.db.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not find token in database", err)
		return
	}

	newAccessToken, err := auth.MakeJWT(foundToken.UserID, cfg.secretKey, accessTokenExpiry)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating new access token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": newAccessToken})
}

func (cfg *APIConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not find refresh token in header", err)
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error revoking refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
