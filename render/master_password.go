package render

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
)

func DrawMasterPassword(bima *bima.Bima, done func(*bima.Bima)) {
	passwordEntry := &widget.Entry{
		PlaceHolder: "Enter Master Password",
	}

	passwordField := widget.NewButton("Unlock", func() {
		bima.Registry.SaveMasterPassword(passwordEntry.Text)
		done(bima)
	})
	passwordForm := widget.NewVBox(
		layout.NewSpacer(),
		passwordEntry, passwordField,
		layout.NewSpacer(),
	)

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	container.AddObject(layout.NewSpacer())
	container.AddObject(passwordForm)
	container.AddObject(layout.NewSpacer())

	bima.UI.Window.SetContent(container)
}
