package bima

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	//"strconv"

	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"

	"github.com/yeo/bima/dto"
	"github.com/yeo/bima/shield"
)

type Registry struct {
	// AppID is to identify who this is when syncing with our backend
	// App on different platform shares this to sync data
	AppID string

	// We used secretkey+masterpassword to encrypt data
	MasterPassword     []byte // this field won't persist to db
	SecretKey          []byte // this field is persisted to sqlite db but never send to server. used take care of backup it
	encryptedSecretKey []byte // raw encrypted string in the database before we decrypt
	combineKey         []byte // the concat of secretkey+masterpassword for final encryption key

	ApiURL string
	Email  string

	dbFile string
	debug  bool
}

func NewRegistry() *Registry {
	r := Registry{
		ApiURL: "http://bima.getopty.com/api",
		debug:  false,
		dbFile: "bima.db",
	}

	if debugFlag := os.Getenv("DEBUG"); debugFlag == "1" {
		r.debug = true
		r.dbFile = "dev.db"
	}

	// To avoid the annoying of password when debugging, we support set password via env.
	if password := os.Getenv("BIMAPASS"); password != "" {
		r.MasterPassword = []byte(password)
	}

	return &r
}

func (r *Registry) LoadConfigsFromDB() error {
	configs, err := dto.LoadConfigs()

	if err != nil {
		// TODO: Show ui error
		log.Error().Msg("Error when loading config")

		return err
	}

	if configs.AppID != "" {
		log.Debug().Str("appid", configs.AppID).Msg("Found existed appid")
		r.AppID = configs.AppID

		r.encryptedSecretKey = []byte(configs.SecretKey)
		if len(r.MasterPassword) > 0 {
			r.CombineEncryptionKey()
		}
	} else {
		u, _ := uuid.NewV4()
		log.Debug().Str("appid", u.String()).Msg("Generated appid")
		r.AppID = u.String()
		dto.UpdateConfig(dto.CfgAppId, r.AppID)
	}

	// The URL in db has highest priority
	if configs.ApiURL != "" {
		r.ApiURL = configs.ApiURL
	}

	// If we happen to have URL in env, pick it but not save to db
	if url := os.Getenv("BIMAURL"); url != "" {
		r.ApiURL = url
	}

	return nil
}

func (r *Registry) IsDebug() bool {
	return r.debug
}

func (r *Registry) Save() error {
	// Re-encrypt our secret key
	dto.UpdateConfig(dto.CfgAppId, r.AppID)
	r.encryptedSecretKey = shield.Encrypt(r.SecretKey, r.MasterPassword)
	dto.UpdateConfig(dto.CfgSecretKey, string(r.encryptedSecretKey))
	r.CombineEncryptionKey()

	return nil
}

func (r *Registry) HasOnboard() bool {
	return r.AppID != "" && r.encryptedSecretKey != nil
}

// EncryptionKeyPair returns concat of secretkey+use-master-password
func (r *Registry) CombineEncryptionKey() []byte {
	if len(r.combineKey) > 0 {
		return r.combineKey
	}

	if s, e := shield.Decrypt(r.encryptedSecretKey, r.MasterPassword); e == nil {
		r.SecretKey = s
		r.combineKey = append(r.SecretKey, r.MasterPassword...)
	}

	return r.combineKey
}

func (r *Registry) SaveMasterPassword(password string) error {
	//TODO: Error handlera

	// When user just finish onboard or unlock the app. As in, this is the first time they enter password when the app starts
	if r.MasterPassword == nil {
		log.Info().Msg("Save fresh new password")
		r.MasterPassword = []byte(password)

		if r.encryptedSecretKey == nil {
			log.Info().Msg("Generate Secret Key")
			// Save the encryption key
			// Generate a 24 byte keys length
			c := 24
			b := make([]byte, c)
			_, err := rand.Read(b)
			if err != nil {
				log.Error().Err(err).Msg("Cannot generate encryption key")
				return err
			}
			r.SecretKey = b
			encryptedSecretKey := shield.Encrypt(b, r.MasterPassword)
			dto.UpdateConfig(dto.CfgSecretKey, string(encryptedSecretKey))
		} else {
			// We had encryptied secret key, has the master password, attempt to decrypt the secret key itself
			r.CombineEncryptionKey()
		}

		if r.SecretKey != nil {
			return nil
		}

		r.MasterPassword = nil
		return errors.New("Invalid password")
	}

	// At this point, a password already in memory with all data, when save a new password here, we have to change and re-encrypt data
	log.Info().Msg("Update an existing password")
	newEncryptionKey := append(r.SecretKey, []byte(password)...)
	if err := dto.ChangePassword(r.CombineEncryptionKey(), newEncryptionKey); err == nil {
		r.MasterPassword = []byte(password)
		// Re-encrypt our secret key
		encryptedSecretKey := shield.Encrypt(r.SecretKey, r.MasterPassword)
		r.combineKey = newEncryptionKey
		dto.UpdateConfig(dto.CfgSecretKey, string(encryptedSecretKey))
	} else {
		return err
	}

	return nil
}

func (r *Registry) GetSetupKit() string {
	encryptedSecretKey := shield.Encrypt(r.SecretKey, r.MasterPassword)
	encryptedAppID := shield.Encrypt([]byte(r.AppID), r.MasterPassword)

	payload := fmt.Sprintf("{\n\"appID\":\n\"%s\",\n\"secretKey\":\n\"%s\"\n}",
		base64.StdEncoding.EncodeToString(encryptedAppID), base64.StdEncoding.EncodeToString(encryptedSecretKey))

	return payload
}
