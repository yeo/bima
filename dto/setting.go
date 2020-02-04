package dto

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Setting struct {
	Name  string
	Value string
	Scope string
}

func GetConfig(configName string, configScope string) (*Setting, error) {
	stmt, err := dbConn.Prepare("select name, value, scope from config where name = ? and scope = ?")

	if err != nil {
		return nil, fmt.Errorf("Cannot prepare statement: %+w", err)
	}

	defer stmt.Close()

	var name, value string
	err = stmt.QueryRow(configName, configScope).Scan(&name, &value)

	if err != nil {
		return nil, fmt.Errorf("Cannot query data: %+w", err)
	}

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &Setting{
			Name:  name,
			Value: value,
			Scope: configScope,
		}, nil
	default:
		return nil, fmt.Errorf("Error when querying database %+w", err)
	}
}

func UpdateConfig(configName, configValue, configScope string) error {
	if dbConn == nil {
		fmt.Println("DB is no tinit")
		return fmt.Errorf("DB is not init")
	}

	tx, err := dbConn.Begin()

	if err != nil {
		return fmt.Errorf("Cannot init transaction %+w", err)
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO config(name, value, scope) VALUES(?, ?, ?)")
	_, err = stmt.Exec(configName, configValue, configScope)

	if err != nil {
		return fmt.Errorf("Cannot commit %+w", err)
	}
	tx.Commit()

	return nil
}
