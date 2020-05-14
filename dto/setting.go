package dto

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CfgAppId     = "app_id"
	CfgSecretKey = "secret_key"
	CfgApiURL    = "api_url"
	CfgEmail     = "email"
)

type Setting struct {
	Name  string
	Value string
}

type SettingMap struct {
	AppID     string
	SecretKey string

	ApiURL string
	Email  string
}

func LoadConfigs() (*SettingMap, error) {
	configs := SettingMap{}

	stmt, err := dbConn.Prepare("select name, value from config")
	if err != nil {
		return nil, fmt.Errorf("Cannot prepare statement: %+w", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("Cannot prepare statement: %+w", err)
	}
	defer stmt.Close()

	for rows.Next() {
		var name, value string
		err = rows.Scan(&name, &value)

		switch name {
		case CfgAppId:
			configs.AppID = value
		case CfgSecretKey:
			configs.SecretKey = value
		case CfgApiURL:
			configs.ApiURL = value
		}

		if err != nil {
			return nil, fmt.Errorf("Error fetching config", err)
		}
	}

	return &configs, nil
}

func GetConfig(configName string) (*Setting, error) {
	stmt, err := dbConn.Prepare("select name, value, scope from config where name = ? and scope = ?")

	if err != nil {
		return nil, fmt.Errorf("Cannot prepare statement: %+w", err)
	}

	defer stmt.Close()

	var name, value string
	err = stmt.QueryRow(configName).Scan(&name, &value)

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
		}, nil
	default:
		return nil, fmt.Errorf("Error when querying database %+w", err)
	}
}

func UpdateConfig(name, value string) error {
	if dbConn == nil {
		fmt.Println("DB is not init")
		return fmt.Errorf("DB is not init")
	}

	tx, err := dbConn.Begin()

	if err != nil {
		return fmt.Errorf("Cannot init transaction %+w", err)
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO config(name, value) VALUES(?, ?)")
	_, err = stmt.Exec(name, value)

	if err != nil {
		return fmt.Errorf("Cannot commit %+w", err)
	}
	tx.Commit()

	return nil
}

func UpdateConfigMap(p map[string]string) error {
	for k, v := range p {
		UpdateConfig(k, v)
	}

	return nil
}
