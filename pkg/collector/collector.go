package collector

import (
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/MikhailKlemin/gerzson.boros/pkg/boiler"
	"github.com/MikhailKlemin/gerzson.boros/pkg/browser"
	"github.com/MikhailKlemin/gerzson.boros/pkg/client"
	"github.com/MikhailKlemin/gerzson.boros/pkg/postprocess"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

//Entity holds all information
type Entity struct {
	MainDomain  string
	RedirectsTo string
	Links       []string
	Texts       []Resp
}

//Domain holds all scraped data
//This is internal struct
type Domain struct {
	subs    []string
	dlink   string
	cleanRe *regexp.Regexp
}

//Resp is holds information per link
type Resp struct {
	Link    string
	RawHTML string
	RawText string
	Boiler  string
}

//NewDomainCollector creates Collector Instances
func NewDomainCollector() *Domain {
	var d Domain
	d.cleanRe = regexp.MustCompile(`\s+`)
	d.subs = []string{"kapcsolat", "rolunk", "ceginformacio", "cegunkrol", "contact", "bemutatkozas", "elerhetoseg", "about", "elerhetosegeink", "elerhetosegek", "cegunkrol", "magunkrol", "contacts", "fooldal", "home", "szolgaltatasok", "index", "elerhetoseg", "rolam", "cegunkrol", "fooldal", "impresszum", "jogi-nyilatkozat", "cookie-szabalyzat", "adatkezelesi-szabalyzat", "szerzodesi-feltetelekszerzodesi-feltetelek", "feltetelekszerzodesi-feltetelek", "adatvedelmi-nyilatkozat", "adatkezelesi-tajekoztato", "adatvedelmi-tajekoztato", "adatvedelem", "adatkezeles", "nyilatkozat", "terms-and-conditions", "aszf", "privacy-policy"}
	return &d
}

//Start starts scraping
func (d *Domain) Start(dlink string) (data []Resp) {
	links, redirectedTo, err := d.collectLinks(dlink)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(redirectedTo)

	d.dlink = dlink
	sem := make(chan bool, 10)
	var mu sync.Mutex
	for _, link := range links {
		sem <- true
		go func(link string) {
			defer func() { <-sem }()
			//fmt.Println("Processing\t", link)
			d, bp, err := d.bpLinkWithChrome(link)
			if err != nil {
				log.Println(err)
				return
			}

			text, err := postprocess.Tokenize(strings.NewReader(d))
			if err != nil {
				log.Println(err)
				return
			}

			mu.Lock()
			data = append(data, Resp{link, d, text, bp})
			mu.Unlock()
		}(link)
		//fmt.Println(data)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return
}

//collectLinks parses 1st page without creating Chrome instance,
//and collects all links which are belogs to same host, and filter out
//links to pdf and doc files.
func (d *Domain) collectLinks(domainLink string) (links []string, redirectedTo string, err error) {
	mclient := client.CreateClient()

	um := make(map[string]bool)

	var (
		doc *goquery.Document
	)

	doc, redirectedTo, err = mclient.GetRedirect(domainLink)
	if err != nil {
		return links, "", errors.Wrap(err, "Cannot get domain link")
	}

	dURL, err := url.ParseRequestURI(redirectedTo)
	if err != nil {
		return links, "", errors.Wrap(err, "Cannot parse domain URL "+domainLink)
	}
	links = append(links, redirectedTo)

	//doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	doc.Find(`a`).Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr(`href`); ok {
			lURL, err := url.Parse(href)
			if err != nil {
				log.Println(errors.Wrap(err, "problem parsing URL"))
				return
			}
			link := dURL.ResolveReference(lURL).String()
			if d.contains(link) && d.samehost(domainLink, link) &&
				!strings.HasSuffix(link, ".pdf") &&
				!strings.HasSuffix(link, ".doc") {
				if _, ok := um[link]; !ok {
					um[link] = true
					links = append(links, link)

				}
			}
		}
	})

	return
}

func (d *Domain) bpLinkWithChrome(link string) (data string, bp string, err error) {
	data, err = browser.GetText(link)
	if err != nil {
		return data, bp, errors.Wrap(err, "can't get text for link:"+link)
	}
	data = d.cleanRe.ReplaceAllString(data, " ")
	bp, err = boiler.Getboiler(strings.NewReader(data))
	//bp, err = boiler.Tika(strings.NewReader(data))
	if err != nil {
		log.Println(err)
	}
	return

}

func (d *Domain) contains(link string) bool {

	for _, sub := range d.subs {
		if strings.Contains(link, sub) {
			return true
		}
	}
	return false

}

//samehost just comparing two link and figure if they are belong to same host
//in a way that http://a.mysite.com will be same host as http://b.mysite.com
func (d *Domain) samehost(dlink string, link string) bool {
	geth := func(link string) (string, error) {
		u, err := url.ParseRequestURI(link)
		if err != nil {
			//log.Fatal(err)
			return "", err
		}
		parts := strings.Split(u.Hostname(), ".")
		if len(parts) >= 2 {
			return parts[len(parts)-2] + "." + parts[len(parts)-1], nil
		}
		return "", errors.New("Not enought parts in URL\t" + link)

	}
	h1, err := geth(dlink)
	if err != nil {
		return false
	}

	h2, err := geth(link)
	if err != nil {
		return false
	}

	if h1 != h2 {
		return false
	}
	return true
}
