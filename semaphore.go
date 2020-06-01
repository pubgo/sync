package sync

import "sync"

type semaphore struct {
	// Number of Acquires - Releases. When this goes to zero, this structure is removed from map.
	// Only updated inside blockingKeyCountLimit.lk lock.
	refs int

	max   int
	value int
	wait  sync.Cond
}

func newSemaphore(max int) *semaphore {
	return &semaphore{
		max:  max,
		wait: sync.Cond{L: new(sync.Mutex)},
	}
}

func (s *semaphore) Running() int {
	s.wait.L.Lock()
	defer s.wait.L.Unlock()

	return s.value
}

func (s *semaphore) Acquire() {
	s.wait.L.Lock()
	defer s.wait.L.Unlock()
	for {
		if s.value+1 <= s.max {
			s.value++
			return
		}
		s.wait.Wait()
	}
}

func (s *semaphore) Release() {
	s.wait.L.Lock()
	defer s.wait.L.Unlock()
	s.value--
	if s.value < 0 {
		panic("semaphore Release without Acquire")
	}
	s.wait.Signal()
}

type blockingKeyCountLimit struct {
	mutex   sync.RWMutex
	current map[string]*semaphore
	limit   int
}

func NewBlockingKeyCountLimit(n int) *blockingKeyCountLimit {
	return &blockingKeyCountLimit{current: make(map[string]*semaphore), limit: n}
}

func (l *blockingKeyCountLimit) Running() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	all := 0
	for _, v := range l.current {
		all += v.Running()
	}
	return all
}

func (l *blockingKeyCountLimit) Acquire(key2 []byte) {
	key := string(key2)

	l.mutex.Lock()
	kl, ok := l.current[key]
	if !ok {
		kl = newSemaphore(l.limit)
		l.current[key] = kl

	}
	kl.refs++
	l.mutex.Unlock()

	kl.Acquire()
}

func (l *blockingKeyCountLimit) Release(key2 []byte) {
	key := string(key2)

	l.mutex.Lock()
	kl, ok := l.current[key]
	if !ok {
		panic("key not in map. Possible reason: Release without Acquire.")
	}
	kl.refs--
	if kl.refs < 0 {
		panic("internal error: refs < 0")
	}
	if kl.refs == 0 {
		delete(l.current, key)
	}
	l.mutex.Unlock()

	kl.Release()
}
