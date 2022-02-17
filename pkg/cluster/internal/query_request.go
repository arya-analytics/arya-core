package internal

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

// |||| REQUEST ||||

type QueryRequest struct {
	Variant QueryVariant
	Model   *model.Reflect
	opts    map[string]interface{}
}

func NewQueryRequest(variant QueryVariant, model *model.Reflect) *QueryRequest {
	return &QueryRequest{
		Variant: variant,
		Model:   model,
		opts:    map[string]interface{}{},
	}
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

const pkQueryOptKey = "PKQueryOpt"

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

func NewPKQueryOpt(q *QueryRequest, args ...interface{}) {
	panicWhenAlreadySpecified(q, pkQueryOptKey)

	// Handling a single vs multi PK query
	pks := args[0]
	if reflect.TypeOf(args[0]).Kind() != reflect.Slice {
		pks = args
	}

	qo := pkQueryOpt{model.NewPKChain(pks)}
	q.opts[pkQueryOptKey] = qo
}

// || FIELDS ||

const fieldQueryOptKey = "FieldQueryOpt"

type Fields map[string]interface{}

func FieldsQueryOpt(q *QueryRequest) Fields {
	qo, ok := q.opts[fieldQueryOptKey]
	if !ok {
		return Fields{}
	}
	return qo.(Fields)
}

func NewFieldsQueryOpt(q *QueryRequest, ops Fields) {
	panicWhenAlreadySpecified(q, fieldQueryOptKey)
	q.opts[fieldQueryOptKey] = ops
}
