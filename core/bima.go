package bima

import (
	"database/sql"

	"fyne.io/fyne"
	//"fyne.io/fyne/widget"

	"github.com/yeo/bima/sync"
)

type AppState int

type AppModel struct {
	FilterText    string
	CurrentScreen string
}

const (
	Init    AppState = 0 // Initialize the whole application
	Pending AppState = 1 // Finish init, pending drawing ui
	Ready   AppState = 2 // Everything is loaded, ready to draw ui
)

type UI struct {
	Window        fyne.Window
	Header        *fyne.Container
	MainContainer *fyne.Container
}

type Bima struct {
	Registry *Registry
	UI       *UI
	DB       *sql.DB
	Sync     *sync.Sync
	AppState AppState
	AppModel *AppModel
}

func New(w fyne.Window, db *sql.DB) *Bima {
	registry := NewRegistry()

	return &Bima{
		Registry: registry,
		UI: &UI{
			Window: w,
		},
		DB:       db,
		Sync:     sync.New(registry.AppID, registry.SyncURL),
		AppModel: &AppModel{},
	}
}
