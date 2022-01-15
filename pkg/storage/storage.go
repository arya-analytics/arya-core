package storage

import (
	log "github.com/sirupsen/logrus"
)

// |||| ENGINE CONFIG ||||

type EngineConfig map[EngineRole]BaseEngine

func (ec EngineConfig) retrieve(r EngineRole) BaseEngine {
	return ec[r]
}

func (ec EngineConfig) mdEngine() MDEngine {
	md, ok := ec.retrieve(EngineRoleMD).(MDEngine)
	if !ok {
		log.Fatalln("Could not bind meta data engine.")
	}
	return md
}

// |||| STORAGE ||||

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

func (s *Storage) NewMigrate() *migrateQuery {
	return newMigrate(s)
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

func (s *Storage) adapter(r EngineRole) (a Adapter) {
	var err error
	a, err = s.pooler.Retrieve(s.cfg.retrieve(r))
	if err != nil {
		log.Fatalln(err)
	}
	return a
}
