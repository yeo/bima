package render

import (
	"fmt"
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

		urlLbl := canvas.NewText(token.URL, color.RGBA{135, 0, 16, 255})
		urlLbl.TextSize = 20
		nameLbl := canvas.NewText(token.Name, color.RGBA{135, 0, 16, 255})
		refreshLbl := canvas.NewText("", color.RGBA{135, 0, 16, 255})

		otpCode, _ := totp.GenerateCode(token.DecryptToken(bima.Registry.MasterPassword), time.Now())
		otpLbl := canvas.NewText(otpCode, color.RGBA{135, 0, 16, 255})
		otpLbl.TextSize = 40

		done := make(chan bool)
		go func() {
			secs := time.Now().Unix()
			remainder := secs % 30
			//time.Sleep(time.Duration(30-remainder) * time.Second)
			secondToRefresh := 30 - remainder
			ticker := time.NewTicker(1 * time.Second)
			refreshLbl.Text = fmt.Sprintf("Regenerate in %2d s", secondToRefresh)
			refreshLbl.Refresh()
			for {
				select {
				case v := <-done:
					if v {
						log.Debug().Msg("Back to main screen")
						return
					}
				case <-ticker.C:
					secondToRefresh -= 1
					if secondToRefresh <= 0 {
						otpCode, _ = totp.GenerateCode(token.DecryptToken(bima.Registry.MasterPassword), time.Now())
						otpLbl.Text = otpCode
						otpLbl.Refresh()
						secondToRefresh = 30
						log.Debug().Str("url", token.URL).Msg("Re-generator otp token")
					}
					refreshLbl.Text = fmt.Sprintf("Regenerate in %d s", secondToRefresh)
					refreshLbl.Refresh()
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

		container := fyne.NewContainerWithLayout(layout.NewGridLayout(1),
			widget.NewHBox(
				layout.NewSpacer(), urlLbl, layout.NewSpacer(),
			),
			widget.NewHBox(
				layout.NewSpacer(), nameLbl, layout.NewSpacer(),
			),
			widget.NewHBox(
				layout.NewSpacer(), otpLbl, layout.NewSpacer(),
			),
			widget.NewHBox(
				layout.NewSpacer(), refreshLbl, layout.NewSpacer(),
			),
			actionButtons,
			layout.NewSpacer(),

			widget.NewHBox(
				layout.NewSpacer(),
				DrawEditCode(bima, token),
				widget.NewButton("Back", func() {
					done <- true
					log.Debug().Str("button", "code_detail.back").Msg("Click button")
					DrawCode(bima)
				}),
				layout.NewSpacer(),
			),
			layout.NewSpacer(),
		)

		bima.UI.Window.SetContent(container)
		container.Refresh()
	})

	return button
}

type ListCodeComponent struct {
	bima          *bima.Bima
	tokens        []*dto.Token
	codeContainer *widget.Group
	Container     fyne.CanvasObject
}

func (c *ListCodeComponent) Render() fyne.CanvasObject {
	return c.Container
}

func (c *ListCodeComponent) Remove() {
	return
}

func NewListCodeComponent(bima *bima.Bima) *ListCodeComponent {
	bima.AppModel.CurrentScreen = "token/list"
	header := bima.UI.Header

	tokens, err := dto.LoadTokens()
	codeContainer := widget.NewGroupWithScroller("Tokens")
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
				widget.NewVBox(
					widget.NewHBox(
						layout.NewSpacer(),
						canvas.NewText(token.URL, color.RGBA{135, 0, 16, 255}),
						layout.NewSpacer(),
					),

					widget.NewHBox(
						layout.NewSpacer(),
						canvas.NewText(token.Name, color.RGBA{135, 0, 16, 255}),
						layout.NewSpacer(),
						viewButton,
					),
					layout.NewSpacer(),
				)

			codeContainer.Append(row)
		}
	}

	tokenList := fyne.NewContainerWithLayout(layout.NewFixedGridLayout(fyne.NewSize(320, 560)), codeContainer)
	c := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		header,
		tokenList)

	return &ListCodeComponent{
		bima:          bima,
		tokens:        tokens,
		Container:     c,
		codeContainer: codeContainer,
	}
}

func DrawCode(bima *bima.Bima) {
	c := NewListCodeComponent(bima)

	bima.Push("token/list", c)
}
