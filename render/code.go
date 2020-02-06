package render

import (
	"image/color"
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/pquerna/otp/totp"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/dto"
)

func DrawCode(bima *bima.Bima) {
	w := bima.UI.Window
	header := bima.UI.Header

	tokens, err := dto.LoadTokens()

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(4))
	if err == nil {
		for _, token := range tokens {
			otpCode, _ := totp.GenerateCode(token.DecryptToken(bima.Registry.MasterPassword), time.Now())
			log.Println("Render for", token.Name)

			lbl := canvas.NewText(token.Name, color.RGBA{38, 41, 45, 0})
			codeContainer.AddObject(lbl)
			t := canvas.NewText(otpCode, color.RGBA{10, 200, 200, 0})
			codeContainer.AddObject(t)

			var btn *widget.Button
			btn = widget.NewButton("Copy", func() {
				w.Clipboard().SetContent(otpCode)
				btn.SetText("Copied √")
				timer2 := time.NewTimer(time.Second * 10)
				go func() {
					<-timer2.C
					if btn != nil {
						btn.SetText("Copy")
					}
				}()
			})

			codeContainer.AddObject(btn)

			editButton := DrawEditCode(bima, token)
			codeContainer.AddObject(editButton)
		}
	}

	bima.UI.MainContainer = fyne.NewContainerWithLayout(layout.NewGridLayout(1), header, codeContainer)
	w.SetContent(bima.UI.MainContainer)
}
