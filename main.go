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
	state := &config.State{
		ApiConfig: &config.ApiConfig{},
	}
	if err := state.ConnectDatabase(); err != nil {
		log.Printf("Error ConnectDatabase: %v\n", err)
	}

	handlerFileServ := state.ApiConfig.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app"))))
	mux := http.NewServeMux()
	mux.HandleFunc("POST /admin/reset", state.ApiConfig.HandleMetricsReset)
	mux.HandleFunc("GET /admin/metrics", state.ApiConfig.HandleMetricsLoad)
	mux.HandleFunc("GET /api/healthz", config.HandleHealthzStatus)

	mux.HandleFunc("POST /api/validate_chirp", config.HandleValidateChirp)
	mux.Handle("/app/", handlerFileServ)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
