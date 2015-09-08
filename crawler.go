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

type Crawler struct {
	//Links that checked or started to check
	chLinks  mapset.Set
	result   mapset.Set
	depth    int
	search   bool
	parallel bool
	wq       sync.WaitGroup
}

func NewCrawler(depth int, search, parallel bool) *Crawler {
	return &Crawler{mapset.NewSet(), mapset.NewSet(), depth, search,
		parallel, *new(sync.WaitGroup)}
}

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

func (c *Crawler) htmlParser(curr *url.URL, cdepth int) {
	log.Println("Start parsing", curr)
	var resp *http.Response
	defer c.wq.Add(-1)

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
						if curr.Host != newlink.Host && !c.search {
							break
						}
						c.result.Add(newlink.String())
						if cdepth < c.depth && c.chLinks.Add(newlink.String()) {
							c.wq.Add(1)
							if c.parallel {
								go c.htmlParser(newlink, cdepth+1)
							} else {
								c.htmlParser(newlink, cdepth+1)
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
	// Init values for the standart logger
	// Lshortfile - file name and file number
	log.SetFlags(log.Lshortfile)
	var adr = flag.String("url", "http://xmpp.org", "http address")
	var depth = flag.Int("depth", 5, "depth of searching")
	var search = flag.Bool("search", false, "search in all hostname")
	var parallel = flag.Bool("parallel", false, "perfom in parallel")
	flag.Parse()
	log.Println("get properties from command line",
		*adr, *depth, *search, *parallel)
	u, err := url.Parse(*adr)
	if err != nil {
		log.Println("Bad link", err)
	} else {
		//Create struct Crawler
		c := NewCrawler(*depth, *search, *parallel)
		c.wq.Add(1)
		c.chLinks.Add(u.String())
		c.htmlParser(u, 0)
		c.wq.Wait()
		log.Println("print links")
		for val := range c.result.Iter() {
			fmt.Println(val)
		}
	}
}
