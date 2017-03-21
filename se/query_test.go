package main

import (
	"testing"
	"fmt"
	"github.com/inszva/miniSE/DictionaryGenerator"
)

func Test_search(t *testing.T) {
	var hCode [20]byte
	word := DictionaryGenerator.Parse("bitterness")
	fmt.Println(word)
	hashCode(word, &hCode)
	for i := 0; i < 20; i++ {
		fmt.Printf("%x ", hCode[i])
	}
	fmt.Println(search("bitterness"))
}
