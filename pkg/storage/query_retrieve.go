package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type RetrieveQuery struct {
	baseQuery
	modelRfl *model.Reflect
}

// |||| CONSTRUCTOR ||||

func newRetrieve(s *Storage) *RetrieveQuery {
	r := &RetrieveQuery{}
	r.baseInit(s)
	return r
}

// |||| INTERFACE ||||

func (r *RetrieveQuery) WherePK(pk interface{}) *RetrieveQuery {
	r.setMDQuery(r.mdQuery().WherePK(pk))
	return r
}

func (r *RetrieveQuery) WherePKs(pks interface{}) *RetrieveQuery {
	r.setMDQuery(r.mdQuery().WherePKs(pks))
	return r
}

func (r *RetrieveQuery) Model(m interface{}) *RetrieveQuery {
	r.modelRfl = model.NewReflect(m)
	r.setMDQuery(r.mdQuery().Model(r.modelRfl.Pointer()))
	return r
}

func (r *RetrieveQuery) Exec(ctx context.Context) error {
	r.catcher.Exec(func() error {
		err := r.mdQuery().Exec(ctx)
		return err
	})
	if r.storage.cfg.objEngine().InCatalog(r.modelRfl.Pointer()) {
		r.catcher.Exec(func() error {
			return r.objQuery().Model(r.modelRfl.Pointer()).WherePKs(r.modelRfl.PKChain().Raw()).
				Exec(ctx)
		})
	}
	return r.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (r *RetrieveQuery) mdQuery() MDRetrieveQuery {
	if r.baseMDQuery() == nil {
		r.setMDQuery(r.baseMDEngine().NewRetrieve(r.baseMDAdapter()))
	}
	return r.baseMDQuery().(MDRetrieveQuery)
}

func (r *RetrieveQuery) setMDQuery(q MDRetrieveQuery) {
	r.baseSetMDQuery(q)
}

// || OBJECT ||

func (r *RetrieveQuery) objQuery() ObjectRetrieveQuery {
	if r.baseObjQuery() == nil {
		r.setObjQuery(r.baseObjEngine().NewRetrieve(r.baseObjAdapter()))
	}
	return r.baseObjQuery().(ObjectRetrieveQuery)
}

func (r *RetrieveQuery) setObjQuery(q ObjectRetrieveQuery) {
	r.baseSetObjQuery(q)
}
