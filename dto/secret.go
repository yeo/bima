package dto

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"

	"github.com/yeo/bima/shield"
)

type Token struct {
	ID      string
	Name    string
	Token   []byte
	Version int
}

func (t *Token) DecryptToken() string {
	return string(shield.Decrypt(t.Token, "godfather"))

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
		var id, name string
		var token []byte
		err = rows.Scan(&id, &name, &token)

		if err != nil {
			fmt.Println("Error fetching", err)
		}
		tokens = append(tokens, &Token{id, name, token, 1})
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	return tokens, nil
}

func AddSecret(name, token string) error {
	tx, err := dbConn.Begin()

	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO secret(id, name, token) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	u := uuid.NewV4()
	//	if err != nil {
	//		fmt.Printf("Something went wrong: %s", err)
	//		return fmt.Errorf("Error when generating uuid %+w", err)
	//	}

	encryptedToken := shield.Encrypt([]byte(token), "godfather")

	_, err = stmt.Exec(u.String(), name, string(encryptedToken))
	if err != nil {
		return fmt.Errorf("Error when executing statement %+w", err)
	}
	tx.Commit()

	return nil
}
