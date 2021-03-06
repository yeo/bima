package render

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/pquerna/otp/totp"
	"github.com/rs/zerolog/log"

	"github.com/yeo/bima/core"
	"github.com/yeo/bima/dto"
)

type CodeDetailComponent struct {
	bima      *bima.Bima
	Container fyne.CanvasObject
	token     *dto.Token
}

func (c *CodeDetailComponent) Render() fyne.CanvasObject {
	return c.Container
}

func (c *CodeDetailComponent) Remove() {
	//TODO: Remove timer, close channel
}

func NewCodeDetailComponent(bima *bima.Bima, tokenID string) *CodeDetailComponent {
	w := bima.UI.Window

	var token *dto.Token
	for _, v := range bima.AppModel.Tokens {
		if v.ID == tokenID {
			token = v
		}
	}

	urlLbl := canvas.NewText(token.URL, color.RGBA{0, 173, 181, 255})
	urlLbl.TextSize = 20
	nameLbl := canvas.NewText(token.Name, color.RGBA{54, 79, 107, 255})
	refreshLbl := canvas.NewText("", color.RGBA{57, 62, 70, 255})

	otpCode, _ := totp.GenerateCode(token.DecryptToken(bima.Registry.CombineEncryptionKey()), time.Now())
	otpLbl := canvas.NewText(otpCode, color.RGBA{252, 81, 133, 255})
	otpLbl.TextSize = 40

	done := make(chan bool)
	go func() {
		secs := time.Now().Unix()
		remainder := secs % 30
		secondToRefresh := 30 - remainder
		ticker := time.NewTicker(1 * time.Second)
		refreshLbl.Text = fmt.Sprintf("Regenerate in %2d s", secondToRefresh)
		refreshLbl.Refresh()
		for {
			select {
			case v := <-done:
				if v {
					log.Debug().Msg("Back to main screen")
					ticker.Stop()
					return
				}
			case <-ticker.C:
				secondToRefresh -= 1
				if secondToRefresh <= 0 {
					otpCode, _ = totp.GenerateCode(token.DecryptToken(bima.Registry.CombineEncryptionKey()), time.Now())
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

	var copyButton *widget.Button
	copyButton = widget.NewButton("Copy", func() {
		w.Clipboard().SetContent(otpCode)
		copyButton.SetText("Copied")
		//copyButton.Style = widget.PrimaryButton
		timer2 := time.NewTimer(time.Second * 3)
		go func() {
			<-timer2.C
			if copyButton != nil {
				copyButton.SetText("Copy")
			}
			timer2 = nil
		}()
	})

	backButton := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("Back", func() {
			done <- true
			log.Debug().Str("button", "code_detail.back").Msg("Click button")
			DrawCode(bima)
		}),
		layout.NewSpacer(),
	)

	container := fyne.NewContainerWithLayout(layout.NewGridLayout(1),
		container.NewHBox(
			layout.NewSpacer(), urlLbl, layout.NewSpacer(),
		),
		container.NewHBox(
			layout.NewSpacer(), nameLbl, layout.NewSpacer(),
		),
		container.NewHBox(
			layout.NewSpacer(), otpLbl, copyButton, layout.NewSpacer(),
		),
		container.NewHBox(
			layout.NewSpacer(), refreshLbl, layout.NewSpacer(),
		),
		container.NewHBox(
			layout.NewSpacer(),
			widget.NewButton("Delete", func() {
				dialog.ShowConfirm("Are you sure to delete?", "URL: "+token.URL+"\nName: "+token.Name, func(confirm bool) {
					if confirm == true {
						token.DeletedAt = time.Now().Unix()
						dto.DeleteSecret(token)
						DrawCode(bima)
					}
				}, bima.UI.Window)
			}),
			DrawEditCode(bima, token),
			layout.NewSpacer(),
		),
		layout.NewSpacer(),
		backButton,
	)

	return &CodeDetailComponent{
		bima:      bima,
		Container: container,
		token:     token,
	}
}

type ListCodeComponent struct {
	bima *bima.Bima
	// TODO: 2.0 fix
	codeContainer *fyne.Container
	Container     *fyne.Container

	done       chan (bool)
	ticker     *time.Ticker
	codeFilter string
}

func (c *ListCodeComponent) Render() fyne.CanvasObject {
	c.renderCode()

	return c.Container
}

func (c *ListCodeComponent) Remove() {
	c.done <- true
	c.ticker.Stop()
	return
}

func (c *ListCodeComponent) Refresh() {
	go func() {
		for {
			select {
			case <-c.done:
				return
			case <-c.ticker.C:
				if c.codeFilter != c.bima.AppModel.FilterText {
					c.renderCode()
				}
			}
		}
	}()
}

func (c *ListCodeComponent) renderCode() {
	bima := c.bima
	tokens := c.bima.AppModel.Tokens

	//codeContainer := widget.NewGroupWithScroller("Tokens")
	// TODO: 2.0 fix
	codeContainer := container.NewVBox()

	c.codeFilter = strings.Trim(bima.AppModel.FilterText, " ")

	for _, token := range tokens {
		if c.codeFilter != "" {
			if !strings.Contains(token.Name, bima.AppModel.FilterText) &&
				!strings.Contains(token.URL, bima.AppModel.FilterText) {
				continue
			}
		}

		viewButton := widget.NewButton("View", func(t *dto.Token) func() {
			return func() {
				c := NewCodeDetailComponent(bima, t.ID)
				bima.Push("token/view", c)
			}
		}(token))

		urlLbl := canvas.NewText(token.URL, color.RGBA{0, 173, 181, 255})
		urlLbl.TextSize = 17

		nameLbl := canvas.NewText(token.Name, color.RGBA{54, 79, 107, 255})
		row := container.NewVBox(
			container.NewHBox(urlLbl),
			container.NewHBox(
				nameLbl,
				layout.NewSpacer(),
				viewButton,
			),
			layout.NewSpacer(),
			canvas.NewLine(color.RGBA{34, 40, 49, 50}),
		)

		codeContainer.Add(row)
	}

	//lastRow := fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(320, 560)), codeContainer)
	lastRow := fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(320, 560)), container.NewVScroll(codeContainer))

	// TODO: See if we can avoid setting to nil and check memory leak
	c.Container.Objects[1] = nil
	c.Container.Objects[1] = lastRow
	c.Container.Refresh()
}

func NewListCodeComponent(bima *bima.Bima) *ListCodeComponent {
	header := bima.UI.Header

	if tokens, err := dto.LoadTokens(); err == nil {
		bima.AppModel.Tokens = tokens
	} else {
		bima.AppModel.Tokens = []*dto.Token{}
	}

	//codeContainer := widget.NewGroupWithScroller("Tokens")
	// TODO: 2.0 fix
	codeContainer := container.NewVBox()

	c := fyne.NewContainerWithLayout(layout.NewVBoxLayout(),
		header,
		container.NewVScroll(codeContainer))
	//fyne.NewContainerWithLayout(layout.NewGridWrapLayout(fyne.NewSize(320, 560)), container.NewVScroll(codeContainer)))

	p := &ListCodeComponent{
		bima:          bima,
		Container:     c,
		codeContainer: codeContainer,
		ticker:        time.NewTicker(500 * time.Millisecond),
		done:          make(chan bool),
	}

	p.renderCode()
	p.Refresh()
	return p
}

func DrawCode(bima *bima.Bima) {
	c := NewListCodeComponent(bima)

	bima.Push("token/list", c)
}
