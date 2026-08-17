package api

import (
	"net/http"
)

func NewRouter(cfg *APIConfig) *http.ServeMux {
	const filepathRoot = "."

	mux := http.NewServeMux()

	fileServerHandler := http.FileServer(http.Dir(filepathRoot))
	mux.Handle("/app/", http.StripPrefix("/app/", cfg.middlewareMetricsInc(fileServerHandler)))

	mux.HandleFunc("GET /admin/metrics", cfg.metricsHandler)
	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getAllChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpHandler)

	mux.HandleFunc("POST /admin/reset", cfg.resetUsersHandler)
	mux.HandleFunc("POST /api/users", cfg.createUserHandler)
	mux.HandleFunc("POST /api/chirps", cfg.addChirpHandler)
	mux.HandleFunc("POST /api/login", cfg.loginHandler)
	mux.HandleFunc("POST /api/refresh", cfg.refreshHandler)
	mux.HandleFunc("POST /api/revoke", cfg.revokeHandler)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.userUpgradeHandler)

	mux.HandleFunc("PUT /api/users", cfg.updateEmailPasswordHandler)

	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.deleteChirpHandler)

	return mux
}
