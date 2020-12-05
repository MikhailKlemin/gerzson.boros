package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/MikhailKlemin/gerzson.boros/collector"
)

//var domains = []string{"femina.hu", "totalcar.hu", "velvet.hu", "telekom.hu", "rtl.hu", "emag.hu", "portfolio.hu", "eropolis.hu", "ripost.hu", "argep.hu", "t-online.hu", "prohardver.hu", "napi.hu", "nosalty.hu", "bme.hu", "sorozatjunkie.hu", "mestermc.hu", "love.hu", "keptelenseg.hu", "e-kreta.hu", "oktatas.hu", "blogstar.hu", "csubakka.hu", "mozanaplo.hu", "hwsw.hu", "liked.hu", "hupont.hu", "jysk.hu", "aczelauto.hu", "aczelestarsa.hu", "aczelpetra.hu", "ad.hu", "ad-media.hu", "ad6kap6.hu", "adab.hu"}

//Config is configuration
type Config struct {
	OutDir     string `json:"output_directory"`
	DomainPath string `json:"path_to_domain_file"`
}

func main() {

	var c Config
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal("Cannot read config file " + err.Error())
	}
	err = json.Unmarshal(b, &c)
	if err != nil {
		log.Fatal("Cannot read json file " + err.Error())

	}
	start(c.OutDir, c.DomainPath)

}

func start(outDir, domainPath string) {

	domains := loaddomains(domainPath)
	domains = domains[:10]
	t := time.Now()
	tt := time.Now()
	sem := make(chan bool, 40)

	fmt.Println("Total domains to scrap:\t", len(domains))

	for i, dl := range domains {
		if i%100 == 0 && i != 0 {
			fmt.Printf("Processing: %d time per batch: %s \n",
				i, time.Since(t))
			t = time.Now()
		}
		sem <- true
		go func(dl string) {
			defer func() { <-sem }()
			d := collector.NewCollector("http://" + dl)
			data := d.Start()
			sample, _ := json.MarshalIndent(data, "", "    ")
			if err := ioutil.WriteFile(filepath.Join(outDir, dl+".json"), sample, 0600); err != nil {
				log.Println(err)
			}
		}(dl)

	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	fmt.Printf("Took %s to do\n", time.Since(tt))
}

//loaddomains loadsdomain from the text file
func loaddomains(path string) (dls []string) {

	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	raw := strings.Split(string(b), "\n")

	for _, rd := range raw {
		rd = strings.TrimSpace(rd)
		_, err := url.ParseRequestURI("http://" + rd)
		if err != nil {
			log.Println("Bad URI\t", rd)
		} else {
			dls = append(dls, rd)
		}
	}
	return
}
