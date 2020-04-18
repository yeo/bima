package main

import (
	"os"

	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"

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
	a.Settings().SetTheme(theme.LightTheme())
	w := a.NewWindow("Bima")

	dbCon, err := db.Setup()
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot setup db")
	}
	dto.SetDB(dbCon)

	bima := bima.New(w, dbCon)
	render.Render(bima)

	cleanup(bima)
}

func cleanup(bima *bima.Bima) {
	bima.DB.Close()
	bima.Sync.Done <- true
}
