package render

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
)

func DrawSetting(bima *bima.Bima) *widget.Button {
	bima.AppModel.CurrentScreen = "setting"
	//canvas := bima.UI.Window.Canvas()
	button := widget.NewButton("Setting", func() {

		appIDWidget := widget.NewVBox(
			widget.NewLabel("App ID"),
			&widget.Entry{
				Text: bima.Registry.AppID,
			},
		)

		syncEntry := &widget.Entry{
			Text: bima.Registry.SyncURL,
		}
		backend := widget.NewVBox(
			widget.NewLabel("Sync URL"),
			syncEntry,
		)

		actionButtons := widget.NewHBox(
			widget.NewButton("Save", func() {
				bima.Registry.ChangeSyncURL(syncEntry.Text)
				bima.UI.Window.SetContent(bima.UI.MainContainer)
				DrawCode(bima)
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
		container.AddObject(backend)
		container.AddObject(syncWidget)
		container.AddObject(actionButtons)
		container.AddObject(layout.NewSpacer())

		bima.UI.Window.SetContent(container)
		container.Refresh()
	})

	return button
}
