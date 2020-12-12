package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/headzoo/surf"
	"github.com/headzoo/surf/browser"
)

//MyClient is http client
type MyClient struct {
	client *http.Client
}

//ErrRetry - when client rich maximum  retries
var ErrRetry = errors.New("Maxium retry reached")

//CreateClient creates new http client
func CreateClient() *MyClient {
	var m MyClient
	var netTransport = &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 120 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 120 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	//netTransport.TLSClientConfig =

	var netClient = &http.Client{
		Timeout:   time.Second * 120,
		Transport: netTransport,
	}
	m.client = netClient
	m.client = &http.Client{}
	return &m
}

//CreateClient2 creates new http client
func CreateClient2() *browser.Browser {
	bow := surf.NewBrowser()
	bow.HistoryJar().SetMax(1)
	return bow
}

//Get is Get request
func (m *MyClient) Get(link string) (doc *goquery.Document, err error) {
	counter := 0
	for {
		if counter > 1 {
			return doc, fmt.Errorf("%s: %w", link, ErrRetry)
		}
		counter++
		resp, err := m.client.Get(link)
		if err != nil {
			log.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}
		if resp.StatusCode != 200 {
			log.Println(resp.Status)
		}

		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println(err)

			time.Sleep(10 * time.Second)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		break
	}
	return

}

//GetByte is Get request
func (m *MyClient) GetByte(link string) (b []byte, err error) {
	counter := 0

	req, err := http.NewRequest("GET",
		link,
		nil)

	if err != nil {
		//return nil, errors.Wrap(err, "failed to create request instance")
		log.Fatal(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("DNT", "1")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("Cache-Control", "max-age=0")
	//	req.Header.Add("Cookie", "mobile_detect=desktop; inx_checker2=1; INX_CHECKER2=1; PHPSESSID=ghfdeejhkf9jdrhah6st5pd6jq")

	for {
		if counter > 1 {
			return b, fmt.Errorf("%s: %w", link, ErrRetry)
		}
		counter++
		//resp, err := m.client.Get(link)
		resp, err := m.client.Do(req)
		if err != nil {
			log.Println(link, "\t:\t", err)
			time.Sleep(10 * time.Second)
			continue
		}
		if resp.StatusCode != 200 {

			log.Println(link, "\t:\t", resp.Status)
			time.Sleep(2 * time.Second)
			continue

		}

		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			//log.Println(err)
			time.Sleep(10 * time.Second)
			resp.Body.Close()
			continue
		}
		resp.Body.Close()
		break
	}
	return

}

//GetByte2 is Get request
func (m *MyClient) GetByte2(link string) (b []byte, err error) {
	bow := surf.NewBrowser()
	bow.HistoryJar().SetMax(1)

	err = bow.Open(link)
	if err != nil {
		return
	}

	body, err := bow.Dom().Html()
	if err != nil {
		return
	}

	return []byte(body), nil

}

//GetRedirect is Get request
func (m *MyClient) GetRedirect(link string) (doc *goquery.Document, redirectedTo string, err error) {
	counter := 0
	for {
		if counter > 1 {
			return doc, "", fmt.Errorf("%s: %w", link, ErrRetry)
		}
		counter++
		resp, err := m.client.Get(link)
		if err != nil {
			counter++
			//log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			counter++
			time.Sleep(2 * time.Second)
			//log.Println(resp.Status)
		}
		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			//log.Println(err)
			time.Sleep(10 * time.Second)
			resp.Body.Close()
			counter++
			continue
		}
		resp.Body.Close()
		redirectedTo = resp.Request.URL.String()
		break
	}
	return

}
