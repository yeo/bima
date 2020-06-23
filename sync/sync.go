// Sync makes sure our data is backed up the backend
//
package sync

import (
	"bytes"
	"encoding/json"
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

func New(appID, appVersion, apiURL string) *Sync {
	return &Sync{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Done:       make(chan bool),
		AppID:      appID,
		apiURL:     apiURL,
		AppVersion: appVersion,
	}
}

// Watch setup a timer to routily send a Sync request to serer
func (s *Sync) Watch() {
	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-s.Done:
			return
		case <-ticker.C:
			//log.Debug().Msg("Disable sync for now")
			s.Do()
		}
	}
}

func (s *Sync) buildRequest(method, path string, payload io.Reader) (*http.Request, error) {
	// Check trailing slash
	url := s.apiURL + "/" + path
	req, err := http.NewRequest(method, url, payload)
	log.Debug().Str("url", url).Msg("Build Request")
	req.Header.Set("User-Agent", "bima")
	req.Header.Set("AppID", s.AppID)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("App-Version", s.AppVersion)
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

// Currently this use http sync. but we should switch to websocket
func (s *Sync) Do() {
	// Send a POST request with all of our token to backend
	// Backend return a diff of action to let us know what to do.
	// Eg
	// - Remove from server
	// - Updated field
	// - Add new
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
	}
}
