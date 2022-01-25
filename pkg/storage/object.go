package storage

type Object interface {
	Read(b []byte) (n int, err error)
	Size() int64
	//Stat() (ObjectInfo, errutil)
	//ReadAt(b []byte, offset int64) (n int, err errutil)
	//Seek(offset int64, whence int) (n int64, err errutil)
	//Close() (err errutil)
}

type ObjectInfo struct {
	ID   string
	Size int64
}
