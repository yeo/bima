package bima

import (
	"database/sql"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/sync"
)

type UI struct {
	Window        fyne.Window
	Header        *widget.Box
	MainContainer *fyne.Container
}

type Bima struct {
	Registry *Registry
	UI       *UI
	DB       *sql.DB
	Sync     *sync.Sync
}

func New(w fyne.Window, db *sql.DB) *Bima {
	registry := NewRegistry()

	return &Bima{
		Registry: registry,
		UI: &UI{
			Window: w,
		},
		DB:   db,
		Sync: sync.New(registry.AppID),
	}
}
