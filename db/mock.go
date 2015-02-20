package db

import (
	"io"
	"os"
	"time"

	"github.com/steveyen/gkvlite"
)

type mockFile struct {
	b   []byte
	mod time.Time
}

// NewMockFile creates an empty gkvlite.StoreFile contained in memory.
func NewMockFile() gkvlite.StoreFile {
	return &mockFile{mod: time.Now()}
}

func (f *mockFile) ReadAt(p []byte, off int64) (n int, err error) {
	if off > int64(len(f.b)) {
		return 0, io.EOF
	}
	n = copy(p, f.b[off:])
	if n != len(p) {
		err = io.EOF
	}
	return
}

func (f *mockFile) WriteAt(p []byte, off int64) (n int, err error) {
	if int64(len(f.b)) < off {
		f.Truncate(off)
	}
	if int64(len(f.b)) < off+int64(len(p)) {
		f.b = append(f.b[:off], p...)
	} else {
		copy(f.b[off:], p)
	}
	f.mod = time.Now()

	return len(p), nil
}

func (f *mockFile) Stat() (fi os.FileInfo, err error) {
	return mockStat{
		size: int64(len(f.b)),
		mod:  f.mod,
	}, nil
}

func (f *mockFile) Truncate(size int64) error {
	if size > int64(len(f.b)) {
		f.b = append(f.b, make([]byte, size-int64(len(f.b)))...)
	} else {
		f.b = f.b[:size]
	}
	f.mod = time.Now()
	return nil
}

type mockStat struct {
	size int64
	mod  time.Time
}

func (fi mockStat) Name() string       { return "mock.db" }
func (fi mockStat) Size() int64        { return fi.size }
func (fi mockStat) Mode() os.FileMode  { return 0666 }
func (fi mockStat) ModTime() time.Time { return fi.mod }
func (fi mockStat) IsDir() bool        { return fi.Mode().IsDir() }
func (fi mockStat) Sys() interface{}   { return nil }
