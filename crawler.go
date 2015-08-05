package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
)

func linkMaker(root, curr *string, l string) string {

	log.Println("original link", l)

	if strings.HasPrefix(l, "http://") ||
		(strings.HasPrefix(l, "https://")) {
		log.Println("fixed link", l)
		return l
	}
	if strings.HasPrefix(l, "/") {
		log.Println("fixed link", *root+l)
		return (*root + l)
	} else {
		log.Println("fixed link", *curr+"/"+l)
		return (*curr + "/" + l)
	}

}

func htmlParser(root, curr *string) []string {

	result := make([]string, 0, 50)

	//make reguest
	resp, _ := http.Get(*curr)

	//return new html tokenizer for the given Reader
	tz := html.NewTokenizer(resp.Body)
	for {
		// scan the next token and return its type
		token := tz.Next()
		switch token {

		// End of the document or error
		case html.ErrorToken:
			resp.Body.Close()
			log.Println("End of document or error on the page ",
				*curr)
			return result
		case html.StartTagToken, html.SelfClosingTagToken:
			tag, _ := tz.TagName()
			isAnchor := string(tag)

			if isAnchor == "a" {

				//k-key, v-value,
				//mattr-there are more attributes
				k, v, mattr := tz.TagAttr()
				for mattr && string(k) != "href" {
					k, v, mattr = tz.TagAttr()
				}
				result = append(result,
					linkMaker(root, curr, string(v)))
			}
		}
	}

}

func main() {

	// Init values for the standart logger
	// Lshortfile - file name and file number
	log.SetFlags(log.Lshortfile)
	var adr = flag.String("url", "http://xmpp.org", "http address")
	flag.Parse()
	log.Println("get properties from command line")

	larr := htmlParser(adr, adr)

	log.Println("print links")
	for _, val := range larr {
		fmt.Println(val)
	}
}
