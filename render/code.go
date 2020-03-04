package render

import (
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog/log"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/dto"
)

func DrawCode(bima *bima.Bima) {
	bima.AppModel.CurrentScreen = "token/list"
	w := bima.UI.Window
	header := bima.UI.Header

	tokens, err := dto.LoadTokens()
	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	if err == nil {
		for _, token := range tokens {
			otpCode, _ := totp.GenerateCode(token.DecryptToken(bima.Registry.MasterPassword), time.Now())

			if bima.AppModel.FilterText != "" {
				if !strings.Contains(token.Name, bima.AppModel.FilterText) &&
					!strings.Contains(token.URL, bima.AppModel.FilterText) {
					log.Debug().Str("Filter", bima.AppModel.FilterText).Str("Name", token.Name).Msg("Not match filter. Skip")
					continue
				}
			}

			otpLbl := canvas.NewText(otpCode, color.RGBA{0, 0, 0, 0})
			var btn *widget.Button
			btn = widget.NewButton("Copy", func() {
				w.Clipboard().SetContent(otpCode)
				btn.SetText("Copied")
				btn.Style = widget.PrimaryButton
				timer2 := time.NewTimer(time.Second * 10)
				go func() {
					<-timer2.C
					if btn != nil {
						btn.SetText("Copy")
					}
				}()
			})

			editButton := DrawEditCode(bima, token)
			row := fyne.NewContainerWithLayout(layout.NewGridLayout(3),
				widget.NewVBox(
					canvas.NewText(token.Name, color.RGBA{10, 10, 10, 0}),
					canvas.NewText(token.URL, color.RGBA{10, 10, 10, 0}),
				), widget.NewVBox(
					otpLbl,
					btn,
				), widget.NewVBox(
					editButton,
				))

			codeContainer.AddObject(row)
			codeContainer.AddObject(widget.NewVBox(canvas.NewLine(color.RGBA{50, 50, 50, 0})))
		}
	}

	s := codeContainer
	//s.Resize(fyne.Size{500, 500})
	c := fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		widget.NewScrollContainer(widget.NewVBox(header, s)))
	bima.UI.MainContainer = c
	w.SetContent(bima.UI.MainContainer)
}
