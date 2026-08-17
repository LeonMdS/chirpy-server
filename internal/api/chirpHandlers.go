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
		respondWithError(w, http.StatusBadRequest, "Error decoding request when adding chirp", err)
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
	lookupID := r.URL.Query().Get("author_id")
	parsedID, err := uuid.Parse(lookupID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error parsing ID from string", err)
		return
	}

	var response []database.Chirp
	if lookupID == "" {
		response, err = cfg.db.GetAllChirps(r.Context())
	} else {
		response, err = cfg.db.GetChirpsByAuthor(r.Context(), parsedID)
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error getting chirps from database", err)
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

func (cfg *APIConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	foundID, err := auth.GetUIDFromHeaderToken(r.Header, cfg.secretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Could not get authorization token", err)
		return
	}

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

	if chirp.UserID != foundID {
		respondWithError(w, http.StatusForbidden, "You are not allowed to delete this chirp", nil)
		return
	}

	if err = cfg.db.DeleteChirp(r.Context(), chirpID); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting chirp", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
