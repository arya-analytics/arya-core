package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
)

type whereBaseQuery struct {
	baseQuery
	pkChain model.PKChain
}

func (w *whereBaseQuery) whereBasePK(pk interface{}) {
	w.pkChain = append(w.pkChain, model.NewPK(pk))
}

func (w *whereBaseQuery) whereBasePKs(pks interface{}) {
	w.pkChain = model.NewPKChain(pks)
}

func (w *whereBaseQuery) whereBaseValidateReq() {
	w.baseValidateReq()
	w.baseExec(func() error { return whereBaseQueryReqValidator.Exec(w) })
}

var whereBaseQueryReqValidator = validate.New([]validate.Func{
	validatePKProvided,
})

func validatePKProvided(v interface{}) error {
	q := v.(*whereBaseQuery)
	if (len(q.pkChain)) == 0 {
		return storage.Error{Type: storage.ErrTypeInvalidArgs,
			Message: "no PK provided to where query"}
	}
	return nil
}
