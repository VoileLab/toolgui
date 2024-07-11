package framework

import (
	"encoding/json"
	"sync"

	"github.com/mudream4869/toolgui/toolgui/tgutil"
)

type State struct {
	values map[string]any
	files  map[string][]byte

	clickID string

	rwLock sync.RWMutex
}

func NewState() *State {
	return &State{
		values: make(map[string]any),
		files:  make(map[string][]byte),
	}
}

// Destroy release the resource hold by State
func (s *State) Destroy() {
}

// Copy do a swallow copy on State
func (s *State) Copy() *State {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	ret := NewState()
	for k, v := range s.values {
		ret.values[k] = v
	}
	return ret
}

func (s *State) SetClickID(id string) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.clickID = id
}

func (s *State) GetClickID() string {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	return s.clickID
}

func (s *State) Set(key string, v any) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	s.values[key] = v
}

func (s *State) GetObject(key string, out any) error {
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

func (s *State) GetString(key string) string {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	val, ok := s.values[key]
	if !ok {
		return ""
	}
	return val.(string)
}

func (s *State) GetInt(key string) int {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	val, ok := s.values[key]
	if !ok {
		return 0
	}
	return val.(int)
}

func (s *State) GetBool(key string) bool {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	val, ok := s.values[key]
	if !ok {
		return false
	}
	return val.(bool)
}

func (s *State) SetFile(key string, bs []byte) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	s.files[key] = bs
}

func (s *State) GetFile(key string) []byte {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	return s.files[key]
}
