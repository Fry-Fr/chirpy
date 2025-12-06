package main

import (
	"log"
	"net/http"

	"github.com/Fry-Fr/chirpy/config"
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
	mux.HandleFunc("POST /admin/reset", cfg.HandleMetricsReset)
	mux.HandleFunc("GET /admin/metrics", cfg.HandleMetricsLoad)
	mux.HandleFunc("GET /api/healthz", config.HandleHealthzStatus)

	mux.HandleFunc("POST /api/login", cfg.HandleLogin)
	mux.HandleFunc("GET /api/chirps/{chirpId}", cfg.HandleGetChirp)
	mux.HandleFunc("GET /api/chirps", cfg.HandleGetChirps)
	mux.HandleFunc("POST /api/chirps", cfg.HandleCreateChirp)
	mux.HandleFunc("POST /api/users", cfg.HandleCreateUser)
	mux.Handle("/app/", handlerFileServ)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
