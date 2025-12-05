package main

import (
	"net/http"

	"github.com/Fry-Fr/chirpy/config"
)

func main() {
	cfg := &config.ApiConfig{}

	handlerFileServ := cfg.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app"))))
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", cfg.HandleHealthzStatus)
	mux.HandleFunc("GET /api/metrics", cfg.HandleMetricsLoad)
	mux.HandleFunc("POST /api/reset", cfg.HandleMetricsReset)
	mux.Handle("/app/", handlerFileServ)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
