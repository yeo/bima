package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"

	"github.com/rs/zerolog/log"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/db"
	"github.com/yeo/bima/dto"
	"github.com/yeo/bima/render"
)

func main() {
	bima.InitLog()

	dbCon, err := db.Setup()
	dto.SetDB(dbCon)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot setup db")
	}

	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())
	bima := bima.New(a, dbCon)
	render.Render(bima)

	bima.Cleanup()
}
