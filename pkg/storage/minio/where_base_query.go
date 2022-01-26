package minio

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
	"reflect"
)

type whereBaseQuery struct {
	baseQuery
	PKs []model.PK
}

func (r *whereBaseQuery) whereBasePK(pk interface{}) {
	r.PKs = append(r.PKs, model.NewPK(pk))
}

func (r *whereBaseQuery) whereBasePKs(pks interface{}) {
	rfl := reflect.ValueOf(pks)
	for i := 0; i < rfl.Len(); i++ {
		r.whereBasePK(rfl.Index(i).Interface())
	}
}

func (r *whereBaseQuery) whereBaseValidateReq() {
	r.baseValidateReq()
	r.catcher.Exec(func() error { return whereBaseQueryReqValidator.Exec(r) })
}

var whereBaseQueryReqValidator = validate.New([]validate.Func{
	validatePKProvided,
})

func validatePKProvided(v interface{}) error {
	q := v.(*whereBaseQuery)
	if (len(q.PKs)) == 0 {
		return storage.Error{Type: storage.ErrTypeInvalidArgs,
			Message: fmt.Sprintf("no PK provided to where query")}
	}
	return nil
}
