package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/LeonMdS/chirpy-server/internal/auth"
	"github.com/LeonMdS/chirpy-server/internal/database"
	"github.com/google/uuid"
)

type user struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func removeSensitiveUserData(u database.User) user {
	outputUser := user{
		ID:          u.ID,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
		Email:       u.Email,
		IsChirpyRed: u.IsChirpyRed,
	}

	return outputUser
}

func (cfg *APIConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	reqJSON := req{}
	if err := decoder.Decode(&reqJSON); err != nil {
		respondWithError(w, http.StatusBadRequest, "Something went wrong", err)
		return
	}

	hashedPassword, err := auth.HashPassword(reqJSON.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
	}

	creationParams := database.CreateUserParams{
		Email:          reqJSON.Email,
		HashedPassword: hashedPassword,
	}

	createdUser, err := cfg.db.CreateUser(r.Context(), creationParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	outputUser := removeSensitiveUserData(createdUser)
	respondWithJSON(w, http.StatusCreated, outputUser)
}

func (cfg *APIConfig) resetUsersHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Not allowed", nil)
		return
	}
	cfg.fileserverHits.Store(0)
	err := cfg.db.DeleteAllUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type updateEmailPasswordRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (cfg *APIConfig) updateEmailPasswordHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get authorization token", err)
		return
	}

	foundID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token", err)
		return
	}

	req := updateEmailPasswordRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not parse request", err)
		return
	}

	newPwdHashed, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not hash password", err)
		return
	}
	updatedUser, err := cfg.db.UpdateUserEmailPassword(r.Context(), database.UpdateUserEmailPasswordParams{
		ID:             foundID,
		Email:          req.Email,
		HashedPassword: newPwdHashed,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not update database", err)
		return
	}

	outputUser := removeSensitiveUserData(updatedUser)
	respondWithJSON(w, http.StatusOK, outputUser)
}
