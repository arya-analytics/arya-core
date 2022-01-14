package storage

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type EngineConfig map[EngineRole]BaseEngine

type Storage struct {
	cfg    EngineConfig
	pooler *Pooler
}

func New(cfg EngineConfig) *Storage {
	return &Storage{
		cfg:    cfg,
		pooler: NewPooler(),
	}
}

func (s *Storage) Migrate(ctx context.Context) error {
	return s.retrieveMDEngine().NewMigrate(s.adapter(EngineRoleMD)).Exec(ctx)
}

func (s *Storage) NewRetrieve() *retrieveQuery {
	return newRetrieve(s)
}

func (s *Storage) NewCreate() *createQuery {
	return newCreate(s)
}

func (s *Storage) NewDelete() *deleteQuery {
	return newDelete(s)
}

func (s *Storage) retrieveEngine(r EngineRole) BaseEngine {
	return s.cfg[r]
}

func (s *Storage) retrieveMDEngine() MDEngine {
	return s.retrieveEngine(EngineRoleMD).(MDEngine)
}

func (s *Storage) adapter(role EngineRole) (a Adapter) {
	var err error
	switch role {
	case EngineRoleMD:
		a, err = s.pooler.Retrieve(s.retrieveMDEngine())
	}
	if err != nil {
		log.Fatalln(err)
	}
	return a
}
