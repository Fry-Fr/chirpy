package main

import (
	"net/http"

	"github.com/Fry-Fr/chirpy/config"
)

func main() {
	state := &config.ApiConfig{}

	handlerFileServ := state.MiddlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("./app"))))
	mux := http.NewServeMux()
	mux.HandleFunc("POST /admin/reset", state.HandleMetricsReset)
	mux.HandleFunc("GET /admin/metrics", state.HandleMetricsLoad)
	mux.HandleFunc("GET /api/healthz", config.HandleHealthzStatus)

	mux.HandleFunc("POST /api/validate_chirp", config.HandleValidateChirp)
	mux.Handle("/app/", handlerFileServ)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	server.ListenAndServe()
}
