package main

import (
	"testing"
	"fmt"
)

func TestInvertFile_addPair(t *testing.T) {
	var ifile invertFile
	var e error
	var mf *mmFile
	mf, e = newMMFile("./test/test.inv", 1024)
	if e != nil { return }
	ifile.ReadClosableMMFile = newReadClosableMMFile(mf)
	if e != nil {
		t.Error(e)
	}
	for i := 0; i < 2048; i++ {
		fmt.Println(ifile.addPair([]pair{{rawFId: 97, rawOff: 98}}))
	}
	ifile.Close()
}
