package db

// Index is like a Go map[string]string but is safe for concurrent use
// by multiple goroutines without additional locking or coordination.
type Index struct {
	locker Locker
	list   map[string]string
}

func NewIndex() *Index {
	return &Index{
		list: make(map[string]string),
	}
}

// Load returns the value stored in the map for a key,
// or empty string if no value is present.
// The ok result indicates whether value was found in the map.
func (i *Index) Load(key string) (value string, ok bool) {
	i.locker.RLock()
	value, ok = i.list[key]
	i.locker.RUnlock()
	return
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value.
// The loaded result is true if the value was loaded, false if stored.
func (i *Index) LoadOrStore(key, value string) (actual string, loaded bool) {

	actual, loaded = i.Load(key)
	if loaded {
		return
	}

	i.locker.Lock()
	actual, loaded = i.list[key]
	if loaded {
		i.locker.Unlock()
		return
	}
	i.list[key] = value
	i.locker.Unlock()

	return value, false

}

// Store sets the value for a key.
func (i *Index) Store(key, value string) {
	i.locker.Lock()
	i.list[key] = value
	i.locker.Unlock()
}

// Delete deletes the value for a key.
// The result return true if value was found in the map and deleted.
func (i *Index) Delete(key string) bool {

	_, loaded := i.Load(key)
	if !loaded {
		return false
	}

	i.locker.Lock()
	_, loaded = i.list[key]
	if !loaded {
		i.locker.Unlock()
		return false
	}
	delete(i.list, key)
	i.locker.Unlock()

	return true

}
