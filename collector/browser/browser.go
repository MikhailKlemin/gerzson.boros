package browser

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
)

//GetText starts ChromeInstance and get the link
func GetText(link string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()
	// run task list
	var res string
	if err := chromedp.Run(ctx,
		chromedp.Navigate(link),
		chromedp.WaitVisible(`body`, chromedp.ByQuery),
		chromedp.OuterHTML(`html`, &res),
	); err != nil {
		return res, errors.Wrap(err, "Chrome timed out")
	}

	return res, nil
}
