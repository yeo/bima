package render

import (
	"os"
	//"time"

	"fyne.io/fyne"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
)

func AddHeader(bima *bima.Bima) {
	searchBox := &widget.Entry{
		PlaceHolder: "Search",
		MultiLine:   false,
		OnChanged: func(t string) {
			bima.AppModel.FilterText = t
		},
	}

	addButton := DrawNewCode(bima)
	settingButton := widget.NewButton("Settings", func() {
		s := NewSettingComponent(bima)
		bima.Push("settings", s)
	})

	headerWidget := widget.NewHBox(searchBox, layout.NewSpacer(), addButton, settingButton)

	header := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 50}), headerWidget)

	bima.UI.Header = header
}

func DrawMainUI(bima *bima.Bima) {
	bima.AppModel.CurrentScreen = "token/list"
	DrawCode(bima)
}

// main Entrypoint to render whole ui
func Render(bima *bima.Bima) {
	AddHeader(bima)

	// To avoid the annoying of password when debugging, we support set password via env.
	if password := os.Getenv("BIMAPASS"); password != "" {
		bima.Registry.MasterPassword = password
	}

	if bima.Registry.MasterPassword == "" {
		DrawMasterPassword(bima, DrawMainUI)
	} else {
		DrawMainUI(bima)
	}

	go bima.Sync.Watch()

	bima.UI.Window.Resize(fyne.NewSize(320, 640))
	bima.UI.Window.ShowAndRun()
}
