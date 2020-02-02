package render

import (
	"image/color"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"github.com/pquerna/otp/totp"

	"github.com/yeo/bima/dto"
)

func DrawCode(w fyne.Window, header *widget.Box) {
	tokens, err := dto.LoadTokens()

	codeContainer := fyne.NewContainerWithLayout(layout.NewGridLayout(2))
	if err == nil {
		for _, token := range tokens {
			otpCode, _ := totp.GenerateCode(token.Token, time.Now())

			t := canvas.NewText(otpCode, color.RGBA{10, 200, 200, 0})
			codeContainer.AddObject(t)
		}
	}

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(2), header, codeContainer)
	w.SetContent(container)
}
