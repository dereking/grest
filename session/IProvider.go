package session

type IProvider interface {
	SessionInit(sid string) (*SessionBase, error)
	SessionRead(sid string) (*SessionBase, error)
	SessionDestroy(sid string) error
	SessionGC(maxLifeTime int64)
}
