package main

import (
	"reflect"
	"unsafe"
	"bufio"
	"errors"
)

var hCodeDebug [20]byte

func init() {
	hashCode("sonnet", &hCodeDebug)
}

type indexEntry struct {
	hashCode [20]byte
	invertFileId int
	invertFileOff int
}

type ReadClosableMMFile struct {
	mmFile
	writer *bufio.Writer
}

func (rcmmFile *ReadClosableMMFile) Write(b []byte) (n int, e error) {
	if n, e = rcmmFile.mmFile.Write(b); e != nil {
		an, e := rcmmFile.Append(b[n:])
		return an + n, e
	}
	return
}

func (rcmmFile *ReadClosableMMFile) Close() {
	if rcmmFile.writer != nil {
		rcmmFile.writer.Flush()
		rcmmFile.writer = nil
	}
}

func newReadClosableMMFile(mf *mmFile) *ReadClosableMMFile {
	ret := &ReadClosableMMFile{
		mmFile: *mf,
	}
	ret.writer = bufio.NewWriter(ret)
	return ret
}

type indexFile struct {
	*ReadClosableMMFile
}

func (inf *indexFile) addIndex(hCode *[20]byte, invOff uint32) (err error) {
	b := reflect.SliceHeader{}
	b.Cap = 20
	b.Len = b.Cap
	b.Data = uintptr(unsafe.Pointer(&(*hCode)[0]))
	n, e := inf.writer.Write(*(*[]byte)(unsafe.Pointer(&b)))
	if n != 20 || e != nil {
		return e
	}
	b.Cap = 4
	b.Len = 4
	b.Data = uintptr(unsafe.Pointer(&invOff))
	n, e = inf.writer.Write(*(*[]byte)(unsafe.Pointer(&b)))
	if n != 4 || e != nil {
		return e
	}
	return nil
}

var index [256]indexFile

type invertFile struct {
	*ReadClosableMMFile
}

func (inv *invertFile) addPair(pairs []pair) (int, error) {
	if pairs == nil || len(pairs) == 0 {
		return 0, errors.New("No pair to add.")
	}
	off := inv.ReadClosableMMFile.writer.Buffered() + inv.ReadClosableMMFile.idx
	var length = uint64(len(pairs))
	b := reflect.SliceHeader{}
	b.Cap = 8
	b.Len = 8
	b.Data = uintptr(unsafe.Pointer(&length))
	if n, e := inv.writer.Write(*(*[]byte)(unsafe.Pointer(&b)));
		n == b.Len && e == nil {
	} else {
		return 0, e
	}
	b.Cap = len(pairs) * 8
	b.Len = b.Cap
	b.Data = uintptr(unsafe.Pointer(&pairs[0].rawFId))
	if n, e := inv.writer.Write(*(*[]byte)(unsafe.Pointer(&b)));
		n == b.Len && e == nil {
		return off, nil
	} else {
		return 0, e
	}
}

var invert [256]invertFile