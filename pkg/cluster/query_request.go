package cluster

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

// |||| REQUEST ||||

type QueryRequest struct {
	Variant QueryVariant
	Model   *model.Reflect
	opts    map[string]interface{}
}

// |||| VARIANT ||||

type QueryVariant int

const (
	QueryVariantCreate QueryVariant = iota
	QueryVariantRetrieve
	QueryVariantUpdate
	QueryVariantDelete
)

// |||| QUERY OPTS ||||

// || UTILITIES ||

func panicWhenAlreadySpecified(q *QueryRequest, optKey string) {
	_, ok := q.opts[optKey]
	if ok {
		panic(fmt.Sprintf("%s already specified. There must be a duplicate method call in your query!", optKey))
	}
}

// || PK ||

const pkQueryOptKey = "PK Query Opt"

func PKQueryOpt(q *QueryRequest) (model.PKChain, bool) {
	qo, ok := q.opts[pkQueryOptKey]
	if !ok {
		return model.PKChain{}, false
	}
	return qo.(pkQueryOpt).PKChain, true
}

type pkQueryOpt struct {
	PKChain model.PKChain
}

func newPkQueryOpt(q *QueryRequest, args ...interface{}) {
	panicWhenAlreadySpecified(q, pkQueryOptKey)
	qo := pkQueryOpt{model.NewPKChain(args)}
	q.opts[pkQueryOptKey] = qo
}

// || FIELDS ||

const fieldQueryOptKey = "Field Query Opt"

type Fields map[string]interface{}

func FieldsQueryOpt(q *QueryRequest) Fields {
	qo, ok := q.opts[fieldQueryOptKey]
	if !ok {
		return Fields{}
	}
	return qo.(Fields)
}

func newFieldQueryOpts(q *QueryRequest, ops Fields) {
	panicWhenAlreadySpecified(q, fieldQueryOptKey)
	q.opts[fieldQueryOptKey] = ops
}
