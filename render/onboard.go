package render

//
//import (
//	"image/color"
//	"time"
//
//	"fyne.io/fyne"
//	"fyne.io/fyne/canvas"
//	"fyne.io/fyne/layout"
//	"fyne.io/fyne/widget"
//	"github.com/pquerna/otp/totp"
//
//	"github.com/yeo/bima/core"
//	"github.com/yeo/bima/dto"
//)
//
//type Onboard struct {
//	Bima *bima.Bima
//
//	ScratchButton *widget.Button
//	MigrateButton *widget.Button
//
//	container fyne.Container
//}
//
//func NewOnboard(bima *bima.App) *Onboad {
//	o := Onboard{bima: bima}
//	o.ScratchButton = widget.NewButton("Starting from scrach", func() {
//	})
//
//	o.MigrateButton = widget.NewButton("Already use Bima? Sync from other", func() {
//	})
//
//	o.container = fyne.NewContainerWithLayout(layout.NewGridLayout(1))
//
//	container.AddObject(layout.NewSpacer())
//	container.AddObject(o.ScratchButton)
//	container.AddObject(o.MigrateButton)
//	container.AddObject(layout.NewSpacer())
//
//	return &o
//}
//
//func (o *Onboard) Render() {
//	o.Bima.UI.Window.SetContent(container)
//}
