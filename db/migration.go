package db

import (
	"database/sql"
	"fmt"
	"log"
)

const (
	schema_version = 1
)

type MigrationUnit struct {
	Name  string
	Query string
}

func setupMigration(db *sql.DB) error {
	sqlStmt := "CREATE TABLE migration (version int, name string);"
	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}
	return err
}

func checkMirationRan(db *sql.DB, version int) (bool, error) {
	stmt, err := db.Prepare("select version from migration where version = ?")
	if err != nil {
		return false, err
	}

	defer stmt.Close()
	var foundVersion int
	err = stmt.QueryRow(version).Scan(&foundVersion)

	switch err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return foundVersion == version, nil
	default:
		return false, err
	}
}

func runMigrations(db *sql.DB) error {
	sqls := []*MigrationUnit{
		&MigrationUnit{
			Name:  "create config table",
			Query: "CREATE TABLE config (name text not null primary key, value text, scope text);",
		},
		&MigrationUnit{
			Name: "create secret tables",
			Query: `
						CREATE TABLE secret (id text not null primary key, name text, url text, token blob, version integer DEFAULT 1, deleted_at INTEGER);
						`,
		},
	}

	for i, sqlStmt := range sqls {
		ran, err := checkMirationRan(db, i+1)
		log.Println("Migrtion", i, "status is", ran, "return error", err)
		if !ran && err == nil {
			log.Println("Run migration", sqlStmt.Query)
			_, err = db.Exec(sqlStmt.Query)
			if err != nil {
				return err
			}
			db.Exec(fmt.Sprintf("INSERT INTO migration (version, name) VALUES (%d, '%s')", i+1, sqlStmt.Name))
		}
	}

	return nil
}
