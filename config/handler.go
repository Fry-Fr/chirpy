package config

import (
	"fmt"
	"net/http"
)

func (cfg *ApiConfig) HandleHealthzStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) HandleMetricsReset(w http.ResponseWriter, r *http.Request) {
	cfg.FileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *ApiConfig) HandleMetricsLoad(w http.ResponseWriter, r *http.Request) {
	hits := cfg.FileserverHits.Load()
	display := fmt.Sprintf("Hits: %v", hits)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(display))
}
