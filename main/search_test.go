package main

import (
	"testing"
)

var testDatas = []struct{
	s string
	p pair
}{{"remember", pair{0, 125}},
	{"poor", pair{0, 194}},
	{"fantastical", pair{7, 721}},
	{"delivered", pair{13, 2412}},
	{"patience", pair{17, 95289}},
	{"murderous", pair{25, 229}},
}

func MustContain(t *testing.T, pairs []pair, p pair) {
	for _, pr := range pairs {
		if pr.rawFId == p.rawFId && pr.rawOff == p.rawOff {
			return
		}
	}
	t.Error("Must contain:", p, "get:", pairs)
}

func Test_search(t *testing.T) {
	for _, d := range testDatas {
		MustContain(t, search(d.s), d.p)
	}
}