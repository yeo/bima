package bima

import (
	"database/sql"

	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type UI struct {
	Window fyne.Window
	Header *widget.Box
}

type Bima struct {
	Registry *Registry
	UI       *UI
	DB       *sql.DB
}

func New(w fyne.Window, db *sql.DB) *Bima {
	registry := NewRegistry()

	return &Bima{
		Registry: registry,
		UI: &UI{
			Window: w,
		},
		DB: db,
	}
}
