package main

import (
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
	}

	addButton := render.DrawNewCode(bima)
	addButton.Resize(fyne.NewSize(50, 20))
	settingButton := render.DrawSetting(bima)

	header := widget.NewHBox(searchBox, addButton, settingButton)
	header.Resize(fyne.NewSize(350, 50))
	bima.UI.Header = header

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	bima.UI.MainContainer = fyne.NewContainerWithLayout(layout.NewGridLayout(1), header, codeContainer)

	if bima.Registry.MasterPassword == "" {
		render.DrawMasterPassword(bima, render.DrawMainUI)
	} else {
		render.DrawMainUI(bima)
	}

	go bima.Sync.Watch()

	w.Resize(fyne.NewSize(350, 600))
	w.ShowAndRun()

	cleanup(bima)
}

func cleanup(bima *bima.Bima) {
	bima.DB.Close()
	bima.Sync.Done <- true
}
