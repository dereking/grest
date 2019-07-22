package session

import (
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/dereking/grest/log"
	"go.uber.org/zap"

	"github.com/dereking/grest/security"
)

type SessionManager struct {
	//all         map[string]ISession
	cookieName  string     //private cookiename
	lock        sync.Mutex // protects session
	provider    IProvider
	maxlifetime int64
}

var provides = make(map[string]IProvider)
var manager *SessionManager

func newManager(provideName, cookieName string, maxlifetime int64) (*SessionManager, error) {
	provider, ok := provides[provideName]
	if !ok {
		return nil, fmt.Errorf("session: unknown provide %q (forgotten import?)", provideName)
	}
	ret := &SessionManager{
		provider:    provider,
		cookieName:  cookieName,
		maxlifetime: maxlifetime,
		//all: make(map[string]ISessionStore),
	}

	ret.startGC()
	return ret, nil
}

// Register makes a session provider available by the provided name.
// If a Register is called twice with the same name or if the driver is nil,
// it panics.
func Register(name string, provider IProvider) {
	if provider == nil {
		panic("session: Register provide is nil")
	}
	if _, dup := provides[name]; dup {
		panic("session: Register called twice for provide " + name)
	}
	provides[name] = provider
}

//get default memory sessionmanager.
func GetSessionManager() *SessionManager {

	if manager == nil {
		var err error
		manager, err = newManager("memory", "gosessionid", 3600)
		if err != nil {
			log.Logger().Error("SessionManager init", zap.Error(err))
		}
	}

	return manager
}

func (manager *SessionManager) SessionStart(w http.ResponseWriter, r *http.Request) (session *SessionBase) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	cookie, err := r.Cookie(manager.cookieName)
	if err != nil || cookie.Value == "" {
		sid := security.GenSessionID()
		session, _ = manager.provider.SessionInit(sid)
		cookie := http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
		http.SetCookie(w, &cookie)

	} else {
		sid, _ := url.QueryUnescape(cookie.Value)
		session, _ = manager.provider.SessionRead(sid)
	}
	return
}

//GC timeout sessions.
func (manager *SessionManager) startGC() {
	go func() {
		timer := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-timer.C:
				//fmt.Println("Timer has expired.")
				manager.provider.SessionGC(manager.maxlifetime)
			}
		}
		timer.Stop()
	}()
}
