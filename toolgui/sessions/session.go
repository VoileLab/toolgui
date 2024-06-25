package sessions

import (
	"sync"

	"github.com/google/uuid"
)

// Sessions provide a goroutine-safe mapping from UUID to T
type Sessions[T any] interface {
	// New create a (id -> T) mapping and return id
	New() string

	// Get return T by id, return nil if id does not exist
	Get(id string) *T

	// Del delete uuid
	Del(id string)
}

type sessions[T any] struct {
	lock        sync.RWMutex
	data        map[string]*T
	constructor func() *T
}

// NewSessions create Sessions by providing the constructor of T
func NewSessions[T any](ctor func() *T) Sessions[T] {
	return &sessions[T]{
		data:        make(map[string]*T),
		constructor: ctor,
	}
}

// New create a (id -> T) mapping and return id
func (ss *sessions[T]) New() string {
	id := uuid.New().String()
	ss.lock.Lock()
	defer ss.lock.Unlock()
	ss.data[id] = ss.constructor()
	return id
}

// Get return T by id, return nil if id does not exist
func (ss *sessions[T]) Get(id string) *T {
	ss.lock.RLock()
	defer ss.lock.RUnlock()
	return ss.data[id]
}

// Del delete uuid
func (ss *sessions[T]) Del(id string) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	delete(ss.data, id)
}
