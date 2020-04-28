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
}

type HeaderComponent struct {
	bima      *bima.Bima
	Container *fyne.Container
}

func (h *HeaderComponent) Render() *fyne.Container {
	return h.Container
}

func (h *HeaderComponent) Remove() {
	return
}

func NewHeaderComponent(bima *bima.Bima) *HeaderComponent {
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
	container := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.Size{300, 50}), headerWidget)

	h := HeaderComponent{
		bima:      bima,
		Container: container,
	}

	return &h
}

func DrawMainUI(bima *bima.Bima) {
	c := NewListCodeComponent(bima)
	bima.Push("token/list", c)
}

// main Entrypoint to render whole ui
func Render(bima *bima.Bima) {
	h := NewHeaderComponent(bima)
	bima.UI.Header = h.Render()

	// If never see onboard yet, we should show up onboard screen to enter email and setup password
	if bima.Registry.HasSetPassword == "" {
		c := NewPasswordComponent(bima, NewPasswordForm)
		bima.Push("onboard", c)
	} else {
		// To avoid the annoying of password when debugging, we support set password via env.
		if password := os.Getenv("BIMAPASS"); password != "" {
			bima.Registry.MasterPassword = password
		}

		if bima.Registry.MasterPassword == "" {
			c := NewPasswordComponent(bima, EnterPasswordForm)
			bima.Push("unlock", c)
		} else {
			DrawMainUI(bima)
			//Check if sync is enable or not. If yes, we will start syncing
		}
	}

	bima.UI.Window.Resize(fyne.NewSize(320, 640))
	bima.UI.Window.ShowAndRun()
}
