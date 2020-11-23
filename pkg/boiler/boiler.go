package boiler

import (
	"io"

	"github.com/jlubawy/go-boilerpipe"
	"github.com/pkg/errors"
)

/*
//I do not use Tika extra dependance for not much ofimporvement.

var client = tika.NewClient(nil, "http://localhost:9998/")

//Tika connects to Tika srv and do the thing
func Tika(b io.Reader) (string, error) {
	body, err := client.Parse(context.Background(), b)
	if err != nil {
		return "", errors.Wrap(err, "can't send req to Tika")
	}
	ps, _ := client.Parsers(context.Background())
	fmt.Printf("%#v\n", ps)
	return body, nil
}
*/

//Getboiler parse Go boilepipe
func Getboiler(r io.Reader) (string, error) {
	doc, err := boilerpipe.ParseDocument(r)
	if err != nil {
		//log.Fatal(err)
		return "", errors.Wrap(err, "can't parse doc with goboiler")
	}
	boilerpipe.ArticlePipeline.Process(doc)
	//return doc.Text(true, true), nil
	return doc.Content(), nil

}
