// Sync makes sure our data is backed up the backend
//
package sync

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/yeo/bima/dto"
)

type Sync struct {
	Client *http.Client
	Done   chan bool
	AppID  string
}

type SyncResponse struct {
	Add    []*dto.Token
	Delete []string
	Update []*dto.Token
}

func New(appID string) *Sync {
	return &Sync{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
		Done:  make(chan bool),
		AppID: appID,
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

	payload, err := json.Marshal(tokens)
	if err != nil {
		log.Println("Cannot marshal", tokens)
		return
	}

	url := "https://a19e6def.ngrok.io/sync"

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
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

	log.Println(diff)
	if diff.Add != nil {
		log.Println("We will add", diff.Add)
	}
}
