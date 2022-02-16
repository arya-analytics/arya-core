package cluster

import "github.com/arya-analytics/aryacore/pkg/util/model"

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

type queryOptFunc func(q *QueryRequest, args ...interface{})

// || PK ||

const wherePKKey = "PK"

func PKQueryOpt(q *QueryRequest) (model.PKChain, bool) {
	qo, ok := q.opts[wherePKKey]
	if !ok {
		return model.PKChain{}, false
	}
	return qo.(pkQueryOpt).PKChain, true
}

type pkQueryOpt struct {
	PKChain model.PKChain
}

func newPkQueryOpt(q *QueryRequest, args ...interface{}) {
	qo := pkQueryOpt{model.NewPKChain(args)}
	q.opts[wherePKKey] = qo
}

// || FIELDS ||

const whereFieldsKey = "Field"

type Fields map[string]interface{}

func FieldsQueryOpt(q *QueryRequest) Fields {
	qo, ok := q.opts[whereFieldsKey]
	if !ok {
		return Fields{}
	}
	return qo.(Fields)
}

func newFieldQueryOpts(q *QueryRequest, ops Fields) {
	q.opts[whereFieldsKey] = ops
}
