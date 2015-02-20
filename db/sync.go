package db

import (
	"os"
	"sync"

	"github.com/steveyen/gkvlite"
)

type syncFile struct {
	store gkvlite.StoreFile
	mtx   sync.RWMutex
}

// NewSyncFile creates a wrapper around the given gkvlite.StoreFile that
// only allows one write/truncate or any number of read/stat operations to
// run at the same time. A waiting write/truncate will block new read/stat
// operations from starting.
func NewSyncFile(store gkvlite.StoreFile) gkvlite.StoreFile {
	return &syncFile{store: store}
}

func (f *syncFile) ReadAt(p []byte, off int64) (n int, err error) {
	f.mtx.RLock()
	defer f.mtx.RUnlock()

	return f.store.ReadAt(p, off)
}

func (f *syncFile) WriteAt(p []byte, off int64) (n int, err error) {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.store.WriteAt(p, off)
}

func (f *syncFile) Stat() (fi os.FileInfo, err error) {
	f.mtx.RLock()
	defer f.mtx.RUnlock()

	return f.store.Stat()
}

func (f *syncFile) Truncate(size int64) error {
	f.mtx.Lock()
	defer f.mtx.Unlock()

	return f.store.Truncate(size)
}
