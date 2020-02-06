package main

import (
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
	settingButton := render.DrawSetting(bima, loadMainUI)

	header := widget.NewHBox(searchBox, addButton, settingButton)
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

func cleanup(bima *bima.Bima) {
	fmt.Println("TODO: Cleanup")

	bima.DB.Close()
	bima.Sync.Done <- true
}
