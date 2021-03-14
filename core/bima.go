package bima

import (
	"database/sql"
	//"os"

	"fyne.io/fyne/v2"
	//"fyne.io/fyne/widget"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yeo/bima/db"
	"github.com/yeo/bima/dto"
	"github.com/yeo/bima/sync"
)

const (
	AppVersion = "0.1.0"
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

func New(a fyne.App) *Bima {
	registry := NewRegistry()
	InitLog(registry)
	db, err := InitDB(registry)

	if err != nil {
		log.Fatal().Err(err).Msg("Cannot setup db")
	}

	registry.LoadConfigsFromDB()

	title := "Bima " + AppVersion
	if registry.IsDebug() {
		title = "Bima Debug " + AppVersion
	}
	w := a.NewWindow(title)

	syncer := sync.New(registry.AppID, AppVersion, registry.ApiURL)
	return &Bima{
		Registry: registry,
		UI: &UI{
			Window: w,
		},
		DB:       db,
		Sync:     syncer,
		AppModel: &AppModel{},
	}
}

func (b *Bima) Push(name string, c Component) {
	log.Info().Str("path", name).Msg("Push component")
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

func (b *Bima) Cleanup() {
	b.DB.Close()
	b.Sync.Done <- true
}

func InitLog(r *Registry) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if r.IsDebug() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
}

func InitDB(r *Registry) (*sql.DB, error) {
	dbCon, err := db.Setup(r.dbFile)
	dto.SetDB(dbCon)
	if err != nil {
		return nil, err
	}

	return dbCon, nil
}
