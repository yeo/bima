// Sync makes sure our data is backed up the backend
//
package sync

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yeo/bima/dto"
)

type Sync struct {
	Client  *http.Client
	Done    chan bool
	AppID   string
	SyncURL string
}

type SyncResponse struct {
	Current []*dto.Token `json:"current"`
	Removed []string     `json:"removed"`
}

type SyncRequest struct {
	Current []*dto.Token `json:"current"`
	Removed []string     `json:"removed"`
}

func New(appID string) *Sync {
	syncURL := os.Getenv("SYNC_URL")
	if syncURL == "" {
		syncURL = "http://localhost:4000/api/sync"
	}

	return &Sync{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Done:    make(chan bool),
		AppID:   appID,
		SyncURL: syncURL,
	}
}

// Watch setup a timer to routily send a Sync request to serer
func (s *Sync) Watch() {
	ticker := time.NewTicker(30 * time.Second)

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
		log.Println("Cannot fetch token")
		return
	}

	removeTokenIDs := make([]string, 0)
	if removedTokens, err := dto.LoadDeleteTokens(); err == nil {
		for _, t := range removedTokens {
			removeTokenIDs = append(removeTokenIDs, t.ID)
		}
	}

	syncRequest := SyncRequest{
		Current: tokens,
		Removed: removeTokenIDs,
	}

	payload, err := json.Marshal(syncRequest)
	if err != nil {
		log.Println("Cannot marshal", tokens)
		return
	}

	req, err := http.NewRequest("POST", syncURL, bytes.NewBuffer(payload))
	req.Header.Set("User-Agent", "bima")
	req.Header.Set("AppID", s.AppID)
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.Do(req)

	if err != nil {
		log.Println("Cannot post to backend. Retry later")
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	log.Println("Response", string(body))
	var diff SyncResponse
	err = json.Unmarshal(body, &diff)

	if resp.StatusCode == 200 {
		// Actually removed deleted token since the delete request are synced
		//dto.CommitDeleteToken()
	}

	if diff.Current != nil {
		for _, t := range diff.Current {
			log.Println("token", t)
		}
	}
}
