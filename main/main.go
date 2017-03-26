package main

import (
	"flag"
	"fmt"
)

var c = flag.Bool("c", false, "create index")
var p = flag.Int("p", 1, "process number")
var m = flag.Int("m", 0, "multi-process mod")

func main() {
	flag.Parse()
	var err error
	if *c {
		rft, err = initRawFileTable()
		if err != nil {
			fmt.Println(err); return
		}
		e := createFiles(readWordWithoutStopWords, rft.num)
		if e != nil {
			fmt.Println(e); return
		}
	} else {
		var query string
		fmt.Print("Search>")
		for fmt.Scanf("%s", &query); query != "!"; fmt.Scanf("%s", &query) {
			 pairs := search(query)
			 for _, p := range pairs {
				 fmt.Println("FileNo:", p.rawFId, "Offset:", p.rawOff)
			 }
			fmt.Print("Search>")
		}
	}
}