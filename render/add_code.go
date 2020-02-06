package render

import (
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/dto"
)

func DrawFormCode(bima *bima.Bima, token *dto.Token, done func(token *dto.Token)) *widget.Box {

	nameEntry := &widget.Entry{
		PlaceHolder: "Name, eg: awesome_username",
		Text:        token.Name,
	}

	urlEntry := &widget.Entry{
		PlaceHolder: "URL, eg: github.com",
		Text:        token.URL,
	}

	codeEntry := &widget.Entry{
		PlaceHolder: "OTP Secret",
	}
	if token.ID != "" {
		codeEntry.Hide()
	}

	content := widget.NewVBox(
		nameEntry,
		urlEntry,
		codeEntry,
		widget.NewButton("Save", func() {
			token.Name = nameEntry.Text
			token.URL = urlEntry.Text
			if token.ID == "" {
				// When token is already save, we don't allow to change token anymore. One has to delete and resync
				token.RawToken = codeEntry.Text
			}

			done(token)
		}),
	)

	return content
}

func DrawNewCode(bima *bima.Bima) *widget.Button {
	var popup *widget.PopUp
	canvas := bima.UI.Window.Canvas()
	content := DrawFormCode(bima, &dto.Token{}, func(token *dto.Token) {
		if dto.AddSecret(token, bima.Registry.MasterPassword) == nil {
			if popup != nil {
				popup.Hide()
				popup = nil
			}
		} else {
			// TODO: Error handler
		}
	})
	addButton := widget.NewButton("Add", func() {
		popup = widget.NewPopUp(content, canvas)
	})

	return addButton
}

func DrawEditCode(bima *bima.Bima, token *dto.Token) *widget.Button {
	var popup *widget.PopUp
	canvas := bima.UI.Window.Canvas()

	content := DrawFormCode(bima, token, func(token *dto.Token) {
		popup.Hide()
	})

	return widget.NewButton("Edit", func() {
		popup = widget.NewPopUp(content, canvas)
	})
}
