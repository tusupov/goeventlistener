package db

type Events struct {
	locker Locker
	list   map[string]*ListenerList
}

func NewEvents() *Events {
	return &Events{
		locker: Locker{},
		list:   make(map[string]*ListenerList),
	}
}

// Load listener list
func (e *Events) Load(event string) (list *ListenerList, loaded bool) {
	e.locker.RLock()
	list, loaded = e.list[event]
	e.locker.RUnlock()
	return
}

// Add listener list
// If event not exists, create new event, and add listener to list
func (e *Events) Store(event string, listener Listener) bool {

	var list *ListenerList

	list, loaded := e.Load(event)
	if !loaded {

		var newLoaded bool

		e.locker.Lock()
		list, newLoaded = e.list[event]
		if !newLoaded {
			list = NewListenerList()
			e.list[event] = list
		}
		e.locker.Unlock()

	}

	return list.Store(listener)

}

// Delete listener from event
func (e *Events) DeleteListener(event string, listenerName string) (loaded bool) {

	// Load listener list
	e.locker.RLock()
	list, loaded := e.list[event]
	e.locker.RUnlock()
	if !loaded {
		return
	}

	// Delete listener from list
	return list.Delete(listenerName)

}
