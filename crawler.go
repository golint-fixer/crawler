package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"

	"github.com/ksheremet/crawler/crawler"
)

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
		c := crawler.NewCrawler(*depth, *search, *parallel)
		c.Crawl(u)
		log.Println("print links")

		for val := range c.GetResult().Iter() {
			fmt.Println(val)
		}
	}
}
