package dto

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/satori/go.uuid"

	"github.com/yeo/bima/shield"
)

type Token struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	RawToken  string // Use when user input the plain text token before we save and encrypt it into Token field
	Token     []byte `json:"token"`
	URL       string `json:"url"`
	Version   int    `json:"version"`
	DeletedAt int64  `json:"deleted_at"`
}

func (t *Token) DecryptToken(masterPassword string) string {
	return string(shield.Decrypt(t.Token, masterPassword))

}

func queryTokens(query string) ([]*Token, error) {
	tokens := make([]*Token, 0)
	rows, err := dbConn.Query(query)
	if err != nil {
		fmt.Println("Query error", err)
		return nil, fmt.Errorf("%w", err)
	}

	defer rows.Close()
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

func LoadDeleteTokens() ([]*Token, error) {
	return queryTokens("select id, name, token, url, version from secret where deleted_at is NOT NULL")
}

func LoadTokens() ([]*Token, error) {
	return queryTokens("select id, name, token, url, version from secret where deleted_at is NULL")
}

func CommitDeleteSecret(id string) error {
	log.Println("Commit to delete", id)

	r, err := dbConn.Exec("DELETE FROM secret WHERE id=?", id)

	log.Println("Delete affect", r, "rows", "error", err)
	return err
}

func LoadSecretByID(id string) (*Token, error) {
	stmt, err := dbConn.Prepare("select name, token, url, version, deleted_at from secret where id = ?")

	if err != nil {
		return nil, fmt.Errorf("Cannot prepare statement: %+w", err)
	}

	defer stmt.Close()

	var name, url string
	var token []byte
	var version int
	var deletedAt sql.NullInt64
	err = stmt.QueryRow(id).Scan(&name, &token, &url, &version, &deletedAt)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:

		t := &Token{
			ID:      id,
			Name:    name,
			Token:   token,
			URL:     url,
			Version: version,
		}

		if deletedAtValue, err2 := deletedAt.Value(); err2 == nil {
			if deletedAtValue != nil {
				t.DeletedAt = deletedAtValue.(int64)
			}
		}

		return t, nil
	default:
		return nil, fmt.Errorf("Error when querying database %+w", err)
	}
}

func DeleteSecret(token *Token) error {
	log.Println("Delete token", token)
	r, err := dbConn.Exec("UPDATE secret SET deleted_at = datetime('now'), version = version + 1 WHERE id=?", token.ID)

	log.Println("Mark for deletion result", r, err)
	return err
}

func UpdateSecret(token *Token) error {
	log.Println(token)
	r, err := dbConn.Exec("UPDATE secret SET name = ?, url = ?, version = version + 1 WHERE id=?", token.Name, token.URL, token.ID)

	log.Println("Update result", r, err)
	return err
}

func InsertOrReplaceSecret(token *Token) error {
	currentToken, err := LoadSecretByID(token.ID)
	if err != nil {
		log.Println("Error when fetching current token for", token, err)
		return fmt.Errorf("Cannot fetch current secret %+w", err)
	}

	if currentToken == nil || (token.Version > currentToken.Version) {
		log.Println("Insert or replace", token)
		r, err := dbConn.Exec("INSERT OR REPLACE INTO secret(id, name, url, token, version) VALUES(?, ?, ?, ?, ?)", token.ID, token.Name, token.URL, token.Token, token.Version)
		log.Println("Insert or Replace result", r, err)
		return err
	}

	log.Println("Token", currentToken.ID, "is already newer. Current version=", currentToken.Version, "Request Version=", token.Version, token.ID)

	return nil
}

func AddSecret(token *Token, masterPassword string) error {
	name := token.Name
	url := token.URL
	rawToken := token.RawToken

	if name == "" || url == "" || rawToken == "" {
		return errors.New("Invalid input")
	}

	tx, err := dbConn.Begin()

	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO secret(id, name, url, token) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	u, _ := uuid.NewV4()
	if err != nil {
		fmt.Printf("Something went wrong: %s", err)
		return fmt.Errorf("Error when generating uuid %+w", err)
	}

	encryptedToken := shield.Encrypt([]byte(rawToken), masterPassword)

	_, err = stmt.Exec(u.String(), name, url, string(encryptedToken))
	if err != nil {
		return fmt.Errorf("Error when executing statement %+w", err)
	}
	tx.Commit()

	return nil
}
