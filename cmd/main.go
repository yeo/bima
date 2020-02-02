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

	"github.com/yeo/bima/db"
	"github.com/yeo/bima/dto"
	"github.com/yeo/bima/render"
)

func main() {
	a := app.New()
	w := a.NewWindow("primary")

	searchBox := &widget.Entry{
		PlaceHolder: "Search",
		MultiLine:   false,
	}

	addButton := render.DrawNewCode(w.Canvas())
	settingButton := render.DrawSetting(w.Canvas())

	header := widget.NewHBox(searchBox, addButton, settingButton)

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(1), header, codeContainer)

	w.SetContent(container)

	dbCon, err := db.Setup()
	if err != nil {
		log.Println("error", err)
	}
	dto.SetDB(dbCon)

	go func() {
		secs := time.Now().Unix()
		remainder := secs % 30
		time.Sleep(time.Duration(30-remainder) * time.Second)
		log.Println("Sleep to rounded time", remainder)
		render.DrawCode(w, header)
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ticker.C:
				render.DrawCode(w, header)
			}
		}
	}()
	render.DrawCode(w, header)

	w.Resize(fyne.NewSize(350, 600))
	w.ShowAndRun()
	cleanup(dbCon)
}

func cleanup(dbCon *sql.DB) {
	fmt.Println("TODO: Cleanup")

	dbCon.Close()
}
