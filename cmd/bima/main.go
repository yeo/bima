package main

import (
	"database/sql"
	"fmt"
	//"image/color"
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	//"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/db"
	"github.com/yeo/bima/dto"
	"github.com/yeo/bima/render"
)

func main() {
	a := app.New()
	w := a.NewWindow("primary")

	dbCon, err := db.Setup()
	if err != nil {
		log.Println("error", err)
	}
	dto.SetDB(dbCon)

	bima := bima.New(w, dbCon)

	searchBox := &widget.Entry{
		PlaceHolder: "Search",
		MultiLine:   false,
	}

	addButton := render.DrawNewCode(bima)
	settingButton := render.DrawSetting(w.Canvas())

	header := widget.NewHBox(searchBox, addButton, settingButton)
	bima.UI.Header = header

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(1), header, codeContainer)

	if bima.Registry.MasterPassword == "" {
		render.DrawMasterPassword(bima, container, loadMainUI)
	} else {
		loadMainUI(bima, container)
	}

	w.Resize(fyne.NewSize(350, 600))
	w.ShowAndRun()
	cleanup(dbCon)
}

func loadMainUI(bima *bima.Bima, content *fyne.Container) {
	bima.UI.Window.SetContent(content)

	go func() {
		secs := time.Now().Unix()
		remainder := secs % 30
		time.Sleep(time.Duration(30-remainder) * time.Second)
		log.Println("Sleep to rounded time", remainder)
		render.DrawCode(bima)
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

func cleanup(dbCon *sql.DB) {
	fmt.Println("TODO: Cleanup")

	dbCon.Close()
}
