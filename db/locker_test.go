package db

import (
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkSyncMap_Add(b *testing.B) {

	b.Run("self", func(b *testing.B) {

		m := sync.Map{}
		x := int32(0)

		for i := 0; i < b.N; i++ {
			newI := atomic.AddInt32(&x, 1)
			v := "listener" + strconv.Itoa(int(newI))
			m.LoadOrStore(v, v)
		}

	})

	b.Run("paralel", func(b *testing.B) {

		m := sync.Map{}
		x := int32(0)

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				newI := atomic.AddInt32(&x, 1)
				v := "listener" + strconv.Itoa(int(newI))
				m.LoadOrStore(v, v)
			}
		})

	})

	b.Run("paralelwg", func(b *testing.B) {

		m := sync.Map{}
		x := int32(0)
		wg := sync.WaitGroup{}

		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				newI := atomic.AddInt32(&x, 1)
				v := "listener" + strconv.Itoa(int(newI))
				m.LoadOrStore(v, v)
				wg.Done()
			}()
		}

		wg.Wait()

	})

}

func BenchmarkMutex_Add(b *testing.B) {

	b.Run("self", func(b *testing.B) {

		m := make(map[string]string)
		mu := sync.RWMutex{}
		x := int32(0)

		for i := 0; i < b.N; i++ {
			newI := atomic.AddInt32(&x, 1)
			v := "listener" + strconv.Itoa(int(newI))

			// Store if not exists
			mu.RLock()
			_, ok := m[v]
			mu.RUnlock()
			if ok {
				continue
			}

			mu.Lock()
			_, ok = m[v]
			if ok {
				mu.Unlock()
				continue
			}
			m[v] = v
			mu.Unlock()

		}

	})

	b.Run("paralel", func(b *testing.B) {

		m := make(map[string]string)
		mu := sync.RWMutex{}
		x := int32(0)

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				newI := atomic.AddInt32(&x, 1)
				v := "listener" + strconv.Itoa(int(newI))

				// Store if not exists
				mu.RLock()
				_, ok := m[v]
				mu.RUnlock()
				if ok {
					continue
				}

				mu.Lock()
				_, ok = m[v]
				if ok {
					mu.Unlock()
					continue
				}
				m[v] = v
				mu.Unlock()

			}
		})

	})

	b.Run("paralelwg", func(b *testing.B) {

		m := make(map[string]string)
		mu := sync.RWMutex{}
		x := int32(0)
		wg := sync.WaitGroup{}

		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				newI := atomic.AddInt32(&x, 1)
				v := "listener" + strconv.Itoa(int(newI))

				// Store if not exists
				mu.RLock()
				_, ok := m[v]
				mu.RUnlock()
				if ok {
					return
				}

				mu.Lock()
				_, ok = m[v]
				if ok {
					mu.Unlock()
					return
				}
				m[v] = v
				mu.Unlock()

				wg.Done()
			}()
		}

		wg.Wait()

	})

}

func BenchmarkLocker_Add(b *testing.B) {

	b.Run("self", func(b *testing.B) {

		m := NewIndex()
		x := int32(0)

		for i := 0; i < b.N; i++ {
			newI := atomic.AddInt32(&x, 1)
			v := "listener" + strconv.Itoa(int(newI))
			m.LoadOrStore(v, v)
		}

	})

	b.Run("paralel", func(b *testing.B) {

		m := NewIndex()
		x := int32(0)

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				newI := atomic.AddInt32(&x, 1)
				v := "listener" + strconv.Itoa(int(newI))
				m.LoadOrStore(v, v)
			}
		})

	})

	b.Run("paralelwg", func(b *testing.B) {

		m := NewIndex()
		x := int32(0)
		wg := sync.WaitGroup{}

		for i := 0; i < b.N; i++ {
			wg.Add(1)
			go func() {
				newI := atomic.AddInt32(&x, 1)
				v := "listener" + strconv.Itoa(int(newI))
				m.LoadOrStore(v, v)
				wg.Done()
			}()
		}

		wg.Wait()

	})

}
