package filesys

import (
	"strings"
	"sync"
)

var (
	inoMgr *inodeManager
)

func init() {
	inoMgr = newInodeManager()
}

type inodeManager struct {
	mtx     sync.RWMutex
	inodes  map[string]uint64
	counter uint64
}

func newInodeManager() *inodeManager {
	return &inodeManager{
		inodes:  make(map[string]uint64),
		counter: 2,
	}
}

func (im *inodeManager) get(keys ...string) (uint64, bool) {
	fullKey := strings.Join(keys, ":")

	// check existing
	im.mtx.RLock()
	id, ok := im.inodes[fullKey]
	im.mtx.RUnlock()

	return id, ok
}

func (im *inodeManager) getOrCreate(keys ...string) uint64 {
	id, ok := im.get(keys...)
	if ok {
		return id
	}

	fullKey := strings.Join(keys, ":")
	// generate new
	im.mtx.Lock()

	id = im.counter
	im.counter++
	im.inodes[fullKey] = id

	im.mtx.Unlock()

	return id
}
