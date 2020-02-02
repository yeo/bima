package dto

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Token struct {
	ID    int
	Name  string
	Token string
}

var dbConn *sql.DB

func SetDB(db *sql.DB) {
	dbConn = db
}

func LoadTokens() ([]*Token, error) {
	rows, err := dbConn.Query("select id, name, token from secret")
	if err != nil {
		fmt.Println("Query error", err)
		return nil, fmt.Errorf("%w", err)
	}

	defer rows.Close()
	tokens := make([]*Token, 0)
	for rows.Next() {
		var id int
		var name, token string
		err = rows.Scan(&id, &name, &token)

		if err != nil {
			fmt.Println("Error fetching", err)
		}
		tokens = append(tokens, &Token{id, name, token})
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	return tokens, nil
}

func AddToken(name, token, url string) error {
	return nil
}
