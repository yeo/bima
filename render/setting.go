package render

import (
	//"time"

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

	newDeviceButton := widget.NewHBox(
		widget.NewButton("Show Setup Kit", func() {
			c := NewSetupKitComponent(bima)
			bima.Push("changepassword", c)
		}),
		widget.NewButton("Get Quick Setup Code", func() {
			s := dialog.NewProgressInfinite("Getting quick setup code", "...", bima.UI.Window)
			go func() {
				// Request server, on finish
				// Close dialog
				// Then show dialog for code
				code, err := bima.Sync.ExchangeBlob(bima.Registry.GetSetupKit())
				s.Hide()
				s = nil

				if err == nil {
					code2 := ""
					for i, c := range code {
						code2 += string(c)
						if (i+1)%3 == 0 && i > 0 && (i+1) < len(code) {
							code2 += " "
						}
					}
					dialog.ShowInformation("Setup Code", code2, bima.UI.Window)
				} else {
					dialog.ShowError(err, bima.UI.Window)
				}
			}()
			s.Show()
		}),
	)

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	container.AddObject(appIDWidget)
	container.AddObject(email)
	container.AddObject(backend)
	container.AddObject(actionButtons)
	container.AddObject(changePasswordButton)
	container.AddObject(exportButtons)
	container.AddObject(newDeviceButton)
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

type SetupKitComponent struct {
	bima      *bima.Bima
	Container *fyne.Container
}

func NewSetupKitComponent(bima *bima.Bima) *SetupKitComponent {
	s := SetupKitComponent{
		bima: bima,
		Container: fyne.NewContainerWithLayout(layout.NewGridLayout(1),
			widget.NewLabel("Setup Kit.\nSave below text into a safe place\nto bootstrap a new device\nand restore data when you have no access\nto any active installation.\n"),
			&widget.Entry{
				Text:      bima.Registry.GetSetupKit(),
				MultiLine: true,
			},
			layout.NewSpacer()),
	}

	s.Container.AddObject(widget.NewButton("Back", func() {
		DrawMainUI(s.bima)
	}))

	return &s
}

func (s *SetupKitComponent) Render() fyne.CanvasObject {
	return s.Container
}

func (s *SetupKitComponent) Remove() {
	return
}
