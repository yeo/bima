package dto

import (
	"database/sql"
)

var dbConn *sql.DB

func SetDB(db *sql.DB) {
	dbConn = db
}
