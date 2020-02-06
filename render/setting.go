package render

import (
	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
)

func DrawSetting(bima *bima.Bima) *widget.Button {
	//canvas := bima.UI.Window.Canvas()

	appIDWidget := widget.NewVBox(
		widget.NewLabel("App ID"),
		&widget.Entry{
			Text: bima.Registry.AppID,
		},
	)

	backend := widget.NewVBox(
		widget.NewLabel("Sync URL"),
		&widget.Entry{
			Text: "https://bima.getopty.com/api",
		},
	)

	saveButton := widget.NewHBox(
		widget.NewButton("Save", func() {
			bima.UI.Window.SetContent(bima.UI.MainContainer)
			bima.UI.MainContainer.Refresh()
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
	container.AddObject(saveButton)

	button := widget.NewButton("Setting", func() {
		bima.UI.Window.SetContent(container)
	})

	return button
}
