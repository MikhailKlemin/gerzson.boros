package database

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/MikhailKlemin/gerzson.boros/collector/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

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

	}
	//fmt.Println("Inserted post with ID:", r.InsertedID)
}
