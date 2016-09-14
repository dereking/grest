package session

import (
	"time"
)

//"github.com/dereking/grest/security"

type SessionBase struct {
	data         map[interface{}]interface{}
	sid          string
	TimeAccessed time.Time // last access time
}

func NewSession(sid string) *SessionBase {
	return &SessionBase{
		sid:          sid,
		TimeAccessed: time.Now(),
		data:         make(map[interface{}]interface{}),
	}
}

func (s *SessionBase) SessionID() string {
	return s.sid
}

func (s *SessionBase) Clear() {
	s.data = make(map[interface{}]interface{})
}
func (s *SessionBase) Remove(key interface{}) error {
	delete(s.data, key)
	return nil
}
func (s *SessionBase) Set(key interface{}, value interface{}) error {
	s.data[key] = value
	return nil
}

func (s *SessionBase) Get(key interface{}) interface{} {
	return s.data[key]
}
func (s *SessionBase) GetString(key interface{}) string {
	r := s.data[key]

	if r != nil {
		return r.(string)
	} else {

		return ""
	}
}

func (s *SessionBase) GetInt64(key interface{}) int64 {
	r := s.data[key]

	if r != nil {
		return r.(int64)
	} else {
		return 0
	}
}

func (s *SessionBase) GetFloat64(key interface{}) float64 {
	r := s.data[key]

	if r != nil {
		return r.(float64)
	} else {
		return 0.0
	}
}

func (s *SessionBase) GetBool(key interface{}) bool {
	r := s.data[key]

	if r != nil {
		return r.(bool)
	} else {
		return false
	}
}

func (s *SessionBase) GetTime(key interface{}) time.Time {
	r := s.data[key]

	if r != nil {
		return r.(time.Time)
	} else {
		return time.Date(1970, time.February, 0, 0, 0, 0, 0, time.Local)
	}
}
