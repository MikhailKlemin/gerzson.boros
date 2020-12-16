package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/MikhailKlemin/gerzson.boros/collector"
	"github.com/MikhailKlemin/gerzson.boros/collector/config"
	"github.com/MikhailKlemin/gerzson.boros/collector/database"
)

//var domains = []string{"femina.hu", "totalcar.hu", "velvet.hu", "telekom.hu", "rtl.hu", "emag.hu", "portfolio.hu", "eropolis.hu", "ripost.hu", "argep.hu", "t-online.hu", "prohardver.hu", "napi.hu", "nosalty.hu", "bme.hu", "sorozatjunkie.hu", "mestermc.hu", "love.hu", "keptelenseg.hu", "e-kreta.hu", "oktatas.hu", "blogstar.hu", "csubakka.hu", "mozanaplo.hu", "hwsw.hu", "liked.hu", "hupont.hu", "jysk.hu", "aczelauto.hu", "aczelestarsa.hu", "aczelpetra.hu", "ad.hu", "ad-media.hu", "ad6kap6.hu", "adab.hu"}
//var db *database.Datastore

func main() {
	//log := logrus.New()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	c := config.LoadGeneralConfig()
	//database.Print(c)
	//os.Exit(1)

	//ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	//defer cancel()
	//db = database.NewDatastore(ctx, c)

	/*dbs, err := db.Client.ListDatabaseNames(ctx, bson.M{"name": primitive.Regex{Pattern: ".*"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.Join(dbs, "\t"))
	*/

	//defer db.Client.Disconnect(ctx)

	//start(c.OutDir, c.DomainPath)
	start(c)

	//db := database.NewDatastore(c, log)
	//db.Session.Connect()
}

//func start(outDir, domainPath string) {
func start(conf config.GeneralConfig) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db := database.NewDatastore(ctx, conf)
	defer db.Client.Disconnect(ctx)

	domains := loaddomains(conf.DomainPath)
	//domains = []string{"index.hu"}
	//domains = domains[:10]

	/*	db, err := leveldb.OpenFile(conf.LevelDBPath, nil)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
	*/
	t := time.Now()

	tt := time.Now()

	//sem := make(chan bool, conf.Concurrency)

	fmt.Println("Total domains to scrap:\t", len(domains))

	tasks := make(chan string, len(domains))
	results := make(chan collector.Entity, len(domains))
	col := collector.NewCollector(conf)

	worker := func(tasks <-chan string, results chan<- collector.Entity) {
		for dlink := range tasks {
			//do shit
			e, err := col.Start("http://" + dlink)
			if err != nil {
				//log.Println(err)
			}
			if e.MainDomain == "" {
				e.MainDomain = "http://" + dlink
			}
			results <- e
			//fmt.Println(dlink)
		}
	}

	for w := 0; w <= conf.Concurrency; w++ {
		go worker(tasks, results)
	}

	for _, d := range domains {
		tasks <- d

	}

	close(tasks)

	counter := 0
	for a := 0; a < len(domains); a++ {
		fmt.Printf("%d                          \r", a)
		if a%1000 == 0 && a != 0 {
			fmt.Printf("Processing: %d time per batch: %s, empty links per batch: %d \n",
				a, time.Since(t),
				counter)
			t = time.Now()
			counter = 0
		}
		entity := <-results
		if len(entity.Links) == 0 {
			counter++
		}
		//b, _ := json.Marshal(entity)
		//err := db.Put([]byte(entity.MainDomain), b, nil)
		db.Insert(entity)
		/*
			if err != nil {
				log.Fatal(err)
			}
		*/

	}

	/*
		for i, dl := range domains {
			if i%100 == 0 && i != 0 {
				fmt.Printf("Processing: %d time per batch: %s \n",
					i, time.Since(t))
				t = time.Now()
			}
			sem <- true
			//fmt.Println("[Started:]\t", dl, "\tLeft:", len(domains)-i)

			go func(dl string) {
				defer func() { <-sem }()
				d := collector.NewCollector("http://"+dl+"/", conf)
				data := d.Start()
				if data.MainDomain == "" {
					log.Println("Empty domain")
					data.MainDomain = dl
				}
				//			fmt.Println(len(data.Texts))
				db.Insert(data)

				//	sample, _ := json.MarshalIndent(data, "", "    ")
				//	if err := ioutil.WriteFile(filepath.Join(conf.OutDir, dl+".json"), sample, 0600); err != nil {
				//		log.Println(err)
				//	}


			}(dl)

		}

		for i := 0; i < cap(sem); i++ {
			sem <- true
		}
	*/
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
		_, err := url.ParseRequestURI("https://" + rd)
		if err != nil {
			log.Println("Bad URI\t", rd)
			continue
		}
		dls = append(dls, rd)

	}
	return
}
