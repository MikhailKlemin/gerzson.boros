package postprocess

import (
	"io"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/net/html"
)

//Tokenize tokenize HTLML
//This function is used to tokenize HTML to
//represent it clean TEXT
func Tokenize(r io.Reader) (string, error) {
	textTags := []string{
		"a",
		"p", "span", "em", "string", "blockquote", "q", "cite",
		"h1", "h2", "h3", "h4", "h5", "h6", "pre", "ul", "li", "ol",
		"mark", "ins", "del", "small", "i", "b",
	}

	tag := ""
	enter := false
	var text []string
	tokenizer := html.NewTokenizer(r)
	for {
		tt := tokenizer.Next()
		token := tokenizer.Token()

		err := tokenizer.Err()
		if err == io.EOF {
			break
		}

		switch tt {
		case html.ErrorToken:
			//log.Fatal(err)
			return "", errors.Wrap(err, "can't parse token")
		case html.StartTagToken, html.SelfClosingTagToken:
			enter = false

			tag = token.Data
			for _, ttt := range textTags {
				if tag == ttt {
					enter = true
					break
				}
			}
		case html.TextToken:
			if enter {
				data := strings.TrimSpace(token.Data)

				if len(data) > 0 {
					//fmt.Println(data)
					text = append(text, data)
				}
			}
		}
	}

	return strings.Join(text, " "), nil
}
