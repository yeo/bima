package main

import (
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/render"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())

	bima := bima.New(a)
	render.Render(bima)

	bima.Cleanup()
}
