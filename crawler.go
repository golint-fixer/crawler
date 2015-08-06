package main

import (
	"flag"
	"fmt"
	"github.com/deckarep/golang-set"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strings"
)

//Links that checked or started to check
var chLinks mapset.Set
var depth *int
var result mapset.Set

func linkMaker(root, curr string, l string) string {

	log.Println("Original link", l)

	if strings.HasPrefix(l, "http://") ||
		(strings.HasPrefix(l, "https://")) {
		log.Println("Fixed link", l)
		return l
	}
	if strings.HasPrefix(l, "/") {
		log.Println("Fixed link", root+l)
		return (root + l)
	} else {
		log.Println("Fixed link", curr+"/"+l)
		return (curr + "/" + l)
	}

}

func htmlParser(root, curr string, cdepth int) {

	log.Println("Start parsing", curr)

	//make reguest
	resp, _ := http.Get(curr)
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
					newlink := linkMaker(root, curr, string(v))
					result.Add(newlink)
					if cdepth < *depth && chLinks.Add(newlink) {
						if strings.HasPrefix(newlink, root) {
							htmlParser(root, newlink, cdepth+1)
						} else {
							htmlParser(curr, newlink, cdepth+1)
						}
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
	htmlParser(*adr, *adr, 0)

	log.Println("print links")
	for val := range result.Iter() {
		fmt.Println(val)
	}

}
