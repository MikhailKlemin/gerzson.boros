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

	req.Header.Add("DNT", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:83.0) Gecko/20100101 Firefox/83.0")
	req.Header.Set("Accept", "image/webp,*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Referer", "http://www.google.com/")
	req.Header.Set("Cache-Control", "max-age=0")

	//req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Referer", "https://www.google.com/")
	//req.Header.Add("Cookie", "session-id=139-2100433-8588745; session-id-time=2082787201l; csm-hit=tb:DPNXSRSRNMPD1PSXG3A1+s-WVFNSJTDXHV9CSQ4REGA|1605104273946&t:1605104273946&adb:adblk_no; ubid-main=132-9103469-6245957; sst-main=Sst1|PQHSea3rFDAFE5GC7N_iQr1ECw117L2m3E5Oj3PMkWm9uxuM6nf_2OGqqsTIKkwXbe0tKyN--dZir_iDlfx1GTLHQnZMk4y8SmbVKIqcre41VqebgPVjnXCJP7_CKSL1NA-cHJnWHbMvEQEVH0mZ5JxczXRAQMlWNnMIHu-fjBTg3X26cqX9Y8UYDmaYbQPW0iywDGqboJcHJtBo3pyksgcui2njLdJIxagTbWUzBWdnU58T9WBU4EdUpZqbIAxBQ5Q4gUXm_DkvDn7Jsst93po7uY7LZP44fBbXiprYRQJbNCAE5-zXKxQcE05tSok72qRjYbKXvz5PUpsrwZ3Ge8G_6g; i18n-prefs=USD; aws-ubid-main=263-2666761-6844736; sid=\"wo19MEEyl2Jo6fEbfaqOug==|y8YB0FJFlo8Fm48PE203C38wB/UWwKr41sbmLZFxG5w=\"; session-token=0P5D8A/2mVo88GDBnwrjqkafoEsjrY7E5MTF1Yn0EhD2diA19l71sRqnsam/mT8kZoVUp2yjRikBCbFkgQFjLPwPGpezpeOtgZBZ9/4UbS96jvEoGpXXtCKWYhseb7DvQOsPJF8LN+eCblWgYjI7qn3xTSEI6o0wcp28RZ6Kg+qL3xM/dxtafZ7DZdrORzXRHCuc5sdzfsVZ2rfWC+3U8e7OdtN7Ae0hKhwPq27itg+nps7OmEIWNH2AGfwO4OjM; aws-priv=eyJ2IjoxLCJldSI6MCwic3QiOjB9; aws-target-static-id=1592234761606-26391; aws-target-visitor-id=1592234761609-999013.38_0; aws-target-data=%7B%22support%22%3A%221%22%7D; s_fid=193D473662E77E28-2C4E7F48C53AA9B2; s_dslv=1596065479166; s_vn=1623770761766%26vn%3D4; regStatus=pre-register; s_vnum=2024644422417%26vn%3D1; skin=noskin; sp-cdn=\"L5Z9:RU\"; ubid-main=132-9103469-6245957; session-id-time=2082787201l; session-id=139-2100433-8588745")
	req.Header.Add("TE", "Trailers")

	for {
		if counter > 1 {
			return b, fmt.Errorf("%s: %w", link, ErrRetry)
		}
		counter++
		//resp, err := m.client.Get(link)
		resp, err := m.client.Do(req)
		if err != nil {
			log.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}
		if resp.StatusCode != 200 {
			log.Println(resp.Status)
		}

		b, err = ioutil.ReadAll(resp.Body)
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
