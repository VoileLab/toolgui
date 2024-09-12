package tgframe

import (
	"encoding/json"
	"maps"
	"runtime"
	"sync"

	"github.com/mudream4869/toolgui/toolgui/tgutil"
)

type State struct {
	values    map[string]any
	files     map[string][]byte
	funcCache map[string]map[string]any

	clickID string

	rwLock sync.RWMutex
}

func NewState() *State {
	return &State{
		values:    make(map[string]any),
		files:     make(map[string][]byte),
		funcCache: make(map[string]map[string]any),
	}
}

// Destroy release the resource.
func (s *State) Destroy() {
}

// Clone do a swallow copy on [State].
func (s *State) Clone() *State {
	s.rwLock.RLock()
	defer s.rwLock.RUnlock()
	return &State{
		values:    maps.Clone(s.values),
		files:     maps.Clone(s.files),
		funcCache: maps.Clone(s.funcCache),
		clickID:   s.clickID,
	}
}

// SetClickID set the id of clicked button.
func (s *State) SetClickID(id string) {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()
	s.clickID = id
}

// GetClickID get the id of clicked button.
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

func (s *State) Default(key string, v any) any {
	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	_, ok := s.values[key]
	if !ok {
		s.values[key] = v
	}

	return s.values[key]
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

func (s *State) SetFuncCache(key string, value any, funcNameParam ...string) {
	funcName := tgutil.ParamsToParam(funcNameParam)
	if funcName == "" {
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
	}

	s.rwLock.Lock()
	defer s.rwLock.Unlock()

	_, ok := s.funcCache[funcName]
	if !ok {
		s.funcCache[funcName] = make(map[string]any)
	}

	s.funcCache[funcName][key] = value
}

func (s *State) GetFuncCache(key string, funcNameParams ...string) any {
	funcName := tgutil.ParamsToParam(funcNameParams)
	if funcName == "" {
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			funcName = runtime.FuncForPC(pc).Name()
		}
	}

	s.rwLock.RLock()
	defer s.rwLock.RUnlock()

	cache, ok := s.funcCache[funcName]
	if !ok {
		return nil
	}

	return cache[key]
}
