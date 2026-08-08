package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/LeonMdS/chirpy-server/internal/auth"
	"github.com/LeonMdS/chirpy-server/internal/database"
	"github.com/google/uuid"
)

type addChirpRequest struct {
	Body string `json:"body"`
}

func (cfg *APIConfig) addChirpHandler(w http.ResponseWriter, r *http.Request) {
	req := addChirpRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding request when adding chirp", err)
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}
	req.Body = chirpCleaner(req.Body)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Could not get authorization token", err)
		return
	}

	foundID, err := auth.ValidateJWT(token, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Error validating token", err)
		return
	}

	newChirpParams := database.AddChirpParams{
		Body:   req.Body,
		UserID: foundID,
	}
	newChirp, err := cfg.db.AddChirp(r.Context(), newChirpParams)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Error adding chirp to database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, newChirp)
}

func (cfg *APIConfig) getAllChirpsHandler(w http.ResponseWriter, r *http.Request) {
	response, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting all chirps from database", err)
		return
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *APIConfig) getChirpHandler(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error getting chirp", err)
		return
	}

	respondWithJSON(w, http.StatusOK, chirp)
}
