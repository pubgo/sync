package sync

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

type spinLock uint32

func (t *spinLock) TryLock(dur time.Duration) error {
	tk := time.Tick(dur)
	for {
		select {
		case <-tk:
			return errors.New("timeout")
		default:
			if atomic.CompareAndSwapUint32((*uint32)(t), 0, 1) {
				return nil
			}
			runtime.Gosched()
		}
	}
}

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
