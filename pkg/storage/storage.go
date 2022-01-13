package storage

import (
	"context"
	log "github.com/sirupsen/logrus"
)

type EngineConfig map[EngineRole]Engine

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
	err := s.retrieveMDEngine().Migrate(ctx, s.adapter(EngineRoleMetaData))
	return err
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

func (s *Storage) retrieveEngine(r EngineRole) Engine {
	return s.cfg[r]
}

func (s *Storage) retrieveMDEngine() MetaDataEngine {
	return s.retrieveEngine(EngineRoleMetaData).(MetaDataEngine)
}

func (s *Storage) adapter(role EngineRole) (a Adapter) {
	var err error
	switch role {
	case EngineRoleMetaData:
		a, err = s.pooler.Retrieve(s.retrieveMDEngine())
	}
	if err != nil {
		log.Fatalln(err)
	}
	return a
}
