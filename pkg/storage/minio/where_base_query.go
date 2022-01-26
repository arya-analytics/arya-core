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

func (w *whereBaseQuery) whereBasePK(pk interface{}) {
	w.PKs = append(w.PKs, model.NewPK(pk))
}

func (w *whereBaseQuery) whereBasePKs(pks interface{}) {
	rfl := reflect.ValueOf(pks)
	for i := 0; i < rfl.Len(); i++ {
		w.whereBasePK(rfl.Index(i).Interface())
	}
}

func (w *whereBaseQuery) whereBaseValidateReq() {
	w.baseValidateReq()
	w.catcher.Exec(func() error { return whereBaseQueryReqValidator.Exec(w) })
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
