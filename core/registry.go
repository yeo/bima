package bima

import (
	//"crypto/rand"

	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"

	"github.com/yeo/bima/dto"
)

const (
	CfgMasterPassword = "master_password"
	CfgAppID          = "app_id"
	CfgSyncURL        = "sync_url"
	ScopeCore         = "bima"
	CfgEncryptionKey  = "encryption_key"
)

type Registry struct {
	// AppID is to identify who this is when syncing with our backend
	// App on different platform shares this to sync data
	AppID string

	MasterPassword string
	HasSetPassword string

	SyncURL string
	Email   string
}

func NewRegistry() *Registry {
	r := Registry{}

	prefs, err := dto.LoadPrefs()
	if err != nil {
		log.Error().Msg("Error when loading user preferences")
	}

	config, err := dto.GetConfig(CfgAppID, ScopeCore)
	if err == nil && config != nil {
		log.Debug().Str("appid", config.Value).Msg("Found existed appid")
		r.AppID = config.Value
	} else {
		u, _ := uuid.NewV4()
		log.Debug().Str("appid", u.String()).Msg("Generated appid")
		dto.UpdateConfig(CfgAppID, u.String(), ScopeCore)
		r.AppID = u.String()
	}

	// default sync url
	r.SyncURL = "http://bima.getopty.com/api/sync"
	if syncURL, err := dto.GetConfig(CfgSyncURL, ScopeCore); err == nil && syncURL != nil {
		r.SyncURL = syncURL.Value
	}

	if prefs["email"] != nil {
		r.Email = prefs["email"].Value
	}
	if prefs["has_set_password"] != nil {
		r.HasSetPassword = prefs["has_set_password"].Value
	}

	return &r
}

func (r *Registry) Save() error {
	return nil
}

func (r *Registry) ChangeSyncURL(url string) error {
	r.SyncURL = url
	dto.UpdateConfig(CfgSyncURL, url, ScopeCore)

	return nil
}

func (r *Registry) SaveMasterPassword(password string) error {
	if r.MasterPassword == "" {
		log.Info().Msg("Save fresh new password")
		// First time we ever set password
		r.MasterPassword = password

		prefs := map[string]string{
			"has_set_password": "y",
		}
		if err := dto.SavePrefs(prefs); err != nil {
			panic("Cannot save to sqlite")
		}

		return nil
	}

	log.Info().Msg("Update an existing password")
	if err := dto.ChangePassword(r.MasterPassword, password); err == nil {
		r.MasterPassword = password
	} else {
		return err
	}

	return nil
}
