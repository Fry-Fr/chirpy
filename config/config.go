package config

import (
	"database/sql"
	"os"
	"sync/atomic"

	"github.com/Fry-Fr/chirpy/internal/database"
)

type ApiConfig struct {
	FileserverHits atomic.Int32
	DB             *database.Queries
}

func (cfg *ApiConfig) ConnectDatabase() error {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	dbQueries := database.New(db)
	cfg.DB = dbQueries
	return nil
}
