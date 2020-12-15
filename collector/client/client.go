package client

import (
	"time"

	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

//CreateClient2 creates new http client
func CreateClient2() *browser.Browser {
	bow := surf.NewBrowser()
	bow.SetTimeout(60 * time.Second)
	bow.HistoryJar().SetMax(1)
	return bow
}
