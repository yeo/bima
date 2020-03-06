package render

import (
	//"time"

	"github.com/yeo/bima/core"
)

type OtpRefresher struct {
	OtpSecret []string
}

func (r *OtpRefresher) Start() {
	//go func() {
	//	secs := time.Now().Unix()
	//	remainder := secs % 30
	//	time.Sleep(time.Duration(30-remainder) * time.Second)
	//	ticker := time.NewTicker(30 * time.Second)
	//	for {
	//		select {
	//		case <-ticker.C:
	//			if bima.AppModel.CurrentScreen == "token/list" {
	//				DrawCode(bima)
	//			}
	//		}
	//	}
	//}()
}
func (r *OtpRefresher) Stop() {
}

func DrawMainUI(bima *bima.Bima) {
	bima.AppModel.CurrentScreen = "token/list"
	DrawCode(bima)
}
