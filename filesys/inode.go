package filesys

import (
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

func (im *inodeManager) get(prefix, key string) (uint64, bool) {
	fullKey := prefix + ":" + key

	// check existing
	im.mtx.RLock()
	id, ok := im.inodes[fullKey]
	im.mtx.RUnlock()

	return id, ok
}

func (im *inodeManager) getOrCreate(prefix, key string) uint64 {
	id, ok := im.get(prefix, key)
	if ok {
		return id
	}

	fullKey := prefix + ":" + key
	// generate new
	im.mtx.Lock()

	id = im.counter
	im.counter++
	im.inodes[fullKey] = id

	im.mtx.Unlock()

	return id
}
