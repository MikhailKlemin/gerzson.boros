package browser

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

/*
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
		Headless(true).
		Set("blink-settings", "imagesEnabled=false").
		Devtools(false)

	defer l.Cleanup() // remove user-data-dir

	url := l.MustLaunch()

	browser := rod.New().
		ControlURL(url).
		Trace(true).
		SlowMotion(1 * time.Second).
		MustConnect()

	defer browser.Close()
	var page *rod.Page
	err := rod.Try(func() {
		page = browser.Timeout(time.Duration(timeout) * time.Second).MustPage(link)

	})

	check(err)
	if err != nil {
		return "", errors.Wrap(err, "an error")
	}


	var html string
	err = rod.Try(func() {
		html = page.MustSearch(`body`).MustHTML()

	})
	check(err)
	if err != nil {
		return "", errors.Wrap(err, "can't search page")
	}

	return html, nil
}
*/
