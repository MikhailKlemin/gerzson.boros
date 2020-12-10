package client

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
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
	}

	netTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	var netClient = &http.Client{
		Timeout:   time.Second * 120,
		Transport: netTransport,
	}
	m.client = netClient
	return &m
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
			log.Println(err)
			time.Sleep(2 * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			counter++
			time.Sleep(2 * time.Second)
			log.Println(resp.Status)
		}
		doc, err = goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println(err)
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
