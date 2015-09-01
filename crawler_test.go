package main

import (
	"net/url"
	"testing"
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

func TestLinkMaker(t *testing.T) {
	for _, pair := range tests {
		ur, _ := url.Parse(pair.url)
		lmres, _ := linkMaker(ur, pair.u)
		res, _ := url.Parse(pair.res)
		if *lmres != *res {
			t.Error("For", pair.url, pair.u, "expected", pair.res, "got", lmres.String())
		}
	}
}
