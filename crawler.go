package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/deckarep/golang-set"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"strings"
)

//Links that checked or started to check
var chLinks mapset.Set
var depth *int
var result mapset.Set

func linkMaker(curr string, l string) (string, error) {

	log.Println("Original link", l)
	if !strings.HasSuffix(curr, "/") {
		curr = curr + "/"
	}

	if v, verr := url.Parse(curr); verr == nil {

		if u, err := v.Parse(l); err == nil {
			return u.String(), nil
		}
	}
	return "", errors.New("no url")

}

func htmlParser(curr string, cdepth int) {

	log.Println("Start parsing", curr)
	var resp *http.Response

	//make reguest
	if xresp, err := http.Get(curr); err != nil {
		log.Println("Bad url", curr, err)
		return
	} else {
		resp = xresp
	}
	defer resp.Body.Close()

	//return new html tokenizer for the given Reader
	tz := html.NewTokenizer(resp.Body)
	for {
		// scan the next token and return its type
		token := tz.Next()
		switch token {

		// End of the document or error
		case html.ErrorToken:
			log.Println("End of document or error on the page ",
				curr)
			return
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
				if string(k) == "href" {
					newlink, err := linkMaker(curr, string(v))
					if err == nil {
						log.Println("Fixed link ", newlink)
						result.Add(newlink)
					} else {
						log.Println("Bad link")
					}
					if cdepth < *depth && chLinks.Add(newlink) {
						htmlParser(newlink, cdepth+1)
					}
				}
			}
		}
	}

}

func main() {
	chLinks = mapset.NewSet()
	result = mapset.NewSet()

	// Init values for the standart logger
	// Lshortfile - file name and file number
	log.SetFlags(log.Lshortfile)
	var adr = flag.String("url", "http://xmpp.org", "http address")
	depth = flag.Int("depth", 5, "depth of searching")
	flag.Parse()
	log.Println("get properties from command line")
	log.Println(*depth)

	chLinks.Add(*adr)
	htmlParser(*adr, 0)

	log.Println("print links")
	for val := range result.Iter() {
		fmt.Println(val)
	}

}
