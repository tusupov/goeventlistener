package db

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"

	"golang.org/x/net/context"
	"golang.org/x/net/context/ctxhttp"
)

const publishWorkers = 50

var (

	// Http request workers count
	workersCh = make(chan struct{}, publishWorkers)

	// Error codes returned by failures
	errListenerExists          = errors.New("Listener exists.")
	errListenerNotExists       = errors.New("Listener not exists.")
	errEventNotExists          = errors.New("Event not exists.")
	errEventListenersNotExists = errors.New("Event listeners not exists.")
)

// Local storage
//
// Events save all listeners by event name
// events is like Go map[string]map[string][]struct
//
// Index save event name by listener name
// For fast search listener, used index
// index is like Go map[string]string
type Storage struct {
	locker Locker
	events *Events
	index  *Index
}

func New() *Storage {
	return &Storage{
		locker: Locker{},
		events: NewEvents(),
		index:  NewIndex(),
	}
}

// Add listener by event name
func (s *Storage) Add(r ListenerRequest) error {

	s.locker.RLock()
	_, loaded := s.index.LoadOrStore(r.Listener, r.Event)
	if loaded {
		s.locker.RUnlock()
		return errListenerExists
	}

	s.events.Store(r.Event, Listener{r.Listener, r.Address})
	s.locker.RUnlock()

	return nil
}

// Delete listener by name
func (s *Storage) DeleteListener(listenerName string) error {

	eventName, loaded := s.index.Load(listenerName)
	if !loaded {
		return errListenerNotExists
	}

	s.locker.Lock()
	deleted := s.events.DeleteListener(eventName, listenerName)
	if !deleted {
		s.locker.Unlock()
		return errListenerNotExists
	}
	s.index.Delete(listenerName)
	s.locker.Unlock()

	return nil

}

// Get listeners list and do http request
func (s *Storage) Publish(cntx context.Context, eventName string) error {

	// Listeners list for eventName
	s.locker.RLock()
	listenerList, loaded := s.events.Load(eventName)
	if !loaded {
		s.locker.RUnlock()
		return errEventNotExists
	}

	// Listeners to slice
	listeners := listenerList.Listeners()
	s.locker.RUnlock()

	if len(listeners) == 0 {
		return errEventListenersNotExists
	}

	// Start do http request
	return s.publishDo(cntx, listeners)

}

// Do http request for all event listeners
// Used goroutines for all request, for quickly get all results
func (s *Storage) publishDo(cntx context.Context, listener []Listener) (err error) {

	wg := sync.WaitGroup{}
	errCnt := int32(0)

	// Listeners list
	for _, l := range listener {

		// Has error
		if err != nil {
			break
		}

		// Do http request
		workersCh <- struct{}{}
		wg.Add(1)

		go func(l Listener) {

			defer func() {
				wg.Done()
				<-workersCh
			}()

			resp, errGet := ctxhttp.Head(cntx, nil, l.Address)
			if errGet != nil {
				if atomic.CompareAndSwapInt32(&errCnt, 0, 1) {
					err = errGet
				}
				return
			}

			if resp.StatusCode != http.StatusOK {
				err = fmt.Errorf("Expect status code %d, but %d, Url: %s", http.StatusOK, resp.StatusCode, l.Address)
				return
			}

		}(l)

	}

	// Wait all request
	wg.Wait()

	return

}
