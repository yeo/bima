package dto

import (
	"errors"
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
	URL     string
	Version int
}

func (t *Token) DecryptToken(masterPassword string) string {
	return string(shield.Decrypt(t.Token, masterPassword))

}

func LoadTokens() ([]*Token, error) {
	rows, err := dbConn.Query("select id, name, token, version, url  from secret")
	if err != nil {
		fmt.Println("Query error", err)
		return nil, fmt.Errorf("%w", err)
	}

	defer rows.Close()
	tokens := make([]*Token, 0)
	for rows.Next() {
		var id, name, url string
		var token []byte
		var version int
		err = rows.Scan(&id, &name, &token, &url, &version)

		if err != nil {
			fmt.Println("Error fetching", err)
		}
		tokens = append(tokens, &Token{ID: id, Name: name, Token: token, URL: url, Version: version})
	}

	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

	return tokens, nil
}

func AddSecret(name, url, token, masterPassword string) error {
	if name == "" || url == "" || token == "" {
		return errors.New("Invalid input")
	}

	tx, err := dbConn.Begin()

	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO secret(id, name, url, token) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	u := uuid.NewV4()
	//	if err != nil {
	//		fmt.Printf("Something went wrong: %s", err)
	//		return fmt.Errorf("Error when generating uuid %+w", err)
	//	}

	encryptedToken := shield.Encrypt([]byte(token), masterPassword)

	_, err = stmt.Exec(u.String(), name, url, string(encryptedToken))
	if err != nil {
		return fmt.Errorf("Error when executing statement %+w", err)
	}
	tx.Commit()

	return nil
}
