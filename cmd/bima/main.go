package main

import (
	//"image/color"
	"time"

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
		render.DrawMasterPassword(bima, loadMainUI)
	} else {
		loadMainUI(bima)
	}

	go bima.Sync.Watch()

	w.Resize(fyne.NewSize(350, 600))
	w.ShowAndRun()

	cleanup(bima)
}

func loadMainUI(bima *bima.Bima) {
	content := bima.UI.MainContainer
	bima.UI.Window.SetContent(content)

	go func() {
		secs := time.Now().Unix()
		remainder := secs % 30
		time.Sleep(time.Duration(30-remainder) * time.Second)
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ticker.C:
				render.DrawCode(bima)
			}
		}
	}()
	render.DrawCode(bima)
}

func cleanup(bima *bima.Bima) {
	bima.DB.Close()
	bima.Sync.Done <- true
}
