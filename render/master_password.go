package render

import (
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
)

func DrawMasterPassword(bima *bima.Bima, done func(*bima.Bima)) {
	passwordEntry := &widget.Entry{
		PlaceHolder: "Enter Master Password",
	}
	passwordContent := widget.NewHBox(
		passwordEntry,
		widget.NewButton("Save", func() {
			bima.Registry.SaveMasterPassword(passwordEntry.Text)
			done(bima)
		}),
	)

	bima.UI.Window.SetContent(passwordContent)
}
