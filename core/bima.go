package bima

import (
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
}

func New(w fyne.Window) *Bima {
	registry := NewRegistry()

	return &Bima{
		Registry: registry,
		UI: &UI{
			Window: w,
		},
	}
}
