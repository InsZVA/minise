// Whoever can not prevent me using GOTO!!!
package main

import (
	"regexp"
	"net/http"
	"io/ioutil"
	"sync/atomic"
	"strconv"
	"sync"
	"os"
	"bufio"
	"log"
)

var fullHtmlURLs = make(chan string, 20)
// Atomic integer to identify the file
var atomicFileNo int32

// The producer: produce the url to download
func getFullHtmlURL() (err error) {
	var regexpCHT, regexpSonnets, regexpPoetry  *regexp.Regexp
	regexpCHT, err = regexp.Compile(`<br><a href="(.*)/index.html">`)
	if err != nil { return }
	regexpSonnets, err = regexp.Compile(`<DT><A HREF="(.*)">`)
	if err != nil { return }
	regexpPoetry, err = regexp.Compile(`<em><a href="(.*)">`)
	if err != nil { return }

	var resp *http.Response
	var dataStr []byte
	// Get Comedy, History, Tragedy full html url
	resp, err = http.Get("http://shakespeare.mit.edu/")
	if err != nil { return }
	defer resp.Body.Close()
	dataStr, err = ioutil.ReadAll(resp.Body)
	if err != nil { return }
	found := regexpCHT.FindAllStringSubmatch(string(dataStr), -1)
	for i := 0; i < len(found); i++ {
		fullHtmlURLs <- `http://shakespeare.mit.edu/` + found[i][1] + `/full.html`
	}
	// Get Poetry(except The Sonnets) full html url
	found = regexpPoetry.FindAllStringSubmatch(string(dataStr), -1)
	for i := 0; i < len(found); i++ {
		fullHtmlURLs <- `http://shakespeare.mit.edu/` + found[i][1]
	}

	// Get The Sonnets full html url
	var resp2 *http.Response
	resp2, err = http.Get("http://shakespeare.mit.edu/Poetry/sonnets.html")
	if err != nil { return }
	defer resp2.Body.Close()
	dataStr, err = ioutil.ReadAll(resp2.Body)
	if err != nil { return }
	found = regexpSonnets.FindAllStringSubmatch(string(dataStr), -1)
	for i := 0; i < len(found); i++ {
		fullHtmlURLs <- `http://shakespeare.mit.edu/Poetry/` + found[i][1]
	}
	fullHtmlURLs <- "23333333"
	return
}

func getURLContentText(url string) (err error) {
	var resp *http.Response
	var dataRaw []byte
	var dataStr string
	var regexpTable, regexpText *regexp.Regexp
reGet:
	resp, err = http.Get(url)
	if err != nil { return }
	dataRaw, err = ioutil.ReadAll(resp.Body)
	if int64(len(dataRaw)) < resp.ContentLength {
		// Reconnect to pull full response
		goto reGet
	}
	dataStr = string(dataRaw)
	if err != nil { return }

	// Remove the Navigation in the HTML
	regexpTable, err = regexp.Compile(`<table width="100%" bgcolor="#CCF6F6">[\s\S]*?</table>`)
	if err != nil { return }
	dataStr = regexpTable.ReplaceAllString(dataStr, "")
	id := atomic.LoadInt32(&atomicFileNo)

	// Peek the Text
	regexpText, err = regexp.Compile(`>\s*(.*?)\s*<`)
	if err != nil { return }
	text := regexpText.FindAllStringSubmatch(dataStr, -1)

	// Write to file
	var f *os.File
	for !atomic.CompareAndSwapInt32(&atomicFileNo, id, id+1) {
		id = atomic.LoadInt32(&atomicFileNo)
	}
	f, err = os.OpenFile("./out/"+ strconv.Itoa(int(id)) + ".txt", os.O_CREATE | os.O_RDWR, 0666)
	if err != nil { return }
	bufWriter := bufio.NewWriter(f)
	for i := 0; i < len(text); i++ {
		bufWriter.WriteString(text[i][1] + "\n")
	}
	bufWriter.Flush()
	f.Close()
	return
}

func main() {
	go func () {
		getFullHtmlURL()
	} ()
	var wp sync.WaitGroup
URL:
	for {
		select {
		case url := <- fullHtmlURLs:
			if url != "23333333" {
				wp.Add(1)
				go func () {
					log.Println("Handle:" + url)
					retry := 0
					for err := getURLContentText(url); err != nil && retry < 3; retry++ {
						log.Println("Retry", retry, err.Error())
						retry++
						err = getURLContentText(url)
					}
					wp.Done()
				} ()
			} else { break URL }
		}
	}
	wp.Wait()
}