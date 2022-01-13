package storage

type base struct {
	storage  *Storage
	mdEngine MetaDataEngine
	_mdQuery MetaDataQuery
}

func (b *base) init(s *Storage) {
	b.storage = s
	b.mdEngine = b.storage.retrieveMDEngine()
}









































