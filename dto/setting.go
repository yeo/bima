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

func LoadConfigsByScope(scope string) (map[string]*Setting, error) {
	configs := make(map[string]*Setting)

	stmt, err := dbConn.Prepare("select name, value, scope from config where scope = ?")
	if err != nil {
		return nil, fmt.Errorf("Cannot prepare statement: %+w", err)
	}

	rows, err := stmt.Query(scope)
	if err != nil {
		return nil, fmt.Errorf("Cannot prepare statement: %+w", err)
	}
	defer stmt.Close()

	for rows.Next() {
		var name, value, scope string
		err = rows.Scan(&name, &value, &scope)
		configs[name] = &Setting{
			Name:  name,
			Value: value,
			Scope: scope,
		}

		if err != nil {
			return nil, fmt.Errorf("Error fetching config", err)
		}
	}

	return configs, nil
}

func LoadPrefs() (map[string]*Setting, error) {
	return LoadConfigsByScope("prefs")
}

func SavePrefs(p map[string]string) error {
	for k, v := range p {
		UpdateConfig(k, v, "prefs")
	}

	return nil
}

func GetConfig(configName string, configScope string) (*Setting, error) {
	stmt, err := dbConn.Prepare("select name, value, scope from config where name = ? and scope = ?")

	if err != nil {
		return nil, fmt.Errorf("Cannot prepare statement: %+w", err)
	}

	defer stmt.Close()

	var name, value, scope string
	err = stmt.QueryRow(configName, configScope).Scan(&name, &value, &scope)

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
			Scope: scope,
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
