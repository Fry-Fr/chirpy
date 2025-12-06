package handlers

import (
	"fmt"
	"net/http"

	"github.com/Fry-Fr/chirpy/internal/config"
)

func AdminMetrics(cfg *config.ApiConfig, w http.ResponseWriter, r *http.Request) {
	hits := cfg.FileserverHits.Load()
	display := fmt.Sprintf(`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, hits)

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(display))
}
