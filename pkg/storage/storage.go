package storage

type Storage struct {
	engines []BaseEngine
	pooler Pooler
}
