package main

import (
	"github.com/edsrzf/mmap-go"
	"os"
	"io"
	"errors"
)

var MUST_APPEND_END = errors.New("mmio: only can append to a end of a file.")

const MAX_AUTO_APPEND_SIZE = 128*1024*1024

type mmFile struct {
	filename string
	f *os.File
	m mmap.MMap
	idx int
}

func newMMFile(filename string, resize int64) (mf *mmFile, e error) {
	mf = &mmFile{}
	mf.filename = filename
	mf.f, e = os.OpenFile(filename, os.O_CREATE | os.O_RDWR | os.O_TRUNC, 0666)
	if e != nil { return }
	if resize > 0 {
		e = mf.f.Truncate(resize)
	}
	if e != nil { return }
	mf.m, e = mmap.Map(mf.f, mmap.RDWR, 0)
	if e != nil { return }
	return
}

func openMMFile(filename string) (mf *mmFile, e error) {
	mf = &mmFile{}
	mf.filename = filename
	mf.f, e = os.OpenFile(filename, os.O_RDONLY, 0666)
	if e != nil { return }
	mf.m, e = mmap.Map(mf.f, mmap.RDONLY, 0)
	if e != nil { return }
	return
}

func (mf *mmFile) Read(b []byte) (n int, e error) {
	if mf.idx + len(b) > len(mf.m) {
		n = len(mf.m) - mf.idx
		e = io.EOF
	} else {
		n = len(b)
	}
	if n > 0 {
		copy(b, mf.m[mf.idx:])
	}
	mf.idx += n
	return
}

func (mf *mmFile) ReadAt(b []byte, off int64) (n int, e error) {
	if int(off) + len(b) > len(mf.m) {
		n = len(mf.m) - mf.idx
		e = io.EOF
	} else {
		n = len(b)
	}
	if n > 0 {
		copy(b, mf.m[off:])
	}
	return
}

func (mf *mmFile) Write(b []byte) (n int, e error) {
	if mf.idx + len(b) > len(mf.m) {
		n = len(mf.m) - mf.idx
		e = io.EOF
	} else {
		n = len(b)
	}
	if n > 0 {
		copy(mf.m[mf.idx:], b)
	}
	mf.idx += n
	return
}

func (mf *mmFile) WriteAt(b []byte, off int64) (n int, e error) {
	if int(off) + len(b) > len(mf.m) {
		n = len(mf.m) - mf.idx
		e = io.EOF
	} else {
		n = len(b)
	}
	if n > 0 {
		copy(mf.m[off:], b)
	}
	return
}

func (mf *mmFile) Close() (e error) {
	e = mf.m.Unmap()
	if e != nil { return }
	e = mf.f.Close()
	mf.m = nil
	return
}

func (mf *mmFile) Lock() error {
	return mf.m.Lock()
}

func (mf *mmFile) UnLock() error {
	return mf.m.Unmap()
}

func (mf *mmFile) Flush() error {
	return mf.m.Flush()
}

func (mf *mmFile) Len() int {
	return len(mf.m)
}

func (mf *mmFile) Truncate(size int64) (e error) {
	e = mf.m.Unmap()
	if e != nil { return }
	e = mf.f.Truncate(size)
	if e != nil { return }
	mf.m, e = mmap.Map(mf.f, mmap.RDWR, 0)
	if e != nil { return }
	return
}

func (mf *mmFile) Seek(off int64, whence int) (ret int64, e error) {
	if mf.idx + int(off) > mf.Len() {
		return int64(mf.idx), io.EOF
	}
	ret, e = mf.f.Seek(off, whence)
	if e != nil { return }
	switch whence {
	case 0:
		mf.idx = int(off)
	case 1:
		mf.idx += int(off)
	case 2:
		mf.idx += int(off)
	}
	return
}

func (mf *mmFile) Append(b []byte) (n int, e error) {
	if mf.idx != mf.Len() {
		return 0, MUST_APPEND_END
	}
	if len(b) > MAX_AUTO_APPEND_SIZE {
		e = mf.Truncate(int64(mf.Len() + len(b)))
	} else
	if mf.Len() < MAX_AUTO_APPEND_SIZE {
		newLen := mf.Len()
		for newLen - mf.idx < len(b) {
			newLen <<= 1
		}
		e = mf.Truncate(int64(newLen))
	} else {
		e = mf.Truncate(int64(mf.Len() + MAX_AUTO_APPEND_SIZE))
	}
	if e != nil { return }
	n, e = mf.Write(b)
	return
}