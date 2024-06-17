package sessions

import (
	"sync"

	"github.com/google/uuid"
)

type Sessions[T any] interface {
	// New return id to session
	New() string

	// Get return session by id
	Get(id string) *T

	// Del delete session by id
	Del(id string)
}

type sessions[T any] struct {
	lock        sync.RWMutex
	data        map[string]*T
	constructor func() *T
}

func NewSessions[T any](ctor func() *T) Sessions[T] {
	return &sessions[T]{
		data:        make(map[string]*T),
		constructor: ctor,
	}
}

func (ss *sessions[T]) New() string {
	id := uuid.New().String()
	ss.lock.Lock()
	defer ss.lock.Unlock()
	ss.data[id] = ss.constructor()
	return id
}

func (ss *sessions[T]) Get(id string) *T {
	ss.lock.RLock()
	defer ss.lock.RUnlock()
	return ss.data[id]
}

func (ss *sessions[T]) Del(id string) {
	ss.lock.Lock()
	defer ss.lock.Unlock()
	delete(ss.data, id)
}
