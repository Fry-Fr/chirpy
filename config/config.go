package config

import (
	"database/sql"
	"os"
	"sync/atomic"

	"github.com/Fry-Fr/chirpy/internal/database"
)

type State struct {
	ApiConfig *ApiConfig
	DB        *database.Queries
}

type ApiConfig struct {
	FileserverHits atomic.Int32
}

func (s *State) ConnectDatabase() error {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return err
	}
	dbQueries := database.New(db)
	s.DB = dbQueries
	return nil
}
