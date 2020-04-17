package render

import (
	//"time"

	"github.com/yeo/bima/core"
)

func DrawMainUI(bima *bima.Bima) {
	bima.AppModel.CurrentScreen = "token/list"
	DrawCode(bima)
}
