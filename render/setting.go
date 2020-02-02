package render

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

func DrawSetting(canvas fyne.Canvas) *widget.Button {
	button := widget.NewButton("Setting", func() {
	})

	return button
}
