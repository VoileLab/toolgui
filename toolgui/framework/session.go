package framework

import "sync"

type Session struct {
	values map[string]any

	rwLock sync.RWMutex
}

func NewSession() *Session {
	return &Session{
		values: make(map[string]any),
	}
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
