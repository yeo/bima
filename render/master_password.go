package render

import (
	"crypto/rand"

	"fyne.io/fyne"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/rs/zerolog/log"

	"github.com/yeo/bima/core"
)

type PasswordFormType int

const (
	NewPasswordForm    PasswordFormType = iota
	ChangePasswordForm PasswordFormType = iota
	EnterPasswordForm  PasswordFormType = iota
)

type PasswordComponent struct {
	formType             PasswordFormType
	bima                 *bima.Bima
	Container            *fyne.Container
	passwordEntry        *widget.Entry
	confirmPasswordEntry *widget.Entry
	actionButton         *widget.Button
}

func (p *PasswordComponent) Render() fyne.CanvasObject {
	return p.Container
}

func (p *PasswordComponent) Remove() {
	return
}

func (p *PasswordComponent) generateEncryptionKey() []byte {
	// Generate a 24 byte keys length
	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		// TODO: Show error to end user
		panic("Cannot generate key")
	}

	return b
}

func (p *PasswordComponent) Save() error {
	if p.formType == ChangePasswordForm || p.formType == NewPasswordForm {
		if (p.passwordEntry.Text == "") || (p.passwordEntry.Text != p.confirmPasswordEntry.Text) {
			dialog.ShowInformation("Password validation fail", "Both pass need to be same and not empty", p.bima.UI.Window)
			return nil
		}
	}

	switch p.formType {
	case ChangePasswordForm:
		log.Debug().Str("password", p.passwordEntry.Text).Msg("Enter New Password")
		p.generateEncryptionKey()
		p.bima.Registry.SaveMasterPassword(p.passwordEntry.Text)
		DrawMainUI(p.bima)
	case NewPasswordForm:
		// Onboard form or enter password form
		log.Debug().Str("password", p.passwordEntry.Text).Msg("Change Password")
		p.bima.Registry.SaveMasterPassword(p.passwordEntry.Text)
		DrawMainUI(p.bima)
	case EnterPasswordForm:
		p.bima.Registry.SaveMasterPassword(p.passwordEntry.Text)
		DrawMainUI(p.bima)
	}

	return nil
}

func NewPasswordComponent(bima *bima.Bima, formType PasswordFormType) *PasswordComponent {
	actionLabel := "Next"
	passwordLabel := "Enter Master Password"
	if formType == ChangePasswordForm {
		actionLabel = "Change Password"
		passwordLabel = "Enter New Password"
	}

	if formType == EnterPasswordForm {
		actionLabel = "Unlock"
	}

	p := PasswordComponent{
		formType:  formType,
		bima:      bima,
		Container: fyne.NewContainerWithLayout(layout.NewGridLayout(1)),
		passwordEntry: &widget.Entry{
			PlaceHolder: passwordLabel,
		},
		confirmPasswordEntry: &widget.Entry{
			PlaceHolder: "Confirm Password",
		},
	}
	p.actionButton = widget.NewButton(actionLabel, func() {
		p.Save()
	})

	passwordForm := widget.NewVBox()

	if formType == ChangePasswordForm || formType == NewPasswordForm {
		passwordForm.Append(widget.NewLabel("Pick a password to encrypt your data.\nIf you forgot this password,\nyour data is lost forever.\nMake sure it is at least 16 character"))
	} else {
		passwordForm.Append(widget.NewLabel("Enter password to decrypt your token\n"))
	}

	passwordForm.Append(p.passwordEntry)

	if formType == ChangePasswordForm || formType == NewPasswordForm {
		passwordForm.Append(p.confirmPasswordEntry)
	}

	passwordForm.Append(p.actionButton)
	passwordForm.Append(layout.NewSpacer())

	if p.formType == ChangePasswordForm {
		passwordForm.Append(
			widget.NewButton("Back", func() {
				DrawCode(bima)
			}))
	}

	p.Container.AddObject(passwordForm)
	p.Container.AddObject(layout.NewSpacer())

	return &p
}
