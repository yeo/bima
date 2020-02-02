package render

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/dto"
)

func DrawNewCode(canvas fyne.Canvas) *widget.Button {
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

			if dto.AddSecret(name, code) == nil {
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
