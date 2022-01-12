package storage

type EngineConfig map[EngineRole]EngineBase

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
