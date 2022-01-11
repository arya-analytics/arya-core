package storage

type Storage struct {
	pooler   Pooler
	cfgChain ConfigChain
}

func NewStorage(cfgChain ConfigChain, pooler Pooler) *Storage {
	return &Storage{
		cfgChain: cfgChain,
		pooler:   pooler,
	}
}



