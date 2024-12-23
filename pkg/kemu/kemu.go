// Package kemu provides locking per-key.
// For example, you can acquire a lock for a specific user ID and all other requests for that user ID
// will block until that entry is unlocked (effectively your work load will be run serially per-user ID),
// and yet have work for separate user IDs happen concurrently.
package kemu

import (
	"fmt"
	"sync"
)

// Mutex wraps a map of mutexes
// and each key will be locked seperately.
type Mutex struct {
	m       sync.Mutex
	mutexes map[interface{}]*entry
}

// mutexEntry hold each entry of keyed mutex.
type entry struct {
	mutex *Mutex
	em    sync.Mutex
	count int
	key   interface{}
}

// Unlocker defines an Unlock method to release the lock.
type Unlocker interface {
	Unlock() error
}

// New create new instance of KeyedMutex.
func New() *Mutex {
	return &Mutex{mutexes: make(map[interface{}]*entry)}
}

// Lock acquires a lock corresponding to the given key.
// This method will never return nil and Unlock method must be called to release the lock.
func (m *Mutex) Lock(key interface{}) Unlocker {
	// Read and create mutexes entry for given key atomically
	m.m.Lock()
	e, f := m.mutexes[key]
	if !f {
		e = &entry{mutex: m, key: key}
		m.mutexes[key] = e
	}

	// Increase mutex entry count
	// and unlock parent mutex
	e.count++
	m.m.Unlock()

	// Acquire lock and blocking
	// until entry count equals with setted value
	e.em.Lock()
	return e
}

// Unlock will release the lock.
func (e *entry) Unlock() error {
	m := e.mutex

	// Do decrement entry count
	// and if necessary remove mutexes entry atomically
	m.m.Lock()
	e, f := m.mutexes[e.key]
	if !f {
		m.m.Unlock()
		return fmt.Errorf("mutex entry for key:%v not found, but unlock requested", e.key)
	}
	e.count--
	if e.count < 1 {
		delete(m.mutexes, e.key)
	}
	m.m.Unlock()

	// Unlocking mutex entry
	e.em.Unlock()
	return nil
}

// List return list of mutexes.
func (m *Mutex) List() map[interface{}]*entry {
	m.m.Lock()
	defer m.m.Unlock()

	return m.mutexes
}
