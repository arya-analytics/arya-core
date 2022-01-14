package storage

type baseQuery struct {
	storage  *Storage
	mdEngine MDEngine
}

func (b *baseQuery) init(s *Storage) {
	b.storage = s
	b.mdEngine = b.storage.cfg.mdEngine()
}