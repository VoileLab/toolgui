package tgexec

import (
	"errors"
	"sync"
	"time"
	"uuid"

	"github.com/voilelab/toolgui/toolgui/tgutil"
)

// errUUIDMapFull is the error [uuidMap.New] returns when the map already
// holds as many items as it may.
var errUUIDMapFull = errors.New("uuidmap is full")

// errUUIDNotFound is the error [uuidMap.Acquire] returns when the map holds
// nothing under the id.
var errUUIDNotFound = errors.New("uuidmap: id not found")

// errUUIDAlive is the error [uuidMap.Acquire] returns when something else
// already holds the id.
var errUUIDAlive = errors.New("uuidmap: id already alive")

type dataPair[T any] struct {
	value     *T
	timestamp time.Time
	alive     bool
}

// uuidMap is a goroutine-safe mapping from UUID to T.
type uuidMap[T any] struct {
	lock        sync.RWMutex
	data        map[string]*dataPair[T]
	constructor func() *T
	destructor  func(*T)

	ttl           time.Duration
	latestCleanup time.Time

	// maxSize is the number of items the map takes, 0 for no limit.
	maxSize int
}

// newUUIDMap create T by providing the constructor and destructor the of T.
// An item nothing has kept alive for ttl is reclaimed. The map takes any
// number of items until [uuidMap.SetMaxSize] gives it a limit.
func newUUIDMap[T any](
	constructor func() *T, destructor func(*T), ttl time.Duration) *uuidMap[T] {

	return &uuidMap[T]{
		data:          make(map[string]*dataPair[T]),
		constructor:   constructor,
		destructor:    destructor,
		ttl:           ttl,
		latestCleanup: time.Now(),
	}
}

// Destroy releases every item.
func (ss *uuidMap[T]) Destroy() {
	for _, d := range ss.data {
		ss.destructor(d.value)
	}
}

// SetMaxSize limits the number of items. 0 or less is no limit, the default.
func (ss *uuidMap[T]) SetMaxSize(size int) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.maxSize = size
}

// Size return the number of item.
func (ss *uuidMap[T]) Size() int {
	ss.lock.RLock()
	defer ss.lock.RUnlock()

	return len(ss.data)
}

// cleanup drop the items that outlived the ttl and that nothing is using. It
// runs at most once per ttl, unless force asks for a sweep now: a map at its
// limit has to know whether it is really full before turning a caller away.
//
// Called with the lock held.
func (ss *uuidMap[T]) cleanup(force bool) {
	if !force && time.Since(ss.latestCleanup) < ss.ttl {
		return
	}

	ss.latestCleanup = time.Now()

	ids := []string{}
	for id, d := range ss.data {
		if time.Since(d.timestamp) > ss.ttl && !d.alive {
			ids = append(ids, id)
		}
	}
	for _, id := range ids {
		// An expired entry is as gone as a deleted one, so it gets the same
		// destructor: whatever it holds outside memory is nobody's otherwise.
		ss.destructor(ss.data[id].value)
		delete(ss.data, id)
	}
}

// New creates an alive item and returns its id. It returns
// [errUUIDMapFull] when the map is at its limit and nothing can be reclaimed.
func (ss *uuidMap[T]) New() (string, error) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	ss.cleanup(false)

	if ss.maxSize > 0 && len(ss.data) >= ss.maxSize {
		// A full map may only be full of items nobody came back for, so it
		// sweeps off-schedule rather than staying full until the next one.
		ss.cleanup(true)

		if len(ss.data) >= ss.maxSize {
			return "", tgutil.Errorf("%w", errUUIDMapFull)
		}
	}

	id := uuid.New().String()
	ss.data[id] = &dataPair[T]{
		value:     ss.constructor(),
		timestamp: time.Now(),
		alive:     true,
	}
	return id, nil
}

// SetAlive flag of id
func (ss *uuidMap[T]) SetAlive(id string, alive bool) {
	// The write lock, not the read one: this writes to the pair.
	ss.lock.Lock()
	defer ss.lock.Unlock()

	d, ok := ss.data[id]
	if !ok {
		return
	}

	d.timestamp = time.Now()
	d.alive = alive
}

// Get returns the item under id and whether it is alive; nil if none.
func (ss *uuidMap[T]) Get(id string) (*T, bool) {
	// The write lock, not the read one: this refreshes the timestamp.
	ss.lock.Lock()
	defer ss.lock.Unlock()

	d, ok := ss.data[id]
	if !ok {
		return nil, false
	}

	d.timestamp = time.Now()
	return d.value, d.alive
}

// Acquire takes the item under id and marks it alive, in one step. It
// returns [errUUIDNotFound] for no item and [errUUIDAlive] when it is held.
func (ss *uuidMap[T]) Acquire(id string) (*T, error) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	d, ok := ss.data[id]
	if !ok {
		return nil, errUUIDNotFound
	}

	// Checking and claiming under one lock is the point: between a Get that
	// reports an idle item and a SetAlive that claims it, a second caller
	// reads the same idle item and both end up sharing one value.
	if d.alive {
		return nil, errUUIDAlive
	}

	d.timestamp = time.Now()
	d.alive = true
	return d.value, nil
}

// Del delete uuid.
func (ss *uuidMap[T]) Del(id string) {
	ss.lock.Lock()
	defer ss.lock.Unlock()

	d, ok := ss.data[id]
	if !ok {
		return
	}

	ss.destructor(d.value)
	delete(ss.data, id)
}
