package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/uptrace/bun"
)

type queryRetrieve struct {
	queryBase
	bunQ *bun.SelectQuery
}

func newRetrieve(db *bun.DB) *queryRetrieve {
	q := &queryRetrieve{bunQ: db.NewSelect()}
	q.baseInit(db)
	return q
}

func (q *queryRetrieve) Model(m interface{}) storage.QueryMDRetrieve {
	q.baseModel(m)
	q.bunQ = q.bunQ.Model(q.Dest().Pointer())
	return q
}

func (q *queryRetrieve) where(query string, args ...interface{}) storage.QueryMDRetrieve {
	q.bunQ = q.bunQ.Where(query, args...)
	return q
}

func (q *queryRetrieve) WherePK(pk interface{}) storage.QueryMDRetrieve {
	return q.where(q.baseSQL().pk(), pk)
}

func (q *queryRetrieve) WherePKs(pks interface{}) storage.QueryMDRetrieve {
	return q.where(q.baseSQL().pks(), bun.In(pks))
}

func (q *queryRetrieve) WhereFields(flds storage.Fields) storage.QueryMDRetrieve {
	for fldName, fldValue := range flds {
		relName, baseFldName := model.SplitLastFieldName(fldName)
		baseFldSQL := q.baseSQL().fieldEquals(baseFldName)
		if relName != "" {
			q.bunQ = q.bunQ.Relation(relName, func(sq *bun.SelectQuery) *bun.SelectQuery {
				return sq.Where(baseFldSQL, fldValue)
			})
		} else {
			q.bunQ = q.bunQ.Where(baseFldSQL, fldValue)
		}
	}
	return q
}

func (q *queryRetrieve) Relation(rel string, fields ...string) storage.QueryMDRetrieve {
	q.bunQ = q.bunQ.Relation(rel, func(sq *bun.SelectQuery) *bun.SelectQuery {
		return sq.Column(
			q.baseSQL().fieldNamesToSQL(fields...)...,
		)
	})
	return q
}

func (q *queryRetrieve) Fields(flds ...string) storage.QueryMDRetrieve {
	q.bunQ = q.bunQ.Column(flds...)
	return q
}

func (q *queryRetrieve) Count(ctx context.Context) (count int, err error) {
	q.baseExec(func() error {
		count, err = q.bunQ.Count(ctx)
		return err
	})
	return count, q.baseErr()
}

func (q *queryRetrieve) Exec(ctx context.Context) error {
	q.baseExec(func() error {
		err := q.bunQ.Scan(ctx)
		return err
	})
	q.baseExchangeToSource()
	return q.baseErr()
}
