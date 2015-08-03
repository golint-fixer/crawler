package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

func linkMaker(root, curr *string, l string) string {
	if strings.HasPrefix(l, "http://") ||
		(strings.HasPrefix(l, "https://")) {
		return l
	}
	if strings.HasPrefix(l, "/") {
		return (*root + l)
	} else {
		return (*curr + "/" + l)
	}

}

func htmlParser(root, curr *string) []string {

	result := make([]string, 50)
	//make reguest
	resp, _ := http.Get(*curr)

	//return new html tokenizer for the given Reader
	tz := html.NewTokenizer(resp.Body)
	for {
		// scan the next token and return its type
		token := tz.Next()

		switch token {
		case html.ErrorToken:
			// End of the document
			return result
		case html.StartTagToken, html.SelfClosingTagToken:
			tag, hattr := tz.TagName()
			isAnchor := string(tag)

			if isAnchor == "a" && hattr {
				_, v, _ := tz.TagAttr()
				result = append(result,
					linkMaker(root, curr, string(v)))
			}
		}
	}
	resp.Body.Close()

	return result
}

func main() {

	var adr = flag.String("flag", "http://xmpp.org", "http address")
	flag.Parse()

	larr := htmlParser(adr, adr)
	for _, val := range larr {
		fmt.Println(val)
	}
}
