package internal

import (
	"context"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/model"
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

func NewPKQueryOpt(qr *QueryRequest, pks interface{}) {
	panicWhenAlreadySpecified(qr, pkQueryOptKey)
	qo := pkQueryOpt{model.NewPKChain(pks)}
	qr.opts[pkQueryOptKey] = qo
}

// || FIELDS ||

const fieldQueryOptKey = "FieldQueryOpt"

func FieldsQueryOpt(qr *QueryRequest) (model.WhereFields, bool) {
	qo, ok := qr.opts[fieldQueryOptKey]
	if !ok {
		return model.WhereFields{}, false
	}
	return qo.(model.WhereFields), true
}

func NewFieldsQueryOpt(q *QueryRequest, ops model.WhereFields) {
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

// || RELATION ||

type RelationQueryOpt struct {
	Rel    string
	Fields []string
}

const relationQueryOptKey = "RelationQueryOpt"

func NewRelationQueryOpt(qr *QueryRequest, rel string, fields ...string) {
	rq := RelationQueryOpt{rel, fields}
	_, ok := qr.opts[relationQueryOptKey]
	if !ok {
		qr.opts[relationQueryOptKey] = []RelationQueryOpt{rq}
	} else {
		qr.opts[relationQueryOptKey] = append(qr.opts[relationQueryOptKey].([]RelationQueryOpt), rq)
	}
}

func RelationQueryOpts(qr *QueryRequest) []RelationQueryOpt {
	opts, ok := qr.opts[relationQueryOptKey]
	if !ok {
		return []RelationQueryOpt{}
	}
	return opts.([]RelationQueryOpt)
}
