package main

import (
	"flag"
	"fmt"
	"github.com/deckarep/golang-set"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

//Links that checked or started to check
var chLinks mapset.Set
var depth *int
var search *bool
var parallel *bool
var result mapset.Set
var wq sync.WaitGroup

func linkMaker(curr *url.URL, l string) (*url.URL, error) {

	log.Println("Original link", l)
	if !strings.HasSuffix(curr.Path, "/") {
		curr.Path += "/"
	}

	if u, err := curr.Parse(l); err == nil {
		u.Fragment = ""
		return u, nil
	} else {
		return nil, err
	}
}

func htmlParser(curr *url.URL, cdepth int) {
	log.Println("Start parsing", curr)
	var resp *http.Response
	defer wq.Add(-1)

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
				if string(k) == "href" && !strings.HasPrefix(string(v), "#") {
					if newlink, err := linkMaker(curr, string(v)); err == nil {
						log.Println("Fixed link", newlink)
						if curr.Host != newlink.Host && !*search {
							break
						}
						result.Add(newlink.String())
						if cdepth < *depth && chLinks.Add(newlink.String()) {
							wq.Add(1)
							if *parallel {
								go htmlParser(newlink, cdepth+1)
							} else {
								htmlParser(newlink, cdepth+1)
							}
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
	search = flag.Bool("search", false, "search in all hostname")
	parallel = flag.Bool("parallel", false, "perfom in parallel")
	flag.Parse()
	log.Println("get properties from command line",
		*adr, *depth, *search, *parallel)
	u, err := url.Parse(*adr)
	if err != nil {
		log.Println("Bad link", err)
	} else {
		chLinks.Add(u.String())
		wq.Add(1)
		htmlParser(u, 0)
		wq.Wait()
		log.Println("print links")
		for val := range result.Iter() {
			fmt.Println(val)
		}
	}
}
