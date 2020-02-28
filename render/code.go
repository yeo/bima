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
	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(3))
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

			codeContainer.AddObject(widget.NewVBox(
				canvas.NewText(token.Name, color.RGBA{38, 41, 45, 0}),
				canvas.NewText(token.URL, color.RGBA{38, 41, 45, 0}),
			))
			t := canvas.NewText(otpCode, color.RGBA{10, 200, 200, 0})
			var btn *widget.Button
			btn = widget.NewButton("Copy", func() {
				w.Clipboard().SetContent(otpCode)
				btn.SetText("Copied")
				timer2 := time.NewTimer(time.Second * 10)
				go func() {
					<-timer2.C
					if btn != nil {
						btn.SetText("Copy")
					}
				}()
			})

			codeContainer.AddObject(widget.NewVBox(
				t,
				btn,
			))

			editButton := DrawEditCode(bima, token)
			codeContainer.AddObject(widget.NewVBox(
				editButton,
			))
		}
	}

	//s := widget.NewScrollContainer(codeContainer)
	s := &widget.ScrollContainer{
		Content: codeContainer,
		Offset:  fyne.Position{0, 400},
	}

	s.Resize(fyne.Size{300, 400})
	c := fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		widget.NewVBox(header, s))
	bima.UI.MainContainer = c
	w.SetContent(bima.UI.MainContainer)
}
