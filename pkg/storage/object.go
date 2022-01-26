package storage

type Object interface {
	Read(b []byte) (n int, err error)
	Size() int64
}

type ObjectInfo struct {
	ID   string
	Size int64
}
