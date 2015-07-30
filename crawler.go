package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func main() {

	var adr = flag.String("flag", "http://xmpp.org", "http address")
	flag.Parse()

	//make reguest
	resp, _ := http.Get(*adr)

	//return new html tokenizer for the given Reader
	tz := html.NewTokenizer(resp.Body)
	for {
		// scan the next token and return its type
		token := tz.Next()

		switch token {
		case html.ErrorToken:
			// End of the document
			return
		case html.StartTagToken, html.SelfClosingTagToken:
			tag, hattr := tz.TagName()
			isAnchor := string(tag)

			if isAnchor == "a" && hattr {
				_, v, _ := tz.TagAttr()
				res := string(v)
				if strings.HasPrefix(res, "http://") {
					fmt.Println(res)
				} else {
					fmt.Println(*adr + res)
				}
			}
		}
	}
	resp.Body.Close()
}
