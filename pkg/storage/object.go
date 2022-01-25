package storage

type Object interface {
	Read(b []byte) (n int, err error)
	Size() int64
	//Stat() (ObjectInfo, error)
	//ReadAt(b []byte, offset int64) (n int, err error)
	//Seek(offset int64, whence int) (n int64, err error)
	//Close() (err error)
}

type ObjectInfo struct {
	ID   string
	Size int64
}
