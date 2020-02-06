package render

import (
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
)

func DrawSetting(bima *bima.Bima, done func(*bima.Bima)) *widget.Button {
	//canvas := bima.UI.Window.Canvas()

	backendEntry := &widget.Entry{
		PlaceHolder: "Enter Master Password",
		Text:        "https://bima.getopty.com/api",
	}
	passwordContent := widget.NewHBox(
		backendEntry,
		widget.NewButton("Save", func() {
			done(bima)
		}),
	)

	button := widget.NewButton("Setting", func() {
		bima.UI.Window.SetContent(passwordContent)
	})

	return button
}
