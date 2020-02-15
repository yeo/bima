package db

import (
	"database/sql"
	"fmt"
	"os"
	"os/user"

	"github.com/rs/zerolog/log"

	_ "github.com/mattn/go-sqlite3"
)

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func Setup() (*sql.DB, error) {
	user, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	dbDir := user.HomeDir + "/.bima"
	os.MkdirAll(dbDir, os.ModePerm)

	dbPath := dbDir + "/bima.db"
	log.Info().Str("dbpath", dbPath).Msg("Load db file")

	needSetup := !fileExists(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot sqlite db file. Check file permission")
	}

	if needSetup {
		err := setupMigration(db)
		if err != nil {
			log.Fatal().Err(err).Msg("Cannot migrate database")
		}
	}

	runMigrations(db)

	return db, nil
}
