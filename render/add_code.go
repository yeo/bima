package render

import (
	"image/color"
	"log"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/dto"
)

func DrawFormCode(bima *bima.Bima, token *dto.Token, done func(token *dto.Token)) fyne.CanvasObject {
	nameEntry := &widget.Entry{
		PlaceHolder: "",
		Text:        token.Name,
	}

	urlEntry := &widget.Entry{
		PlaceHolder: "",
		Text:        token.URL,
	}

	codeEntry := &widget.Entry{
		PlaceHolder: "",
	}
	if token.ID != "" {
		codeEntry.Hide()
	}

	content := widget.NewVBox(
		layout.NewSpacer(),
		canvas.NewText("Name", color.RGBA{135, 0, 16, 255}),
		nameEntry,
		canvas.NewText("URL", color.RGBA{135, 0, 16, 255}),
		urlEntry,
		canvas.NewText("OTP Secret", color.RGBA{135, 0, 16, 255}),
		codeEntry,
		widget.NewButton("Save", func() {
			token.Name = nameEntry.Text
			token.URL = urlEntry.Text
			if token.ID == "" {
				// When token is already save, we don't allow to change token anymore. One has to delete and resync
				token.RawToken = codeEntry.Text
			}

			done(token)
			nameEntry.SetText("")
			urlEntry.SetText("")
			codeEntry.SetText("")
		}),
		layout.NewSpacer(),
		widget.NewButton("Close", func() {
			done(nil)
		}),
	)

	contentLayout := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 400}), content)

	return contentLayout
}

func DrawNewCode(bima *bima.Bima) *widget.Button {
	var popup *widget.PopUp
	canvas := bima.UI.Window.Canvas()
	content := DrawFormCode(bima, &dto.Token{}, func(token *dto.Token) {
		if token == nil {
			popup.Hide()
			popup = nil
			DrawCode(bima)
			return
		}

		if dto.AddSecret(token, bima.Registry.CombineEncryptionKey()) == nil {
			if popup != nil {
				popup.Hide()
				popup = nil
			}
			DrawCode(bima)
			go bima.Sync.Do()
		} else {
			// TODO: Error handler
		}
	})
	addButton := widget.NewButton("Add", func() {
		popup = widget.NewModalPopUp(content, canvas)
	})

	return addButton
}

func DrawEditCode(bima *bima.Bima, token *dto.Token) *widget.Button {
	var popup *widget.PopUp
	canvas := bima.UI.Window.Canvas()

	content := DrawFormCode(bima, token, func(token *dto.Token) {
		if token == nil {
			popup.Hide()
			popup = nil
			DrawCode(bima)
			return
		}
		log.Println("Delete at for token", token.DeletedAt)

		if err := dto.UpdateSecret(token); err == nil {
			go bima.Sync.Do()
			if popup != nil {
				popup.Hide()
			}
		}

		DrawCode(bima)
	})

	return widget.NewButton("Edit", func() {
		popup = widget.NewModalPopUp(content, canvas)
	})
}
