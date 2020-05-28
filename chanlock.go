package sync

import (
	"sync"
)

type chanLock struct {
	c chan struct{}
}

func (t *chanLock) Lock() {
	t.c <- struct{}{}
}

func (t *chanLock) Unlock() {
	<-t.c
}

func NewChanLock() sync.Locker {
	return &chanLock{c: make(chan struct{}, 1)}
}
