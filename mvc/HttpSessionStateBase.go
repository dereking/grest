package mvc

import (
	"github.com/dereking/grest/security"
)

type HttpSessionStateBase struct {
	SessionID string
	Timeout   int //second

	data map[string]interface{}
}

func NewHttpSession() *HttpSessionStateBase {
	ret := &HttpSessionStateBase{
		SessionID: security.GenSessionID(),
		Timeout:   60 * 30,
		data:      make(map[string]interface{}),
	}
	return ret
}

func (s *HttpSessionStateBase) Abandon() {

}
func (s *HttpSessionStateBase) Clear() {
	s.data = make(map[string]interface{})
}
func (s *HttpSessionStateBase) Remove(name string) {
	delete(s.data, name)
}
func (s *HttpSessionStateBase) RemoveAll() {
	s.data = make(map[string]interface{})
}
func (s *HttpSessionStateBase) Add(name string, value interface{}) {
	s.data[name] = value
}
