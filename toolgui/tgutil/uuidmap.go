package tgutil

import (
	"sync"
	"time"

	"github.com/google/uuid"
)

// UUIDMap provide a goroutine-safe mapping from UUID to T.
type UUIDMap[T any] interface {
	// New create a (id -> T) mapping and return id.
	New() string

	// Get return T by id, return nil if id does not exist.
	Get(id string) *T

	// Del delete uuid.
	Del(id string)

	// Size return the number of T.
	Size() int

	// Destroy the resource hold by the State.
	Destroy()
}

type dataPair[T any] struct {
	value     *T
	timestamp time.Time
}

type uuidmap[T any] struct {
	lock        sync.RWMutex
	data        map[string]*dataPair[T]
	constructor func() *T
	destructor  func(*T)

	ttl           time.Duration
	latestCleanup time.Time
}

// NewUUIDMap create T by providing the constructor and destructor the of T.
func NewUUIDMap[T any](
	constructor func() *T, destructor func(*T), ttl time.Duration) UUIDMap[T] {

	return &uuidmap[T]{
		data:          make(map[string]*dataPair[T]),
		constructor:   constructor,
		destructor:    destructor,
		ttl:           ttl,
		latestCleanup: time.Now(),
	}
}

// Destroy release the resources hold by data.
func (ss *uuidmap[T]) Destroy() {
	for _, d := range ss.data {
		ss.destructor(d.value)
	}
}

// Size return the number of item.
func (ss *uuidmap[T]) Size() int {
	return len(ss.data)
}

func (ss *uuidmap[T]) cleanup() {
	// TBD: how to limit the number of T?
	if time.Since(ss.latestCleanup) < 20*ss.ttl {
		return
	}

	ids := []string{}
	for id, d := range ss.data {
		if time.Since(d.timestamp) > ss.ttl {
			ids = append(ids, id)
		}
	}
	for _, id := range ids {
		delete(ss.data, id)
	}
}

// New create a (id -> T) mapping and return id.
func (ss *uuidmap[T]) New() string {
	id := uuid.New().String()
	ss.lock.Lock()
	defer ss.lock.Unlock()
	ss.cleanup()
	ss.data[id] = &dataPair[T]{
		value:     ss.constructor(),
		timestamp: time.Now(),
	}
	return id
}

// Get return T by id, return nil if id does not exist.
func (ss *uuidmap[T]) Get(id string) *T {
	ss.lock.RLock()
	defer ss.lock.RUnlock()

	d, ok := ss.data[id]
	if !ok {
		return nil
	}

	d.timestamp = time.Now()
	return d.value
}

// Del delete uuid.
func (ss *uuidmap[T]) Del(id string) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	d, ok := ss.data[id]
	if !ok {
		return
	}

	ss.destructor(d.value)
	delete(ss.data, id)
}
