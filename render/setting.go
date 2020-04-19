package render

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/dto"
)

func DrawSetting(bima *bima.Bima) *widget.Button {
	bima.AppModel.CurrentScreen = "setting"
	//canvas := bima.UI.Window.Canvas()
	button := widget.NewButton("Setting", func() {

		appIDEntry := &widget.Entry{
			Text: bima.Registry.AppID,
		}
		appIDWidget := widget.NewVBox(
			widget.NewLabel("App ID"),
			appIDEntry,
		)

		syncEntry := &widget.Entry{
			Text: bima.Registry.SyncURL,
		}
		backend := widget.NewVBox(
			widget.NewLabel("Sync URL"),
			syncEntry,
		)

		emailEntry := &widget.Entry{
			Text: bima.Registry.Email,
		}
		email := widget.NewVBox(
			widget.NewLabel("Email"),
			emailEntry,
		)

		actionButtons := widget.NewHBox(
			widget.NewButton("Save", func() {
				bima.Registry.AppID = appIDEntry.Text
				bima.Registry.SyncURL = syncEntry.Text
				bima.Registry.Email = emailEntry.Text
				dto.SavePrefs(map[string]string{
					"email": bima.Registry.Email,
				})
			}),
			widget.NewButton("Back", func() {
				bima.UI.Window.SetContent(bima.UI.MainContainer)
				DrawCode(bima)
			}),
		)
		syncWidget := widget.NewHBox(
			&widget.Check{
				Checked: false,
				Text:    "Enable Sync",
			},
		)

		container := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
		container.AddObject(appIDWidget)
		container.AddObject(email)
		container.AddObject(backend)
		container.AddObject(syncWidget)
		container.AddObject(actionButtons)
		container.AddObject(layout.NewSpacer())

		bima.UI.Window.SetContent(container)
		container.Refresh()
	})

	return button
}
