package api

import (
	"encoding/json"
	"net/http"

	"github.com/LeonMdS/chirpy-server/internal/auth"
	"github.com/google/uuid"
)

type upgradeRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserID uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *APIConfig) userUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	APIKey, err := auth.GetAPIKeyFromHeader(r.Header)
	if err != nil || APIKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	req := upgradeRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding request", err)
		return
	}

	if req.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeToRed(r.Context(), req.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Error fetching user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
