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

func linkMaker(curr *url.URL, l string) (*url.URL, error) {

	log.Println("Original link", l)
	if !strings.HasSuffix(curr.Path, "/") {
		curr.Path += "/"
	}

	if u, err := curr.Parse(l); err == nil {
		return u, nil
	} else {
		return nil, errors.New("Bad link")
	}
}

func htmlParser(curr *url.URL, cdepth int) {

	log.Println("Start parsing", curr)
	var resp *http.Response

	//make reguest
	if xresp, err := http.Get(curr.String()); err != nil {
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
				curr.String())
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
					if newlink, err := linkMaker(curr, string(v)); err == nil {
						log.Println("Fixed link", newlink)
						result.Add(newlink)
						if cdepth < *depth && chLinks.Add(newlink.String()) {
							log.Println("crawled list", chLinks)
							htmlParser(newlink, cdepth+1)
						} else {
							log.Println("Already crawled", newlink)
						}
					} else {
						log.Println(err)
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
	log.Println("get properties from command line", *adr, *depth)
	u, err := url.Parse(*adr)
	if err != nil {
		log.Println("Bad link", err)
	} else {
		chLinks.Add(u.String())
		htmlParser(u, 0)

		log.Println("print links")
		for val := range result.Iter() {
			fmt.Println(val)
		}
	}

}
