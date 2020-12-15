package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/MikhailKlemin/gerzson.boros/collector"
	"github.com/MikhailKlemin/gerzson.boros/collector/config"
	"github.com/syndtr/goleveldb/leveldb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

//Print is
func Print(conf config.GeneralConfig) {
	db, err := leveldb.OpenFile(conf.LevelDBPath, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	count := 0
	tcount := 0
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		tcount++
		//key:=iter.Key()
		val := iter.Value()
		var e collector.Entity

		err := json.Unmarshal(val, &e)
		if err != nil {
			log.Fatal(err)
		}

		if len(e.Texts) == 0 && e.MainDomain != "" {
			count++
		}

	}
	fmt.Printf("Empty %d from %d\n", count, tcount)

}

//Datastore stores DB
type Datastore struct {
	//db     *mongo.Database
	Client *mongo.Client
	once   sync.Once
	conf   config.GeneralConfig
	ctx    context.Context
}

//NewDatastore is
func NewDatastore(ctx context.Context, conf config.GeneralConfig) *Datastore {
	var d Datastore
	d.conf = conf
	d.ctx = ctx
	d.once.Do(func() {
		//fmt.Println("Hello Once")
		var err error
		opts := options.Client().ApplyURI(conf.DatabaseHost)
		d.Client, err = mongo.NewClient(opts)
		if err != nil {
			log.Fatal(err)
		}
		err = d.Client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}

		err = d.Client.Ping(context.TODO(), nil)
		if err != nil {
			log.Fatal(err)
		}

		//d.Client.Database()
	})

	fmt.Println("Hello")

	return &d

}

//Insert inserts
func (db *Datastore) Insert(v interface{}) {
	//fmt.Println("Database: ", db.conf.DatabaseName)
	//fmt.Println("Collection: ", db.conf.DatabaseCollection)

	collection := db.Client.Database(db.conf.DatabaseName).Collection(db.conf.DatabaseCollection)

	_, err := collection.InsertOne(context.TODO(), v)
	if err != nil {
		var merr mongo.WriteException
		var ok bool

		merr, ok = err.(mongo.WriteException)
		if ok {
			errCode := merr.WriteErrors[0].Code
			if errCode == 11000 {
				return
			}
		}
		log.Println(err.Error())
		return
	}

	//fmt.Println()

	//fmt.Println("Inserted post with ID:", r.InsertedID)
}
