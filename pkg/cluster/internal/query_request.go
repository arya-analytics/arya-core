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

func NewPKQueryOpt(qr *QueryRequest, pk interface{}) {
	panicWhenAlreadySpecified(qr, pkQueryOptKey)
	qo := pkQueryOpt{model.NewPKChain([]interface{}{pk})}
	qr.opts[pkQueryOptKey] = qo
}

func NewPKsQueryOpt(qr *QueryRequest, pks interface{}) {
	panicWhenAlreadySpecified(qr, pkQueryOptKey)
	qo := pkQueryOpt{model.NewPKChain(pks)}
	qr.opts[pkQueryOptKey] = qo
}

// || WHERE FIELDS ||

const whereFieldsQueryOptKey = "WhereFieldsQueryOpt"

func WhereFieldsQueryOpt(qr *QueryRequest) (model.WhereFields, bool) {
	qo, ok := qr.opts[whereFieldsQueryOptKey]
	if !ok {
		return model.WhereFields{}, false
	}
	return qo.(model.WhereFields), true
}

func NewWhereFieldsQueryOpt(q *QueryRequest, ops model.WhereFields) {
	panicWhenAlreadySpecified(q, whereFieldsQueryOptKey)
	q.opts[whereFieldsQueryOptKey] = ops
}

// || FIELDS ||

const fieldsQueryOptkey = "RetrieveFieldsQueryOpt"

type FieldsQueryOpt []string

func (fqo FieldsQueryOpt) ContainsAny(flds ...string) (contains bool) {
	for _, qFld := range flds {
		for _, fld := range fqo {
			if qFld == fld {
				contains = true
			}
		}
	}
	return contains
}

func NewFieldsQueryOpt(qr *QueryRequest, flds ...string) {
	panicWhenAlreadySpecified(qr, fieldsQueryOptkey)
	qo := FieldsQueryOpt{}
	qo = append(qo, flds...)
	qr.opts[fieldsQueryOptkey] = qo
}

func RetrieveFieldsQueryOpt(qr *QueryRequest) (FieldsQueryOpt, bool) {
	qo, ok := qr.opts[fieldsQueryOptkey]
	if !ok {
		return FieldsQueryOpt{}, false
	}
	return qo.(FieldsQueryOpt), true
}

// || RELATION ||

type RelationQueryOpt struct {
	Rel    string
	Fields []string
}

const relationQueryOptKey = "RelationQueryOpt"

func NewRelationQueryOpt(qr *QueryRequest, rel string, flds ...string) {
	rq := RelationQueryOpt{rel, flds}
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

type QueryRequestVariantOperations map[QueryVariant]ServiceOperation

func SwitchQueryRequestVariant(ctx context.Context, qr *QueryRequest, qrvo QueryRequestVariantOperations) error {
	op, ok := qrvo[qr.Variant]
	if !ok {
		panic(fmt.Sprintf("%s not supported for model %s", qr.Variant, qr.Model.Type().Name()))
	}
	return op(ctx, qr)
}
