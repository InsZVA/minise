package main

import (
	"github.com/xiang90/splaytree"
	"unsafe"
	"strconv"
	"reflect"
	"github.com/inszva/miniSE/DictionaryGenerator"
)

type cacheNode struct {
	hCode [20]byte
	pairs []pair
}

var cache splaytree.Tree

func hashLess(h1 [20]byte, h2 [20]byte) bool {
	for t := 0; t < 5; t++ {
		if *(*uint32)(unsafe.Pointer(&h1[4*t])) >
			*(*uint32)(unsafe.Pointer(&h2[4*t])) {
			return false
		} else if *(*uint32)(unsafe.Pointer(&h1[4*t])) <
			*(*uint32)(unsafe.Pointer(&h2[4*t])) {
			return true } else { continue }
	}
	return false
}

func (cn cacheNode) Less(cn2 splaytree.Item) bool {
	for t := 0; t < 5; t++ {
		cn2cn := cn2.(cacheNode)
		if *(*uint32)(unsafe.Pointer(&cn.hCode[4*t])) >
			*(*uint32)(unsafe.Pointer(&(cn2cn.hCode[4*t]))) {
			return false
		} else if *(*uint32)(unsafe.Pointer(&cn.hCode[4*t])) <
			*(*uint32)(unsafe.Pointer(&(cn2cn.hCode[4*t]))) {
			return true } else { continue }
	}
	return false
}

func init() {
	cache = splaytree.NewSplayTree()
}

func startupIndex(h byte) error {
	if index[h].ReadClosableMMFile == nil {
		path := "./index/" + strconv.Itoa(int(h)) + ".inv"
		mf, err := openMMFile(path)
		if err != nil { return err }
		index[h].ReadClosableMMFile = newReadClosableMMFile(mf)
	}
	if invert[h].ReadClosableMMFile == nil {
		path := "./invert/" + strconv.Itoa(int(h)) + ".inv"
		mf, err := openMMFile(path)
		if err != nil { return err }
		invert[h].ReadClosableMMFile = newReadClosableMMFile(mf)
	}
	return nil
}

func readPairsFromIndexOff(idx byte, i int) []pair {
	invOff := *(*uint32)(unsafe.Pointer(&index[idx].m[i+20]))
	var length uint64
	b := reflect.SliceHeader{}
	b.Cap = 8
	b.Len = 8
	b.Data = uintptr(unsafe.Pointer(&length))
	n, e := invert[idx].ReadAt(*(*[]byte)(unsafe.Pointer(&b)), int64(invOff))
	if n != 8 || e != nil {
		return nil
	}
	ret := make([]pair, length)
	b.Cap = 8 * int(length)
	b.Len = 8 * int(length)
	b.Data = uintptr(unsafe.Pointer(&ret[0].rawFId))
	n, e = invert[idx].ReadAt(*(*[]byte)(unsafe.Pointer(&b)), int64(invOff) + 8)
	if n != b.Len || e != nil {
		return nil
	}
	return ret
}

func searchFromFile(hCode [20]byte) []pair {
	if startupIndex(hCode[0]) != nil { return nil}
	n := len(index[hCode[0]].m)
	if mod := n % 24; mod != 0 {
		n -= mod
	}
	l, r := 0, n
	for l != r {
		m := (l + r) >> 1
		m = m / 24 * 24
		hCodeFB := make([]byte, 20)
		copy(hCodeFB, index[hCode[0]].m[m:])
		hCodeF := *(*[20]byte)(unsafe.Pointer(&hCodeFB[0]))
		if hashLess(hCode, hCodeF) {
			l = m + 24; continue
		} else
		if hCodeF == hCode {
			return readPairsFromIndexOff(hCode[0], m)
		} else {
			r = m; continue
		}
	}
	return nil
}

func search(word string) []pair {
	word = DictionaryGenerator.Parse(word)
	if word == "" { return nil }
	var hCode [20]byte
	hashCode(word, &hCode)
	cn := cacheNode{hCode:hCode}
	c := cache.Get(cn)
	if c == nil {
		cn.pairs = searchFromFile(hCode)
		if cn.pairs == nil {
			return nil
		}
		cache.Insert(cn)
		if cache.Len() > 1024 {
			cache.DeleteMin()
		}
		return cn.pairs
	}
	return c.(cacheNode).pairs
}