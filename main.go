package main

import (
	"log"
	"net/http"

	"github.com/Fry-Fr/chirpy/internal/config"
	"github.com/Fry-Fr/chirpy/internal/handlers"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()
	cfg := &config.ApiConfig{}
	if err := cfg.ConnectDatabase(); err != nil {
		log.Printf("Error ConnectDatabase: %v\n", err)
	}

	handlerFileServ := cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app"))))
	mux := http.NewServeMux()

	mux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		handlers.AdminReset(cfg, w, r)
	})
	mux.HandleFunc("GET /admin/metrics", func(w http.ResponseWriter, r *http.Request) {
		handlers.AdminMetrics(cfg, w, r)
	})
	mux.HandleFunc("GET /api/healthz", handlers.HealthStatus)

	mux.HandleFunc("POST /api/refresh", func(w http.ResponseWriter, r *http.Request) {
		handlers.RefreshToken(cfg, w, r)
	})
	mux.HandleFunc("POST /api/revoke", func(w http.ResponseWriter, r *http.Request) {
		handlers.RevokeToken(cfg, w, r)
	})
	mux.HandleFunc("POST /api/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.LoginUser(cfg, w, r)
	})
	mux.HandleFunc("GET /api/chirps/{chirpId}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetChirp(cfg, w, r)
	})
	mux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetChirps(cfg, w, r)
	})
	mux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateChirp(cfg, w, r)
	})
	mux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.CreateUser(cfg, w, r)
	})
	mux.HandleFunc("PUT /api/users", func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateUser(cfg, w, r)
	})
	mux.Handle("/app/", handlerFileServ)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
