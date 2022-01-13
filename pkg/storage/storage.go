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

func (s *Storage) NewRetrieve() *retrieve {
	return newRetrieve(s)
}

func (s *Storage) NewCreate() *create {
	return newCreate(s)
}

func (s *Storage) retrieveEngine(r EngineRole) Engine {
	return s.cfg[r]
}

func (s *Storage) retrieveMDEngine() MetaDataEngine {
	e := s.retrieveEngine(EngineRoleMetaData)
	me, ok := e.(MetaDataEngine)
	if !ok {
		log.Fatalln("Couldn't bind engine")
	}
	return me
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
