package db

import (
	"database/sql"
	"log"
)

const (
	schema_version = 1
)

func setupMigration(db *sql.DB) error {
	sqlStmt := "CREATE TABLE migration (version string);"
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	return err
}

func runMigration(db *sql.DB) error {
	sqls := []string{
		`
		CREATE TABLE config (name text not null primary key, value text, scope text);
		CREATE TABLE secret (id text not null primary key, name text, url text, token blob, version integer DEFAULT 1);
		`,
	}

	for _, sqlStmt := range sqls {
		_, err := db.Exec(sqlStmt)
		if err != nil {
			return err
		}
	}

	return nil
}
