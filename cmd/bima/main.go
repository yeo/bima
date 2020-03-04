package main

import (
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	//"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/db"
	"github.com/yeo/bima/dto"
	"github.com/yeo/bima/render"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debugFlag := os.Getenv("DEBUG"); debugFlag == "1" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	a := app.New()
	w := a.NewWindow("Bima")

	dbCon, err := db.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot setup db")
	}
	dto.SetDB(dbCon)

	bima := bima.New(w, dbCon)

	searchBox := &widget.Entry{
		PlaceHolder: "Search",
		MultiLine:   false,
		OnChanged: func(t string) {
			bima.AppModel.FilterText = t
			render.DrawCode(bima)
		},
	}

	addButton := render.DrawNewCode(bima)
	settingButton := render.DrawSetting(bima)

	headerWidget := widget.NewHBox(searchBox, layout.NewSpacer(), addButton, settingButton)

	header := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 50}), headerWidget)
	bima.UI.Header = header

	// To avoid the annoying of password when debugging, we support set password via env.
	if password := os.Getenv("BIMAPASS"); password != "" {
		bima.Registry.MasterPassword = password
	}

	if bima.Registry.MasterPassword == "" {
		render.DrawMasterPassword(bima, render.DrawMainUI)
	} else {
		render.DrawMainUI(bima)
	}
	go bima.Sync.Watch()

	w.Resize(fyne.NewSize(320, 640))
	w.ShowAndRun()

	cleanup(bima)
}

func cleanup(bima *bima.Bima) {
	bima.DB.Close()
	bima.Sync.Done <- true
}
