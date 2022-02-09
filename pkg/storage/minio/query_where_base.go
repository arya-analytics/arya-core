package minio

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/validate"
)

type queryWhereBase struct {
	queryBase
	pkChain model.PKChain
}

func (q *queryWhereBase) whereBasePK(pk interface{}) {
	q.pkChain = append(q.pkChain, model.NewPK(pk))
}

func (q *queryWhereBase) whereBasePKs(pks interface{}) {
	q.pkChain = model.NewPKChain(pks)
}

func (q *queryWhereBase) whereBaseValidateReq() {
	q.baseValidateReq()
	q.baseExec(func() error { return whereBaseQueryReqValidator.Exec(q) })
}

var whereBaseQueryReqValidator = validate.New([]validate.Func{
	validatePKProvided,
})

func validatePKProvided(v interface{}) error {
	q := v.(*queryWhereBase)
	if (len(q.pkChain)) == 0 {
		return storage.Error{Type: storage.ErrorTypeInvalidArgs,
			Message: "no PK provided to where query"}
	}
	return nil
}
