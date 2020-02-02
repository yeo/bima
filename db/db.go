package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/user"

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
	log.Println("Load db path", dbDir)

	os.MkdirAll(dbDir, os.ModePerm)

	dbPath := dbDir + "/bima.db"

	log.Println("Load db path", dbPath)

	needSetup := false
	if !fileExists(dbPath) {
		needSetup = true
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}

	if needSetup {
		sqlStmt := `
	create table secret (id integer not null primary key, name text, token text);
	`
		_, err := db.Exec(sqlStmt)
		if err != nil {
			log.Printf("%q: %s\n", err, sqlStmt)
		}
	}

	return db, nil
}
