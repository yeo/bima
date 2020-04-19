package render

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/dto"
)

type SettingComponent struct {
	bima      *bima.Bima
	Container *fyne.Container
}

func NewSettingComponent(bima *bima.Bima) *SettingComponent {
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

			dialog.ShowInformation("Success", "Your information is saved.", bima.UI.Window)
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

	s := SettingComponent{
		Container: container,
		bima:      bima,
	}

	return &s
}

func (s *SettingComponent) Render() fyne.CanvasObject {
	return s.Container
}
