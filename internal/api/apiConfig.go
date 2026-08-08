package api

import (
	"sync/atomic"

	"github.com/LeonMdS/chirpy-server/internal/database"
)

type APIConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secretKey      string
}

func NewAPIConfig(db *database.Queries, platform, secretKey string) *APIConfig {
	return &APIConfig{
		fileserverHits: atomic.Int32{},
		db:             db,
		platform:       platform,
		secretKey:      secretKey,
	}
}
