package collector

import (
	"log"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/MikhailKlemin/gerzson.boros/collector/boiler"
	"github.com/MikhailKlemin/gerzson.boros/collector/client"
	"github.com/MikhailKlemin/gerzson.boros/collector/config"
	"github.com/MikhailKlemin/gerzson.boros/collector/postprocess"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

//Collector is
type Collector struct {
	Opts Options
	c    config.GeneralConfig
}

//Entity is
type Entity struct {
	MainDomain  string   `bson:"maindomain"`
	RedirectsTo string   `bson:"redirectsto"`
	Links       []string `bson:"links"`
	Texts       []Info   `bson:"texts"`
}

//Info scraped data from domain
type Info struct {
	Link      string    `bson:"link"`
	RawHTML   string    `bson:"rawhtml"`
	RawText   string    `bson:"rawtest"`
	Boiler    string    `bson:"boiler"`
	TimeStamp time.Time `bson:"timestamp"`
}

//Options have default setting
type Options struct {
	Domain   string
	Keywords []string
	re       *regexp.Regexp
}

//NewCollector constructor
func NewCollector(Domain string, conf config.GeneralConfig) *Collector {
	var c Collector
	c.Opts = DefaultOptions(Domain)
	c.c = conf
	return &c

}

//DefaultOptions constructor
func DefaultOptions(Domain string) Options {
	var o Options
	o.Domain = Domain
	o.Keywords = []string{"kapcsolat", "rolunk", "ceginformacio", "cegunkrol", "contact", "bemutatkozas", "elerhetoseg", "about", "elerhetosegeink", "elerhetosegek", "cegunkrol", "magunkrol", "contacts", "fooldal", "home", "szolgaltatasok", "index", "elerhetoseg", "rolam", "cegunkrol", "fooldal", "impresszum", "jogi-nyilatkozat", "cookie-szabalyzat", "adatkezelesi-szabalyzat", "szerzodesi-feltetelekszerzodesi-feltetelek", "feltetelekszerzodesi-feltetelek", "adatvedelmi-nyilatkozat", "adatkezelesi-tajekoztato", "adatvedelmi-tajekoztato", "adatvedelem", "adatkezeles", "nyilatkozat", "terms-and-conditions", "aszf", "privacy-policy"}
	o.re = regexp.MustCompile(`\s+`)

	return o
}

//Start starts scraping
func (c *Collector) Start() (e Entity) {
	links, all, redirectedTo, err := c.collectLinks2()
	if err != nil {
		log.Println(err)
		return
	}

	//fmt.Println(redirectedTo)

	e.MainDomain = c.Opts.Domain
	e.RedirectsTo = redirectedTo
	e.Links = all

	var data []Info

	sem := make(chan bool, 2)
	var mu sync.Mutex
	var bow = client.CreateClient2()
	for _, link := range links {
		sem <- true
		go func(link string) {
			defer func() { <-sem }()
			//fmt.Println("Processing\t", link)
			//d, bp, err := c.bpLinkWithChrome(link, c.c.ChromeTimeout)
			//d, bp, err := c.bpLinkWithHTTPClient(link)
			err = bow.Open(link)
			if err != nil {
				log.Println(err)
				return
			}

			d, err := bow.Dom().Html()
			if err != nil {
				log.Println(err)
				return
			}
			d = c.Opts.re.ReplaceAllString(d, " ")
			bp, _ := boiler.Getboiler(strings.NewReader(d))

			text, err := postprocess.Tokenize(strings.NewReader(d))
			if err != nil {
				log.Println(err)
				return
			}
			//fmt.Println("[INFO]", "Success:\t", link)
			mu.Lock()
			data = append(data, Info{Link: link, RawHTML: d, RawText: text, Boiler: bp, TimeStamp: time.Now()})
			mu.Unlock()
		}(link)
		//fmt.Println(data)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
	e.Texts = data

	/*for _, l := range e.Texts {
		fmt.Println("Iterating", l.Link)
		e.Links = append(e.Links, l.Link)
	}*/
	return
}

//collectLinks parses 1st page without creating Chrome instance,
//and collects all links which are belogs to same host, and filter out
//links to pdf and doc files.
func (c *Collector) collectLinks() (links []string, alllinks []string, redirectedTo string, err error) {
	mclient := client.CreateClient()

	um := make(map[string]bool)
	uma := make(map[string]bool)

	var (
		doc *goquery.Document
	)

	doc, redirectedTo, err = mclient.GetRedirect(c.Opts.Domain)
	if err != nil {
		return links, alllinks, "", errors.Wrap(err, "Cannot get domain link")
	}

	dURL, err := url.ParseRequestURI(redirectedTo)
	if err != nil {
		return links, alllinks, "", errors.Wrap(err, "Cannot parse domain URL "+c.Opts.Domain)
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
			//fmt.Println(lURL.String(), "\t", link)
			if c.samehost(redirectedTo, link) {
				if _, ok := uma[link]; !ok {
					uma[link] = true
					alllinks = append(alllinks, link)
				}
			}
			if c.contains(link) && c.samehost(redirectedTo, link) &&
				!strings.HasSuffix(link, ".pdf") &&
				!strings.HasSuffix(link, ".doc") {
				if _, ok := um[link]; !ok {
					um[link] = true
					links = append(links, link)
				}
			}
		}
	})

	//log.Println("links collected")
	//fmt.Printf("%#v\n", links)
	return
}

