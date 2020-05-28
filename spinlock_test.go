package sync

import (
	"sync"
	"testing"
)

type lock1 struct {
	//sl spinLock
	sl sync.Locker
}

func TestSpinLock(t *testing.T) {
	i := 0

	wg := &sync.WaitGroup{}
	wg.Add(2)

	l := &lock1{sl: NewChanLock()}
	go func() {
		l.sl.Lock()
		defer l.sl.Unlock()

		for _i := 0; _i < 10000; _i++ {
			i += 1
		}

		wg.Done()
	}()

	go func() {
		l.sl.Lock()
		defer l.sl.Unlock()

		for _i := 0; _i < 10000; _i++ {
			i -= 1
		}

		wg.Done()
	}()

	wg.Wait()
	t.Log(i)
}
