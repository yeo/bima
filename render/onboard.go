package render

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	//"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	//"github.com/rs/zerolog/log"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/shield"
)

type SetupKit struct {
	AppID     string `json:"appID"`
	SecretKey string `json:"secretKey"`
}

type OnboardComponent struct {
	bima      *bima.Bima
	Container *fyne.Container
}

func (p *OnboardComponent) Render() fyne.CanvasObject {
	return p.Container
}

func (p *OnboardComponent) Remove() {
	return
}

func NewOnboardComponent(bima *bima.Bima) *OnboardComponent {
	p := OnboardComponent{
		bima: bima,
		Container: fyne.NewContainerWithLayout(layout.NewGridLayout(1),
			layout.NewSpacer(),

			widget.NewVBox(
				widget.NewButton("Starting from scratch", func() {
					c := NewPasswordComponent(bima, NewPasswordForm)
					bima.Push("onboard", c)
				}),
				layout.NewSpacer(),
				widget.NewButton("Use Setup Code from other device", func() {
					c := NewCodeSetupComponent(bima)
					bima.Push("setup/code", c)
				}),
				layout.NewSpacer(),
				widget.NewButton("Or Restore From SetupKit", func() {
				})),

			layout.NewSpacer(),
		),
	}

	return &p
}

type CodeSetupComponent struct {
	bima      *bima.Bima
	Container *fyne.Container
}

func (p *CodeSetupComponent) Render() fyne.CanvasObject {
	return p.Container
}

func (p *CodeSetupComponent) Remove() {
	return
}

func NewCodeSetupComponent(bima *bima.Bima) *CodeSetupComponent {
	codeEntry := widget.NewEntry()
	masterPassword := widget.NewEntry()

	p := CodeSetupComponent{
		bima: bima,
		Container: fyne.NewContainerWithLayout(layout.NewGridLayout(1),
			layout.NewSpacer(),
			widget.NewVBox(
				widget.NewLabel("On the other app, go to\nSetting > Generate Setup Code\n"),
				codeEntry,
				layout.NewSpacer(),
				widget.NewLabel("Enter the same master password that you used on the other app"),
				masterPassword,
				widget.NewButton("Next", func() {
					bima.Registry.MasterPassword = []byte(masterPassword.Text)

					s := dialog.NewProgressInfinite("Getting quick setup code", "...", bima.UI.Window)
					s.Show()
					go func() {
						// Load code
						encryptedBody, err := bima.Sync.GetBlob(codeEntry.Text)
						if err != nil {
							// show error
							s.Hide()
							s = nil
							dialog.ShowError(err, bima.UI.Window)
							return
						}
						var response SetupKit
						err = json.Unmarshal([]byte(encryptedBody), &response)

						if decodedText, err := base64.StdEncoding.DecodeString(response.AppID); err == nil {
							if decryptedText, err := shield.Decrypt(decodedText, bima.Registry.MasterPassword); err == nil {
								bima.Registry.AppID = string(decryptedText)
							}
						}

						if decodedText, err := base64.StdEncoding.DecodeString(response.SecretKey); err == nil {
							if decryptedText, err := shield.Decrypt(decodedText, bima.Registry.MasterPassword); err == nil {
								bima.Registry.SecretKey = decryptedText
							}
						}

						s.Hide()
						if bima.Registry.AppID != "" && len(bima.Registry.SecretKey) > 0 {
							// Update our registry and persist to db
							bima.Registry.Save()
							bima.Sync.AppID = bima.Registry.AppID
							DrawMainUI(bima)
						} else {
							// Show error
							dialog.ShowError(errors.New("Invalid password or setup code"), bima.UI.Window)
						}
					}()
				}),
			),
			layout.NewSpacer(),
		),
	}

	return &p
}
