package bima

import (
	"fmt"

	"github.com/satori/go.uuid"

	"github.com/yeo/bima/dto"
)

const (
	CfgMasterPassword = "master_password"
	CfgAppID          = "app_id"
	ScopeCore         = "bima"
)

type Registry struct {
	// AppID is to identify who this is when syncing with our backend
	// App on different platform shares this to sync data
	AppID          string
	MasterPassword string
}

func NewRegistry() *Registry {
	r := Registry{}

	// TODO: Load
	if config, err := dto.GetConfig(CfgAppID, ScopeCore); err == nil && config != nil {
		r.AppID = config.Value
	} else {
		u := uuid.NewV4()
		fmt.Printf("Generated appID: %s\n", u.String())
		dto.UpdateConfig(CfgAppID, u.String(), ScopeCore)
		r.AppID = u.String()
	}

	return &r
}

func (r *Registry) SaveMasterPassword(password string) error {
	if r.MasterPassword == "" {
		// First time ever set password
		r.MasterPassword = password
		return nil
	}

	// Here we already has a password we need to re-encrypt data with new key
	return nil
}
