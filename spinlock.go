package sync

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type spinLock uint32

func (t *spinLock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(t), 0, 1) {
		runtime.Gosched()
	}
}

func (t *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(t), 0)
}

func NewSpinLock() sync.Locker {
	return new(spinLock)
}
