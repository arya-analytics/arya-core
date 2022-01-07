package storage

type Pooler interface {
	Retrieve(key string) Adapter
}

type Adapter interface {
	Release() error
	Status() ConnStatus
	Conn() interface{}
	close() error
	open() error
}

type ConnStatus int

const (
	Ready ConnStatus = iota
)


