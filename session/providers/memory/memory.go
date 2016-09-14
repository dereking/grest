package memory

import (
	"container/list"
	"sync"
	"time"

	"github.com/dereking/grest/debug"
	"github.com/dereking/grest/session"
)

type Provider struct {
	lock     sync.Mutex               // lock
	sessions map[string]*list.Element // save in memory
	list     *list.List               // gc
}

var pder = &Provider{list: list.New()}

func Initialize() {

	pder.sessions = make(map[string]*list.Element, 0)
	session.Register("memory", pder)

	debug.Debug("memory session storage init ok")
}

func (p *Provider) SessionInit(sid string) (*session.SessionBase, error) {
	p.lock.Lock()
	defer p.lock.Unlock()
	//v := make(map[interface{}]interface{}, 0)
	//newsess := &SessionStore{sid: sid, timeAccessed: time.Now(), value: v}
	newsess := session.NewSession(sid)
	element := p.list.PushBack(newsess)
	p.sessions[sid] = element
	return newsess, nil
}

func (p *Provider) SessionRead(sid string) (*session.SessionBase, error) {
	if element, ok := p.sessions[sid]; ok {
		return element.Value.(*session.SessionBase), nil
	} else {
		sess, err := p.SessionInit(sid)
		return sess, err
	}
	return nil, nil
}

func (p *Provider) SessionDestroy(sid string) error {
	if element, ok := p.sessions[sid]; ok {
		delete(p.sessions, sid)
		p.list.Remove(element)
		return nil
	}
	return nil
}

func (p *Provider) SessionGC(maxlifetime int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for {
		element := p.list.Back()
		if element == nil {
			break
		}
		if (element.Value.(*session.SessionBase).TimeAccessed.Unix() + maxlifetime) < time.Now().Unix() {
			p.list.Remove(element)
			delete(p.sessions, element.Value.(*session.SessionBase).SessionID())
		} else {
			break
		}
	}
}

func (p *Provider) SessionUpdate(sid string) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if element, ok := p.sessions[sid]; ok {
		element.Value.(*session.SessionBase).TimeAccessed = time.Now()
		p.list.MoveToFront(element)
		return nil
	}
	return nil
}
