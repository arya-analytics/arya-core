package internal

import (
	"context"
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

//go:generate stringer -type=QueryVariant
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

func PKQueryOpt(qr *QueryRequest) (model.PKChain, bool) {
	qo, ok := qr.opts[pkQueryOptKey]
	if !ok {
		return model.PKChain{}, false
	}
	return qo.(pkQueryOpt).PKChain, true
}

type pkQueryOpt struct {
	PKChain model.PKChain
}

func NewPKQueryOpt(qr *QueryRequest, args ...interface{}) {
	panicWhenAlreadySpecified(qr, pkQueryOptKey)

	// Handling a single vs multi PK query
	pks := args[0]
	if reflect.TypeOf(args[0]).Kind() != reflect.Slice {
		pks = args
	}

	qo := pkQueryOpt{model.NewPKChain(pks)}
	qr.opts[pkQueryOptKey] = qo
}

// || FIELDS ||

const fieldQueryOptKey = "FieldQueryOpt"

type Fields map[string]interface{}

func (f Fields) Retrieve(fldName string) (interface{}, bool) {
	fld, ok := f[fldName]
	return fld, ok
}

func FieldsQueryOpt(qr *QueryRequest) (Fields, bool) {
	qo, ok := qr.opts[fieldQueryOptKey]
	if !ok {
		return Fields{}, false
	}
	return qo.(Fields), true
}

func NewFieldsQueryOpt(q *QueryRequest, ops Fields) {
	panicWhenAlreadySpecified(q, fieldQueryOptKey)
	q.opts[fieldQueryOptKey] = ops
}

type QueryRequestVariantOperations map[QueryVariant]ServiceOperation

func SwitchQueryRequestVariant(ctx context.Context, qr *QueryRequest, qrvo QueryRequestVariantOperations) error {
	op, ok := qrvo[qr.Variant]
	if !ok {
		panic(fmt.Sprintf("%s not supported for model %s", qr.Variant, qr.Model.Type().Name()))
	}
	return op(ctx, qr)
}
