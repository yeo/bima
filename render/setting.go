package render

import (
	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/dto"
	"github.com/yeo/bima/exporter"
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
		Text: bima.Registry.ApiURL,
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
			bima.Registry.ApiURL = syncEntry.Text
			bima.Registry.Email = emailEntry.Text
			dto.UpdateConfig(dto.CfgEmail, bima.Registry.Email)
			dialog.ShowInformation("Success", "Your information is saved.", bima.UI.Window)
		}),
		layout.NewSpacer(),
		widget.NewButton("Back", func() {
			DrawCode(bima)
		}),
	)

	exportButtons := widget.NewHBox(
		widget.NewButton("Export", func() {
			// TODO: FIX this to make it run on window and write to home
			if err := exporter.Export(bima.Registry.CombineEncryptionKey(), "/tmp/bima.csv"); err == nil {
				dialog.ShowInformation("Success", "Your tokens are exported \nto /tmp/bima.csv", bima.UI.Window)
			} else {
				dialog.ShowInformation("Err", "Export fail", bima.UI.Window)
			}
		}),
		widget.NewButton("Import", func() {
			// TODO: FIX this to make it run on window and write to home
			if err := exporter.Import(bima.Registry.CombineEncryptionKey(), "/tmp/bima.csv"); err == nil {
				dialog.ShowInformation("Success", "Your tokens are imported", bima.UI.Window)
			} else {
				dialog.ShowInformation("Err", "Import fail", bima.UI.Window)
			}
		}),
	)

	changePasswordButton := widget.NewHBox(
		widget.NewButton("Change Password", func() {
			p := NewPasswordComponent(bima, ChangePasswordForm)
			bima.Push("changepassword", p)
		}),
	)

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	container.AddObject(appIDWidget)
	container.AddObject(email)
	container.AddObject(backend)
	container.AddObject(actionButtons)
	container.AddObject(changePasswordButton)
	container.AddObject(exportButtons)
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

func (s *SettingComponent) Remove() {
	return
}
