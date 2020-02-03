package render

import (
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/dto"
)

func DrawNewCode(bima *bima.Bima) *widget.Button {
	canvas := bima.UI.Window.Canvas()

	var popup *widget.PopUp

	nameEntry := &widget.Entry{
		PlaceHolder: "Name",
	}
	codeEntry := &widget.Entry{
		PlaceHolder: "OTP Secret",
	}

	content := widget.NewHBox(
		nameEntry,
		codeEntry,
		widget.NewButton("Save", func() {
			name := nameEntry.Text
			code := codeEntry.Text

			if dto.AddSecret(name, code, bima.Registry.MasterPassword) == nil {
				if popup != nil {
					popup.Hide()
					popup = nil
				}
			} else {
				// TODO: Error handler
			}
		}),
	)

	addButton := widget.NewButton("Add", func() {
		popup = widget.NewModalPopUp(content, canvas)
	})

	return addButton
}
