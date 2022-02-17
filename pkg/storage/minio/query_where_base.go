package minio

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
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
