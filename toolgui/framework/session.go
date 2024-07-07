package framework

import (
	"encoding/json"
	"sync"

	"github.com/mudream4869/toolgui/toolgui/tgutil"
)

type Session struct {
	values map[string]any

	rwLock sync.RWMutex
}

func NewSession() *Session {
	return &Session{
		values: make(map[string]any),
	}
}

// Destroy release the resource hold by Session
func (s *Session) Destroy() {
}

// Copy do a swallow copy on session.Value
func (s *Session) Copy() *Session {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	ret := NewSession()
	for k, v := range s.values {
		ret.values[k] = v
	}
	return ret
}

func (s *Session) Set(key string, v any) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	s.values[key] = v
}

func (s *Session) GetObject(key string, out any) error {
	s.rwLock.RLock()
	val, ok := s.values[key]
	s.rwLock.RUnlock()

	if !ok {
		return nil
	}

	bs, err := json.Marshal(val)
	if err != nil {
		return tgutil.Errorf("%w", err)
	}

	err = json.Unmarshal(bs, out)
	if err != nil {
		return tgutil.Errorf("%w", err)
	}

	return nil
}

func (s *Session) GetString(key string) string {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	val, ok := s.values[key]
	if !ok {
		return ""
	}
	return val.(string)
}

func (s *Session) GetInt(key string) int {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	val, ok := s.values[key]
	if !ok {
		return 0
	}
	return val.(int)
}

func (s *Session) GetBool(key string) bool {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	val, ok := s.values[key]
	if !ok {
		return false
	}
	return val.(bool)
}
