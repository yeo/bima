// Sync makes sure our data is backed up the backend
//
package sync

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/yeo/bima/dto"
)

type Sync struct {
	Client  *http.Client
	Done    chan bool
	AppID   string
	syncURL string
}

type SyncResponse struct {
	Current []*dto.Token `json:"current"`
	Removed []*dto.Token `json:"removed"`
}

type SyncRequest struct {
	Current []*dto.Token `json:"current"`
	Removed []*dto.Token `json:"removed"`
}

func New(appID string, syncURL string) *Sync {
	if url := os.Getenv("SYNC_URL"); url != "" {
		syncURL = url
	}

	return &Sync{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Done:    make(chan bool),
		AppID:   appID,
		syncURL: syncURL,
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
			s.Do()
		}
	}
}

// Do sends a sync request to server and update local db accordingly
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

	log.Printf("Remove token", removedTokens)

	payload, err := json.Marshal(syncRequest)
	if err != nil {
		log.Printf("Cannot marshal", tokens)
		return
	}

	req, err := http.NewRequest("POST", s.syncURL, bytes.NewBuffer(payload))
	req.Header.Set("User-Agent", "bima")
	req.Header.Set("AppID", s.AppID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)

	if err != nil {
		log.Printf("Cannot post to backend. Retry later")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	log.Printf("Response", string(body))
	var diff SyncResponse
	err = json.Unmarshal(body, &diff)

	if resp.StatusCode == 200 {
		if diff.Current != nil {
			for _, t := range diff.Current {
				log.Printf("Get current token", t)
				dto.InsertOrReplaceSecret(t)
			}
		}

		if diff.Removed != nil {
			for _, t := range diff.Removed {
				log.Printf("Remove token", t)
				dto.CommitDeleteSecret(t.ID)
			}
		}
	}
}
