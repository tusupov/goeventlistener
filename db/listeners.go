package db

// ListenerRequest
// With param event name, listener name and http address
type ListenerRequest struct {
	Event    string `json:"event"`
	Listener string `json:"name"`
	Address  string `json:"address"`
}

// Listener
// With listener name and http address
type Listener struct {
	Name    string
	Address string
}

// ListenerList is like a Go map[string]Listener but is safe for concurrent use
// by multiple goroutines without additional locking or coordination.
type ListenerList struct {
	locker Locker
	list   map[string]Listener
}

func NewListenerList() *ListenerList {
	return &ListenerList{
		locker: Locker{},
		list:   make(map[string]Listener),
	}
}

// Listeners get all listeners, copy to slice and return
func (list *ListenerList) Listeners() []Listener {

	list.locker.RLock()
	listener := make([]Listener, len(list.list))
	pos := 0
	for _, l := range list.list {
		listener[pos] = l
		pos++
	}
	list.locker.RUnlock()

	return listener
}

// Store add listener
// If listener by listener name has, return false
func (list *ListenerList) Store(l Listener) bool {

	// Check exists
	list.locker.RLock()
	_, ok := list.list[l.Name]
	list.locker.RUnlock()
	if ok {
		return false
	}

	// Check and add
	list.locker.Lock()
	_, ok = list.list[l.Name]
	if ok {
		list.locker.Unlock()
		return false
	}
	list.list[l.Name] = l
	list.locker.Unlock()

	return true
}

// Delete listener by name
// if listener by name not exists, return false
func (list *ListenerList) Delete(name string) (ok bool) {

	// Check exists
	list.locker.RLock()
	_, ok = list.list[name]
	list.locker.RUnlock()
	if !ok {
		return
	}

	// Check and delete
	list.locker.Lock()
	_, ok = list.list[name]
	if !ok {
		list.locker.Unlock()
		return
	}
	delete(list.list, name)
	list.locker.Unlock()

	return
}
