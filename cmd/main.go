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

	header := widget.NewVBox(
		widget.NewLabel("Hello Fyne!"),
		widget.NewButton("Quit", func() {
			a.Quit()
		}))

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2))
	container := fyne.NewContainerWithLayout(layout.NewGridLayout(2), header, codeContainer)

	w.SetContent(container)

	dbCon, err := db.Setup()
	if err != nil {
		log.Println("error", err)
	}
	dto.SetDB(dbCon)

	go func() {
		render.DrawCode(w, header)
		secs := time.Now().Unix()
		remainder := secs % 30
		fmt.Println("we need to sleep %s", remainder)
		time.Sleep(time.Duration(remainder) * time.Second)
		render.DrawCode(w, header)
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ticker.C:
				render.DrawCode(w, header)
			}
		}
	}()

	w.Resize(fyne.NewSize(250, 600))
	w.ShowAndRun()
	cleanup(dbCon)
}

func cleanup(dbCon *sql.DB) {
	fmt.Println("TODO: Cleanup")

	dbCon.Close()
}
