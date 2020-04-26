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

func (p *PasswordComponent) Save() error {
	log.Debug().Str("password", p.passwordEntry.Text).Msg("Save New Password")

	if (p.passwordEntry.Text == "") || (p.passwordEntry.Text != p.confirmPasswordEntry.Text) {
		dialog.ShowInformation("Password validation fail", "Both pass need to be same and not empty", p.bima.UI.Window)
		return nil
	}

	p.bima.Registry.SaveMasterPassword(p.passwordEntry.Text)
	return nil
}

func NewPasswordComponent(bima *bima.Bima, formType PasswordFormType) *PasswordComponent {
	actionLabel := "Save"
	passwordLabel := "Enter password"
	if formType == ChangePasswordForm {
		actionLabel = "Change Password"
		passwordLabel = "Enter New Password"
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
		widget.NewLabel("Pick a password to encrypt your data"),
		p.passwordEntry,
		p.confirmPasswordEntry,
		p.actionButton,
		layout.NewSpacer(),
		widget.NewButton("Back", func() {
			DrawCode(bima)
		}),
	)
	p.Container.AddObject(passwordForm)
	p.Container.AddObject(layout.NewSpacer())

	return &p
}

func DrawMasterPassword(bima *bima.Bima, done func(*bima.Bima)) {
	bima.AppModel.CurrentScreen = "master_password"

	passwordEntry := &widget.Entry{
		PlaceHolder: "Enter Master Password",
	}

	passwordField := widget.NewButton("Unlock", func() {
		bima.Registry.SaveMasterPassword(passwordEntry.Text)
		done(bima)
	})
	passwordForm := widget.NewVBox(
		layout.NewSpacer(),
		passwordEntry, passwordField,
		layout.NewSpacer(),
	)

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	container.AddObject(layout.NewSpacer())
	container.AddObject(passwordForm)
	container.AddObject(layout.NewSpacer())

	bima.UI.Window.SetContent(container)
}
