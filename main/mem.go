package main

import (
	"errors"
	"strconv"
	"crypto/sha1"
	"reflect"
	"unsafe"
	"sort"
	"io"
	"fmt"
)

func hashCode(s string, code *[20]byte) {
	h := sha1.New()
	h.Write([]byte(s))
	sum := h.Sum(nil)
	var c reflect.SliceHeader
	c.Cap = 20
	c.Len = 20
	c.Data = uintptr(unsafe.Pointer(&(*code)[0]))
	copy(*(*[]byte)(unsafe.Pointer(&c)), sum)
}

type pair struct {
	rawFId uint32
	rawOff uint32
}

var (
	invertSegment map[[20]byte][]pair
	currentHash   byte
	indexSegment map[[20]byte]uint32

	NOT_CURRENT_HASH = errors.New("The first byte of hash code is not running yet.")
)

func memAddWord(word string, rawFId, rawOff uint32) error {
	var hCode [20]byte
	hashCode(word, &hCode)
	if hCode[0] != currentHash {
		return NOT_CURRENT_HASH
	}
	p, ok := invertSegment[hCode]
	if ok {
		invertSegment[hCode] = append(p, pair{rawFId, rawOff})
	} else {
		invertSegment[hCode] = []pair{{rawFId, rawOff}}
	}
	return nil
}

type hashCodeSlice [][20]byte

func (hcs hashCodeSlice) Len() int {
	return len(hcs)
}

func (hcs hashCodeSlice) Less(i, j int) bool {
	for t := 0; t < 5; t++ {
		if *(*uint32)(unsafe.Pointer(&(hcs)[i][4*t])) <
			*(*uint32)(unsafe.Pointer(&(hcs)[j][4*t])) {
			return false
		} else if *(*uint32)(unsafe.Pointer(&(hcs)[i][4*t])) >
			*(*uint32)(unsafe.Pointer(&(hcs)[j][4*t])) {
		 	return true } else { continue }
	}
	return true
}

func (hcs hashCodeSlice) Swap(i, j int) {
	t := (hcs)[i]
	(hcs)[i] = (hcs)[j]
	(hcs)[j] = t
}

func memFlush() (e error) {
	var mf *mmFile
	mf, e = newMMFile("./invert/" + strconv.Itoa(int(currentHash)) + ".inv", 1024)
	if e != nil { return }
	invert[currentHash].ReadClosableMMFile = newReadClosableMMFile(mf)
	keys := make(hashCodeSlice, 0)
	for hCode, pairs := range invertSegment {
		off, e := invert[currentHash].addPair(pairs)
		if e != nil { return e }
		indexSegment[hCode] = uint32(off)
		keys = append(keys, hCode)
	}
	invert[currentHash].Close()

	sort.Sort(&keys)
	mf, e = newMMFile("./index/" + strconv.Itoa(int(currentHash)) + ".inv", 1024)
	if e != nil { return }
	index[currentHash].ReadClosableMMFile = newReadClosableMMFile(mf)
	for i := 0; i < len(keys); i++ {
		e = index[currentHash].addIndex(&keys[i], indexSegment[keys[i]])
		if e != nil { return }
	}
	index[currentHash].Close()
	invertSegment = make(map[[20]byte][]pair)
	indexSegment = make(map[[20]byte]uint32)
	return
}

// create the index file and invert file
// readWord is a function that read a word from a file identified by rawFId,
// returns a word and its offset or EOF error. When it returns EOF,
// it reset its offset so it returns the first word next call.
// rawFNum is the number of raw files.
func createFiles(readWord func(rawFId int) (word string, off uint32, err error), rawFNum int) error {
	invertSegment = make(map[[20]byte][]pair)
	indexSegment = make(map[[20]byte]uint32)
	start := false
	fmt.Print("\n")
	for currentHash = byte(*m); currentHash <= 255; currentHash += byte(*p) {
		if currentHash == 0 {
			if start {
				return nil
			} else {
				start = true
			}
		}
		for id := 0; id < rawFNum; id++ {
			for w, o, e := readWord(id); e != io.EOF; w, o, e = readWord(id) {
				if err := memAddWord(w, uint32(id), o); err != nil {
					continue
				}
			}
		}
		if err := memFlush(); err != nil {
			return err
		}
		fmt.Printf("\r%3d/%3d\n", currentHash, 255)
	}
	return nil
}