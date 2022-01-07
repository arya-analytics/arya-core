package storage

type Storage struct {
	engines []Engine
	pooler Pooler
}
