// Sync makes sure our data is backed up the backend
//
package sync

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yeo/bima/dto"
)

type Me struct {
	Email string
}

type LockBox interface {
	Encrypt(string) string
}

type Sync struct {
	Client     *http.Client
	Me         *Me
	Done       chan bool
	AppID      string
	AppVersion string
	apiURL     string
	LockBox    LockBox
	okToSync   bool
	dbVersion  int
}

type CompareResponse struct {
	AppVersion  string `json:"app_version"`
	DataVersion int    `json:"data_version"`

	RequireDataSync  bool `json:"require_data_sync"`
	RequireAppUpdate bool `json:"require_app_update"`
}

type SyncResponse struct {
	Added   []*dto.Token `json:"added"`
	Removed []*dto.Token `json:"removed"`
	Changed []*dto.Token `json:"changed"`
}

type SyncRequest struct {
	Current []*dto.Token `json:"current"`
	Removed []*dto.Token `json:"removed"`
}

type BlobRequest struct {
	Payload string `json:"payload"`
}

type BlobResponse struct {
	Code    string `json:"code"`
	Payload string `json:"payload"`
}

type MeResponse struct {
	DBVersion int `json:"db"`
}

func New(appID, appVersion, apiURL string) *Sync {
	return &Sync{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Done:       make(chan bool),
		AppID:      appID,
		apiURL:     apiURL,
		AppVersion: appVersion,
		okToSync:   false,
		dbVersion:  0,
	}
}

// Watch setup a timer to routily send a Sync request to serer
func (s *Sync) Watch() {
	ticker := time.NewTicker(15 * time.Second)

	go s.Do()
	for {
		select {
		case <-s.Done:
			return
		case <-ticker.C:
			if s.okToSync {
				s.Do()
			} else {
				log.Debug().Msg("Sync is paused. Not requesting server")
			}
		}
	}
}

func (s *Sync) PauseSync() {
	s.okToSync = false
}

func (s *Sync) ResumeSync() {
	s.okToSync = true
	s.Do()
}

func (s *Sync) buildRequest(method, path string, payload io.Reader) (*http.Request, error) {
	// Check trailing slash
	url := s.apiURL + "/" + path
	req, err := http.NewRequest(method, url, payload)
	log.Debug().Str("url", url).Msg("Build Request")
	req.Header.Set("User-Agent", "bima")
	req.Header.Set("AppID", s.AppID)
	req.Header.Set("AppVersion", s.AppVersion)
	req.Header.Set("Content-Type", "application/json")
	//req.Header.Set("Data-Version", "100")
	return req, err
}

// Compare send a light request(small payload) to server to see if we need to do a sync
func (s *Sync) Compare() *CompareResponse {
	req, err := s.buildRequest("POST", "compare", nil)
	if err != nil {
		// TODO: log error
		return nil
	}

	resp, err := s.Client.Do(req)
	if err != nil {
		log.Printf("Cannot post to backend. Retry later")
		return nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	log.Debug().Str("body", string(body)).Msg("Sync Response")
	var diff CompareResponse
	err = json.Unmarshal(body, &diff)

	return &diff
}

// Getme pings server for update and annoucement
// This is basically the only way for us to communicate with user
func (s *Sync) GetMe() (*MeResponse, error) {
	req, err := s.buildRequest("GET", "me", nil)

	resp, err := s.Client.Do(req)

	if err != nil {
		log.Printf("Cannot get /me from backend.")
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	log.Debug().Str("body", string(body)).Msg("Get Me Response")
	var me MeResponse
	err = json.Unmarshal(body, &me)
	if err != nil {
		return nil, err
	}

	return &me, nil
}

func (s *Sync) BumpDB() {
	s.dbVersion += 1

	// Once we bump db, immediately sync
	go s.Do()
}

// Currently this uses http sync. but we should switch to websocket
func (s *Sync) Do() {
	// Send a POST request with all of our token to backend
	// Backend return a diff of action to let us know what to do.
	// Eg
	// - Remove from server
	// - Updated field
	// - Add new
	if s.AppID == "" || s.AppVersion == "" {
		log.Printf("App is not ready to sync. Missing app id or app version")
		return
	}

	me, err := s.GetMe()
	if err != nil {
		log.Error().Err(err).Msg("Cannot fetch /me")
		return
	}

	if s.dbVersion == me.DBVersion {
		log.Info().Int("me", me.DBVersion).Int("sync", s.dbVersion).Msg("local db version is same as remote, ignore sync")
		return
	}
	log.Info().Int("me", me.DBVersion).Int("sync", s.dbVersion).Msg("local db version is different from remote, perform sync")

	tokens, err := dto.LoadTokens()
	if err != nil {
		log.Printf("Cannot fetch token")
		return
	}

	removedTokens, err := dto.LoadDeleteTokens()
	if err != nil {
		log.Printf("Cannot fetch deleted token", err)

	}

	syncRequest := SyncRequest{
		Current: tokens,
		Removed: removedTokens,
	}

	log.Printf("sync request %v", syncRequest)

	payload, err := json.Marshal(syncRequest)
	if err != nil {
		log.Printf("Cannot marshal", tokens)
		return
	}

	log.Print(string(payload))
	req, err := s.buildRequest("POST", "sync", bytes.NewBuffer(payload))

	resp, err := s.Client.Do(req)

	if err != nil {
		log.Printf("Cannot post to backend. Retry later")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	log.Debug().Str("body", string(body)).Msg("Sync Response")
	var diff SyncResponse
	err = json.Unmarshal(body, &diff)

	if resp.StatusCode == 200 {
		if diff.Added != nil {
			for _, t := range diff.Added {
				dto.InsertOrReplaceSecret(t)
			}
		}

		if diff.Changed != nil {
			for _, t := range diff.Changed {
				dto.UpdateSecretNoVersionBump(t)
			}
		}

		if diff.Removed != nil {
			for _, t := range diff.Removed {
				dto.CommitDeleteSecret(t.ID)
			}
		}

		s.dbVersion = me.DBVersion
	}
}

// Currently this use http sync. but we should switch to websocket
func (s *Sync) ExchangeBlob(content string) (string, error) {
	requestRaw := BlobRequest{
		Payload: content,
	}

	payload, err := json.Marshal(requestRaw)
	req, err := s.buildRequest("POST", "blob", bytes.NewBuffer(payload))
	resp, err := s.Client.Do(req)

	if err != nil {
		log.Printf("Cannot create blob. Retry later")
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	log.Debug().Str("body", string(body)).Msg("Blob Response")

	if resp.StatusCode == 200 {
		var response BlobResponse
		err = json.Unmarshal(body, &response)
		if err == nil {
			return response.Code, nil
		}
	}

	return "", err
}

func (s *Sync) GetBlob(code string) (string, error) {
	req, err := s.buildRequest("GET", "blob/"+code, nil)
	resp, err := s.Client.Do(req)

	if err != nil {
		log.Printf("Cannot create blob. Retry later")
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	log.Debug().Str("body", string(body)).Msg("Blob Response")

	if resp.StatusCode == 200 {
		var response BlobResponse
		err = json.Unmarshal(body, &response)
		return response.Payload, nil
	} else {
		return "", errors.New("Invalid code")
	}

	return "", err
}
