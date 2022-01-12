package storage

type Retrieve struct {
	s        *Storage
	mdEngine MetaDataEngine
	md       MetaDataRetrieve
}

func NewRetrieve(s *Storage) *Retrieve {
	mdEngine := s.bindMetaData(s.retrieveEngine(EngineRoleMetaData))
	return &Retrieve{
		s:        s,
		mdEngine: mdEngine,
	}
}

func (r *Retrieve) retrieveMD() MetaDataRetrieve {
	if r.md == nil {
		a, _ := r.s.pooler.Retrieve(r.mdEngine)
		r.md = r.mdEngine.NewRetrieve(a)
	}
	return r.md

}

func (r *Retrieve) Model(model interface{}) {
	r.md = r.retrieveMD().Model(model)
}

