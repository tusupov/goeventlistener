package db

import (
	"runtime"
	"sync/atomic"
)

var lockerMaxReaders = int32(runtime.NumCPU() * 2)

// A Locker is a reader/writer mutual exclusion lock.
// Effective when goroutines count equal
// The zero value for a locker is an unlocked mutex.
// Greater than zero for a locker mean lock for read,
// -1 value for a locker mean lock for write.
type Locker struct {
	locker int32
}

// RLock
// Increment locker
func (l *Locker) RLock() {
	for {
		locker := atomic.LoadInt32(&l.locker)
		if locker >= 0 && locker <= lockerMaxReaders && atomic.CompareAndSwapInt32(&l.locker, locker, locker+1) {
			return
		} else {
			runtime.Gosched()
		}
	}
}

// RUnlock
// Decrement locker
func (l *Locker) RUnlock() {
	atomic.AddInt32(&l.locker, -1)
}

// Unlock
// Change locker from 0 to -1
func (l *Locker) Lock() {
	for {
		if l.locker == 0 && atomic.CompareAndSwapInt32(&l.locker, 0, -1) {
			return
		} else {
			runtime.Gosched()
		}
	}
}

// Unlock
// Change locker from -1 to 0
func (l *Locker) Unlock() {
	if l.locker != -1 {
		panic("db: Unlock of unlocked Locker")
	}
	atomic.StoreInt32(&l.locker, 0)
}
