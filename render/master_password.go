package render

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
)

func DrawMasterPassword(bima *bima.Bima, container *fyne.Container, done func(*bima.Bima, *fyne.Container)) {
	passwordEntry := &widget.Entry{
		PlaceHolder: "Enter Master Password",
	}
	passwordContent := widget.NewHBox(
		passwordEntry,
		widget.NewButton("Save", func() {
			bima.Registry.SaveMasterPassword(passwordEntry.Text)
			done(bima, container)
		}),
	)

	bima.UI.Window.SetContent(passwordContent)
}
