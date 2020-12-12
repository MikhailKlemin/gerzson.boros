package browser

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/pkg/errors"
)

/*
//GetText starts ChromeInstance and get the link
func GetText(link string) (string, error) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("ignore-certificate-errors", "1"),
	)

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel = context.WithTimeout(context.Background(), 120*time.Second)
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
*/

//GetText2 Same as GetText but with rod
func GetText2(link string, timeout int) (string, error) {

	check := func(err error) {
		var evalErr *rod.ErrEval
		if errors.Is(err, context.DeadlineExceeded) { // timeout error
			fmt.Println("timeout err")
			//return "", errors.Wrap(err, "Time out")
		} else if errors.As(err, &evalErr) { // eval error
			fmt.Println(evalErr.LineNumber)
		} else if err != nil {
			fmt.Println("can't handle", err)
		}
	}

	l := launcher.New().
		//		Set("proxy-server", "socks5://"+p). // add a flag, here we set a http proxy
		Headless(true).
		Set("blink-settings", "imagesEnabled=false").
		Devtools(false)

	defer l.Cleanup() // remove user-data-dir
	//l.ProfileDir("/media/mike/WDC4_1/chrome-profiles/" + p)

	url := l.MustLaunch()

	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(1 * time.Second).
		MustConnect()

	// auth the proxy
	// here we use cli tool "mitmproxy --proxyauth user:pass" as an example
	defer browser.Close()
	//defer time.Sleep(5 * time.Second)
	//page := browser.MustPage(link)
	var page *rod.Page
	err := rod.Try(func() {
		page = browser.Timeout(time.Duration(timeout) * time.Second).MustPage(link)

	})

	check(err)
	if err != nil {
		return "", errors.Wrap(err, "an error")
	}

	/*	if err != nil {
			return "", errors.Wrap(err, fmt.Sprintf("Cannot navigate to %s", link))
		}
	*/
	/*
		elems, err := page.Elements("frame")
		if err != nil {
			return "", nil
		}

		var src []string
		for _, elem := range elems {
			src = append(src, elem.MustHTML())

		}

		if len(elems) > 0 {
			return strings.Join(src, "\n"), nil

		}
	*/
	var html string
	err = rod.Try(func() {
		html = page.MustSearch(`body`).MustHTML()

	})
	check(err)
	if err != nil {
		return "", errors.Wrap(err, "can't search page")
	}

	return html, nil
	//return "", err
}