func (c *Collector) collectLinks2() (links []string, alllinks []string, redirectedTo string, err error) {
	mclient := client.CreateClient2()

	um := make(map[string]bool)
	uma := make(map[string]bool)

	err = mclient.Open(c.Opts.Domain)
	if err != nil {
		return
	}
	redirectedTo = mclient.Url().String()
	//fmt.Println(c.Opts.Domain, "redirected to ", redirectedTo)
	/*
		doc, redirectedTo, err = mclient.GetRedirect(c.Opts.Domain)
		if err != nil {
			return links, alllinks, "", errors.Wrap(err, "Cannot get domain link")
		}*/

	dURL, err := url.ParseRequestURI(redirectedTo)
	if err != nil {
		return links, alllinks, "", errors.Wrap(err, "Cannot parse domain URL "+c.Opts.Domain)
	}
	links = append(links, redirectedTo)

	//doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	mclient.Find(`a`).Each(func(i int, s *goquery.Selection) {
		if href, ok := s.Attr(`href`); ok {
			lURL, err := url.Parse(href)
			if err != nil {
				log.Println(errors.Wrap(err, "problem parsing URL"))
				return
			}
			link := dURL.ResolveReference(lURL).String()
			//fmt.Println(lURL.String(), "\t", link)
			if c.samehost(redirectedTo, link) {
				if _, ok := uma[link]; !ok {
					uma[link] = true
					alllinks = append(alllinks, link)
				}
			}
			if c.contains(link) && c.samehost(redirectedTo, link) &&
				!strings.HasSuffix(link, ".pdf") &&
				!strings.HasSuffix(link, ".doc") {
				if _, ok := um[link]; !ok {
					um[link] = true
					links = append(links, link)
				}
			}
		}
	})

	//log.Println("links collected")
	//fmt.Printf("%#v\n", links)
	return
}

/*
func (c *Collector) bpLinkWithChrome(link string, timeout int) (data string, bp string, err error) {
	data, err = browser.GetText2(link, timeout) //returns rawHTML
	if err != nil {
		return data, bp, errors.Wrap(err, "can't get text for link:"+link)
	}
	data = c.Opts.re.ReplaceAllString(data, " ")
	bp, err = boiler.Getboiler(strings.NewReader(data))
	//bp, err = boiler.Tika(strings.NewReader(data))
	if err != nil {
		log.Println(err)
	}
	return

}
*/

func (c *Collector) bpLinkWithHTTPClient(link string) (data string, bp string, err error) {
	var b []byte
	client := client.CreateClient()
	b, err = client.GetByte2(link) //returns rawHTML
	if err != nil {
		return data, bp, errors.Wrap(err, "can't get text for link:"+link)
	}

	data = string(b)
	data = c.Opts.re.ReplaceAllString(data, " ")
	bp, err = boiler.Getboiler(strings.NewReader(data))
	//bp, err = boiler.Tika(strings.NewReader(data))
	if err != nil {
		log.Println(err)
	}
	return

}

func (c *Collector) contains(link string) bool {

	for _, keyword := range c.Opts.Keywords {
		if strings.Contains(link, keyword) {
			return true
		}
	}
	return false

}

//samehost just comparing two link and figure if they are belong to same host
//in a way that http://a.mysite.com will be same host as http://b.mysite.com
func (c *Collector) samehost(dlink string, link string) bool {
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
