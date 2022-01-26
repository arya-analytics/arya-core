package storage

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
)

type retrieveQuery struct {
	baseQuery
	modelRfl *model.Reflect
	PKs      []model.PK
}

// |||| CONSTRUCTOR ||||

func newRetrieve(s *Storage) *retrieveQuery {
	r := &retrieveQuery{}
	r.baseInit(s)
	return r
}

// |||| INTERFACE ||||

func (r *retrieveQuery) WherePK(pk interface{}) *retrieveQuery {
	r.PKs = append(r.PKs, model.NewPK(pk))
	r.setMDQuery(r.mdQuery().WherePK(pk))
	return r
}

func (r *retrieveQuery) WherePKs(pks interface{}) *retrieveQuery {
	r.PKs = model.NewMultiPK(pks)
	r.setMDQuery(r.mdQuery().WherePKs(pks))
	return r
}

func (r *retrieveQuery) Model(m interface{}) *retrieveQuery {
	r.modelRfl = model.NewReflect(m)
	r.setMDQuery(r.mdQuery().Model(r.modelRfl.Pointer()))
	return r
}

func (r *retrieveQuery) Exec(ctx context.Context) error {
	r.validateReq()
	r.catcher.Exec(func() error {
		err := r.mdQuery().Exec(ctx)
		return err
	})
	if r.objEngine.InCatalog(r.modelRfl.Pointer()) {
		r.catcher.Exec(func() error {
			return r.objQuery().Model(r.modelRfl.Pointer()).WherePKs(r.modelRfl.PKs()).
				Exec(ctx)
		})
	}
	return r.baseErr()
}

// |||| QUERY BINDING ||||

// || META DATA ||

func (r *retrieveQuery) mdQuery() MDRetrieveQuery {
	if r.baseMDQuery() == nil {
		r.setMDQuery(r.mdEngine.NewRetrieve(r.baseMDAdapter()))
	}
	return r.baseMDQuery().(MDRetrieveQuery)
}

func (r *retrieveQuery) setMDQuery(q MDRetrieveQuery) {
	r.baseSetMDQuery(q)
}

// || OBJECT ||

func (r *retrieveQuery) objQuery() ObjectRetrieveQuery {
	if r.baseObjQuery() == nil {
		r.setObjQuery(r.objEngine.NewRetrieve(r.baseObjAdapter()))
	}
	return r.baseObjQuery().(ObjectRetrieveQuery)
}

func (r *retrieveQuery) setObjQuery(q ObjectRetrieveQuery) {
	r.baseSetObjQuery(q)
}

// || VALIDATORS ||

func (r *retrieveQuery) validateReq() {
	r.catcher.Exec(func() error { return retrieveQueryReqValidator.Exec(r) })

}

// ||||| VALIDATORS ||||

var retrieveQueryReqValidator = validate.New([]validate.Func{
	validatePKModelMatch,
})

func validatePKModelMatch(v interface{}) error {
	q := v.(*retrieveQuery)
	if q.modelRfl.IsStruct() && len(q.PKs) > 1 {
		return Error{Type: ErrTypeInvalidArgs,
			Message: "a struct model was provided when querying multiple pks"}
	}
	return nil
}
