package render

import (
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

type SecretKeyComponent struct {
	bima      *bima.Bima
	Container *fyne.Container
}

func (c *SecretKeyComponent) Remove() {
	return
}

func (c *SecretKeyComponent) Render() fyne.CanvasObject {
	return c.Container
}

func NewSecretKeyComponent(bima *bima.Bima) *SecretKeyComponent {
	s := SecretKeyComponent{
		bima: bima,
		Container: fyne.NewContainerWithLayout(layout.NewGridLayout(1),
			widget.NewLabel("Please save this key.\nBima combines this key together\nwith your master password\nto encrypt data.\nBima server has no access to this key and cannot help you to recover it."),
			&widget.Entry{
				Text:      bima.Registry.GetSetupKit(),
				MultiLine: true,
				Wrapping:  fyne.TextWrapBreak,
			},
			layout.NewSpacer()),
	}

	s.Container.AddObject(widget.NewButton("I saved the key securely!", func() {
		DrawMainUI(s.bima)
	}))

	return &s
}

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
		p.bima.Registry.SaveMasterPassword(p.passwordEntry.Text)
		DrawMainUI(p.bima)
	case NewPasswordForm:
		// Onboard form or enter password form
		log.Debug().Str("password", p.passwordEntry.Text).Msg("Change Password")
		p.bima.Registry.SaveMasterPassword(p.passwordEntry.Text)
		c := NewSecretKeyComponent(p.bima)
		p.bima.Push("show_secret_key", c)
	case EnterPasswordForm:
		if e := p.bima.Registry.SaveMasterPassword(p.passwordEntry.Text); e == nil {
			DrawMainUI(p.bima)
		} else {
			dialog.ShowInformation("Err", "Wrong password", p.bima.UI.Window)
		}
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

	passwordForm := widget.NewVBox(
		layout.NewSpacer(),
	)

	if formType == ChangePasswordForm || formType == NewPasswordForm {
		passwordForm.Append(widget.NewLabel("Pick a password to encrypt your data.\nIf you forgot this password,\nyour data is lost forever."))
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
