package minio

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

type whereBaseQuery struct {
	baseQuery
	pks []model.PK
}

func (r *whereBaseQuery) whereBasePK(pk interface{}) {
	r.pks = append(r.pks, model.NewPK(pk))
}

func (r *whereBaseQuery) whereBasePKs(pks interface{}) {
	rfl := reflect.ValueOf(pks)
	for i := 0; i < rfl.Len(); i++ {
		r.whereBasePK(rfl.Index(i).Interface())
	}
}
