package storage

type EngineConfig map[EngineRole]Engine

type Storage struct {
	cfg    EngineConfig
	pooler *Pooler
}

func NewStorage(cfg EngineConfig) *Storage {
	return &Storage{
		cfg:    cfg,
		pooler: NewPooler(),
	}
}
