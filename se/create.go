package main

import (
	"io"
	"github.com/inszva/miniSE/DictionaryGenerator"
	"os"
	"strconv"
	"bufio"
	"strings"
	"path/filepath"
)

type rawFileTable struct {
	idx int
	num int
	f *os.File
	reader *bufio.Reader
	off uint32
}

var rft *rawFileTable

func initRawFileTable() (*rawFileTable, error) {
	rft := &rawFileTable{idx:0, num:0}
	err := filepath.Walk("./raw", func(path string, f os.FileInfo, err error) error {
		if f == nil { return err }
		if f.IsDir() { return nil }
		splited := strings.Split(f.Name(), ".")
		if len(splited) < 2 { return nil }
		if splited[1] == "txt" { rft.num++ }
		return nil
	})
	if err != nil {
		return nil, err
	}
	return rft, nil
}

func readWord(rawFId int) (word string, off uint32, err error) {
	if rft == nil {
		panic("rawFileTable is nil")
	}

	defer func() {
		if err != nil {
			err = io.EOF
			rft.f = nil
		}
	} ()

	if rawFId == rft.idx {
		if rft.f == nil {
			rft.f, err = os.OpenFile("./raw/" + strconv.Itoa(rawFId) + ".txt", os.O_RDONLY, 0666)
			if err != nil { return }
			rft.reader = bufio.NewReader(rft.f)
			rft.off = 0
		}
	} else {
		rft.idx = rawFId
		rft.f, err = os.OpenFile("./raw/" + strconv.Itoa(rawFId) + ".txt", os.O_RDONLY, 0666)
		if err != nil { return }
		rft.reader = bufio.NewReader(rft.f)
		rft.off = 0
	}

	wordBuf := []byte{}
	var e error
	var b byte
	for b, e = rft.reader.ReadByte(); e == nil; b, e = rft.reader.ReadByte() {
		rft.off++
		if b >= 'A' && b <= 'Z' {
			b = b - 'A' + 'a'
		}
		if b >= 'a' && b <= 'z' {
			wordBuf = append(wordBuf, b)
		} else {
			if len(wordBuf) > 0 { break }
			continue
		}
	}
	if len(wordBuf) > 0 {
		word = string(wordBuf)
		off = rft.off - uint32(len(wordBuf))
		return
	}
	return "", 0, io.EOF
}

func readWordWithoutStopWords(rawFId int) (word string, off uint32, err error) {
	for word, off, err = readWord(rawFId); err == nil; word, off, err = readWord(rawFId) {
		word = DictionaryGenerator.Parse(word)
		if word != "" {
			return
		}
	}
	return
}