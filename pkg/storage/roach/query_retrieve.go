package roach

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/uptrace/bun"
)

type queryRetrieve struct {
	queryBase
	bunQ     *bun.SelectQuery
	scanArgs []interface{}
}

func newRetrieve(db *bun.DB) *queryRetrieve {
	q := &queryRetrieve{bunQ: db.NewSelect()}
	q.baseInit(db)
	return q
}

func (q *queryRetrieve) Model(m interface{}) storage.QueryMDRetrieve {
	q.baseModel(m)
	q.bunQ = q.bunQ.Model(q.baseDest().Pointer())
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

func (q *queryRetrieve) WhereFields(flds query.WhereFields) storage.QueryMDRetrieve {
	for fldN, fldV := range flds {
		relN, _ := model.SplitLastFieldName(fldN)
		if relN != "" {
			q.bunQ = q.bunQ.Relation(relN)
		}
		fldExp, args := q.baseSQL().relFldExp(fldN, fldV)
		q.bunQ = q.bunQ.Where(fldExp, args...)
	}
	return q
}

func (q *queryRetrieve) Relation(rel string, fields ...string) storage.QueryMDRetrieve {
	q.bunQ = q.bunQ.Relation(rel, func(sq *bun.SelectQuery) *bun.SelectQuery {
		return sq.Column(q.baseSQL().fieldNames(fields...)...)
	})
	return q
}

func (q *queryRetrieve) Fields(flds ...string) storage.QueryMDRetrieve {
	q.bunQ = q.bunQ.Column(q.baseSQL().fieldNames(flds...)...)
	return q
}

func (q *queryRetrieve) Calculate(c storage.Calculate, fldName string, into interface{}) storage.QueryMDRetrieve {
	q.addScanArg(into)
	q.bunQ = q.bunQ.ColumnExpr(q.baseSQL().calculate(c), bun.Ident(q.baseSQL().fieldName(fldName)))
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
		err := q.bunQ.Scan(ctx, q.scanArgs...)
		return err
	})
	q.baseExchangeToSource()
	return q.baseErr()
}

func (q *queryRetrieve) addScanArg(arg interface{}) {
	q.scanArgs = append(q.scanArgs, arg)
}
