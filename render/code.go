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

func DrawViewCode(bima *bima.Bima, token *dto.Token) *widget.Button {
	button := widget.NewButton("View", func() {
		w := bima.UI.Window

		otpCode, _ := totp.GenerateCode(token.DecryptToken(bima.Registry.MasterPassword), time.Now())
		otpLbl := canvas.NewText(otpCode, color.RGBA{135, 0, 16, 255})
		refreshLbl := canvas.NewText("", color.RGBA{135, 0, 16, 255})
		otpLbl.TextSize = 20

		done := make(chan bool)
		go func() {
			secs := time.Now().Unix()
			remainder := secs % 30
			time.Sleep(time.Duration(30-remainder) * time.Second)
			ticker := time.NewTicker(30 * time.Second)
			for {
				select {
				case v := <-done:
					if v {
						log.Debug().Msg("Back to main screen")
						return
					}
				case <-ticker.C:
					otpCode, _ := totp.GenerateCode(token.DecryptToken(bima.Registry.MasterPassword), time.Now())
					otpLbl.Text = otpCode
					otpLbl.Refresh()
					log.Debug().Msg("Refresh token")
				}
			}
		}()

		var btn *widget.Button
		btn = widget.NewButton("Copy", func() {
			w.Clipboard().SetContent(otpCode)
			btn.SetText("Copied")
			btn.Style = widget.PrimaryButton
			timer2 := time.NewTimer(time.Second * 3)
			go func() {
				<-timer2.C
				if btn != nil {
					btn.SetText("Copy")
				}
			}()
		})

		actionButtons := widget.NewHBox(
			layout.NewSpacer(),
			btn,
			layout.NewSpacer(),
		)

		container := widget.NewHBox(
			widget.NewHBox(
				layout.NewSpacer(),
				canvas.NewText(token.URL, color.RGBA{135, 0, 16, 255}),
				layout.NewSpacer(),
			),
			fyne.NewContainerWithLayout(layout.NewGridLayout(1),
				widget.NewHBox(
					layout.NewSpacer(),
					canvas.NewText(token.Name, color.RGBA{135, 0, 16, 255}),
					layout.NewSpacer(),
				),
				widget.NewHBox(
					layout.NewSpacer(),
					otpLbl,
					layout.NewSpacer(),
				),
				widget.NewHBox(
					layout.NewSpacer(),
					refreshLbl,
					layout.NewSpacer(),
				),
				actionButtons,

				widget.NewHBox(
					layout.NewSpacer(),
					DrawEditCode(bima, token),
					widget.NewButton("Back", func() {
						done <- true
						bima.UI.Window.SetContent(bima.UI.MainContainer)
						DrawCode(bima)
					}),
					layout.NewSpacer(),
				),
				layout.NewSpacer(),
			))

		bima.UI.Window.SetContent(container)
		container.Refresh()
	})

	return button
}

func DrawCode(bima *bima.Bima) {
	bima.AppModel.CurrentScreen = "token/list"
	w := bima.UI.Window
	header := bima.UI.Header

	tokens, err := dto.LoadTokens()
	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(1))
	if err == nil {
		for _, token := range tokens {

			if bima.AppModel.FilterText != "" {
				if !strings.Contains(token.Name, bima.AppModel.FilterText) &&
					!strings.Contains(token.URL, bima.AppModel.FilterText) {
					log.Debug().Str("Filter", bima.AppModel.FilterText).Str("Name", token.Name).Msg("Not match filter. Skip")
					continue
				}
			}
			viewButton := DrawViewCode(bima, token)

			row :=
				fyne.NewContainerWithLayout(layout.NewGridLayout(1),
					widget.NewGroup(token.URL,
						widget.NewHBox(
							canvas.NewText(token.Name, color.RGBA{135, 0, 16, 255}),
							layout.NewSpacer(),
							viewButton,
						),
					),
				)

			codeContainer.AddObject(row)
		}
	}

	s := codeContainer
	c := fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		widget.NewScrollContainer(widget.NewVBox(header, s)))
	bima.UI.MainContainer = c
	w.SetContent(bima.UI.MainContainer)
}
