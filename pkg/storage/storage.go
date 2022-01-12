package storage

import log "github.com/sirupsen/logrus"

type EngineConfig map[EngineRole]Engine

type Storage struct {
	cfg    EngineConfig
	pooler *Pooler
}

func (s *Storage) NewRetrieve() {

}

func (s *Storage) NewCreate() {

}

func (s *Storage) NewDelete() {

}

func (s *Storage) NewUpdate() {

}

func (s *Storage) retrieveEngine(r EngineRole) Engine {
	return s.cfg[r]
}

func (s *Storage) bindMetaData(e interface{}) MetaDataEngine {
	me, ok := e.(MetaDataEngine)
	if !ok {
		log.Fatalln("Couldn't bind engine")
	}
	return me
}