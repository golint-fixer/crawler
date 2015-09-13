package crawler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/deckarep/golang-set"
)

type testpairlink struct {
	url, u, res string
}

var tests = []testpairlink{
	{"http://xmpp.org/value", "value2", "http://xmpp.org/value/value2"},
	{"http://xmpp.org/value", "/value2", "http://xmpp.org/value2"},
	{"http://xmpp.org/value", "#value2", "http://xmpp.org/value/"},
	{"http://xmpp.org/value/", "#value2", "http://xmpp.org/value/"},
}

var page = `<HTML>
<HEAD>
<TITLE>Your Title Here</TITLE>
</HEAD>
<HR>
<a href="%%">error in link</a>
<a href="#anchor">anchor in link</a>
<a>tag a without href</a>
<a class="accessibility" rel=next media="not print" href="#content">Skip to content</a>
<a href="http://somegreatsite.com">Link Name</a>
is a link to another nifty site
<H1>This is a Header</H1>
<H2>This is a Medium Header</H2>
Send me mail at <a href="mailto:support@yourcompany.com">
support@yourcompany.com</a>.
<P> This is a new paragraph!
<P> <B>This is a new paragraph!</B>
<a href="http://xmpp.org/about-xmpp/">About</a>
<a href="http://xmpp.org/about-xmpp/faq/">FAQ</a>
<BR> <B><I>This is a new sentence without a paragraph break, in bold italics.</I></B>
<a href="http://wordpress.org/" rel="generator">WordPress</a>
<HR>
</BODY>
</HTML>`

var resultlinks = []string{"http://somegreatsite.com", "mailto:support@yourcompany.com",
	"http://xmpp.org/about-xmpp/", "http://xmpp.org/about-xmpp/faq/",
	"http://wordpress.org/"}

func TestLinkMaker(t *testing.T) {
	for _, pair := range tests {
		ur, _ := url.Parse(pair.url)
		lmres, _ := linkMaker(ur, pair.u)
		res, _ := url.Parse(pair.res)
		if *lmres != *res {
			t.Error("For", pair.url, pair.u, "expected", pair.res,
				"got", lmres.String())
		}
	}
}

func TestLinkMakerError(t *testing.T) {
	clink := "http://xmpp.org"
	nlink := "%%"
	ur, _ := url.Parse(clink)
	lmres, e := linkMaker(ur, nlink)
	if e == nil {
		t.Error("For", clink, nlink, "expected Error message, got", lmres)
	}
}

func TestHtmlParser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, page)
	}))
	defer ts.Close()

	//Init values
	u, err := url.Parse(ts.URL)
	if err != nil {
		log.Println(err)
	} else {

		var depth = 0
		var search = true
		var parallel = false

		c := NewCrawler(depth, search, parallel)
		c.Crawl(u)
		answ := mapset.NewSet()
		for _, v := range resultlinks {
			answ.Add(v)
		}
		if !answ.Equal(c.GetResult()) {
			t.Error("Expected links and got are not the same")
		}
	}
}
