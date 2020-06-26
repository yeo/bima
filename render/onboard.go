package render

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	//"fmt"
	"image/color"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	//"github.com/rs/zerolog/log"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/shield"
)

const (
	SetupTypeQuickCode = "QuickCode"
	SetupTypeSetupKit  = "SetupKit"
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
			widget.NewVBox(
				layout.NewSpacer(),
				widget.NewButton("Starting from scratch", func() {
					c := NewPasswordComponent(bima, NewPasswordForm)
					bima.Push("onboard", c)
				}),
				&widget.Label{
					Text:     "Use this when you do not have any existed installation to sync data from",
					Wrapping: fyne.TextWrapWord,
				},
				layout.NewSpacer(),
				canvas.NewLine(color.Black),
				layout.NewSpacer(),
				&widget.Label{
					Text:     "If you have data on other devices, you can restore from other app using Setup Code generate from that app. Of if you no longer has access to any of app installation, you can use restore from Setup Kit, which Bima prompted you to save it.",
					Wrapping: fyne.TextWrapWord,
				},

				widget.NewButton("Use other Bima app", func() {
					c := NewCodeSetupComponent(bima, SetupTypeQuickCode)
					bima.Push("setup/code", c)
				}),
				widget.NewButton("Restore with Setup Kit", func() {
					c := NewCodeSetupComponent(bima, SetupTypeSetupKit)
					bima.Push("setup/kit", c)
				}),
				layout.NewSpacer(),
			),
		),
	}

	return &p
}

type CodeSetupComponent struct {
	bima      *bima.Bima
	Container *fyne.Container
	setupType string
}

func (p *CodeSetupComponent) Render() fyne.CanvasObject {
	return p.Container
}

func (p *CodeSetupComponent) Remove() {
	return
}

func NewCodeSetupComponent(bima *bima.Bima, setupType string) *CodeSetupComponent {
	codeEntry := widget.NewEntry()
	masterPassword := &widget.Entry{
		Password: true,
	}

	label := "On the other app, go to\nSetting > Generate Setup Code"
	if setupType == SetupTypeSetupKit {
		label = "Paste your setup kit if you have.\nIt can also found in Setting > Show Setup Kit"
	}

	p := CodeSetupComponent{
		bima: bima,
		Container: fyne.NewContainerWithLayout(layout.NewGridLayout(1),
			layout.NewSpacer(),
			widget.NewVBox(
				widget.NewLabel(label),
				codeEntry,
				layout.NewSpacer(),
				widget.NewLabel("Enter your master password"),
				masterPassword,
				widget.NewButton("Next", func() {
					bima.Registry.MasterPassword = []byte(masterPassword.Text)

					s := dialog.NewProgressInfinite("Syncing data", "...", bima.UI.Window)
					s.Show()
					go func() {
						// Load code
						var encryptedBody string
						var err error

						if setupType == SetupTypeSetupKit {
							encryptedBody = codeEntry.Text
						} else {
							encryptedBody, err = bima.Sync.GetBlob(codeEntry.Text)
						}

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
							go bima.Sync.ResumeSync()
							DrawMainUI(bima)
						} else {
							// Show error
							dialog.ShowError(errors.New("Invalid password or setup code"), bima.UI.Window)
						}
					}()
				}),
				widget.NewButton("Back", func() {
					c := NewOnboardComponent(bima)
					bima.Push("onboard", c)
				}),
			),
			layout.NewSpacer(),
		),
	}

	return &p
}
