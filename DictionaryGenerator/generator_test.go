package DictionaryGenerator

import (
	"testing"
	"os"
	"bufio"
)

var testData = []struct {
		raw string
		expected string
}{
	{"a", ""},
	{"clicked", "click"},
	{"mapping", "map"},
	{"readable", "readabl"},
}

func TestParse(t *testing.T) {
	// Test stop words
	f, err := os.OpenFile("./stopwords.txt", os.O_RDONLY, 0666)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	for s, err := reader.ReadString('\n'); err == nil; s, err = reader.ReadString('\n') {
		if real := Parse(s); real != "" {
			t.Error("Parse:", s, "Expected:", "", "Real:", real)
		}
	}

	// Test stem
	for _, td := range testData {
		if real := Parse(td.raw); real != td.expected {
			t.Error("Parse:", td.raw, "Expected:", "", "Real:", real)
		}
	}
}