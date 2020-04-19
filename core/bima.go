package bima

import (
	"database/sql"

	"fyne.io/fyne"
	//"fyne.io/fyne/widget"

	"github.com/yeo/bima/dto"
	"github.com/yeo/bima/sync"
)

type AppState int

type AppModel struct {
	FilterText    string
	CurrentScreen string
	Tokens        []*dto.Token
}

const (
	Init    AppState = 0 // Initialize the whole application
	Pending AppState = 1 // Finish init, pending drawing ui
	Ready   AppState = 2 // Everything is loaded, ready to draw ui
)

type Component interface {
	Render() fyne.CanvasObject
	Remove()
}

type UI struct {
	Window          fyne.Window
	Header          *fyne.Container
	ActiveComponent Component
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

func (b *Bima) Push(name string, c Component) {
	if b.UI.ActiveComponent != nil {
		b.UI.ActiveComponent.Remove()
		b.UI.ActiveComponent = nil
	}
	b.UI.ActiveComponent = c
	canvasObject := c.Render()
	b.AppModel.CurrentScreen = name
	b.UI.Window.SetContent(canvasObject)
	canvasObject.Refresh()
}
