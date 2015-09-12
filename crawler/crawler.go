package crawler

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/deckarep/golang-set"
	"golang.org/x/net/html"
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
	return &Crawler{chLinks: mapset.NewSet(), result: mapset.NewSet(),
		depth: depth, search: search, parallel: parallel}
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

func (c *Crawler) Crawl(url *url.URL) {
	c.wq.Add(1)
	c.chLinks.Add(url.String())
	c.htmlParser(url, 0)
	c.wq.Wait()

}

func (c *Crawler) GetResult() mapset.Set {
	return c.result
}
