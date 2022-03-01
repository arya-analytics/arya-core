package base

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Service struct {
	storage storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) CanHandle(p *query.Pack) bool {
	return catalog().Contains(p.Model().Pointer())
}

func (s *Service) Exec(ctx context.Context, p *query.Pack) error {
	return query.Switch(
		ctx,
		p,
		query.Ops{
			Create:   s.create,
			Retrieve: s.retrieve,
			Update:   s.update,
			Delete:   s.delete,
		})
}

func (s *Service) create(ctx context.Context, p *query.Pack) error {
	q := s.storage.NewCreate().Model(p.Model().Pointer())
	return q.Exec(ctx)
}

func (s *Service) retrieve(ctx context.Context, p *query.Pack) error {
	q := s.storage.NewRetrieve().Model(p.Model().Pointer())

	// PK

	PKC, ok := query.PKOpt(p)
	if ok {
		q = q.WherePKs(PKC.Raw())
	}

	// WHERE FIELDS

	wFlds, ok := query.WhereFieldsOpt(p)
	if ok {
		q = q.WhereFields(wFlds)
	}

	// FIELDS

	flds, ok := query.RetrieveFieldsOpt(p)
	if ok {
		q = q.Fields(flds...)
	}

	// RELATIONS

	for _, rel := range query.RelationOpts(p) {
		q = q.Relation(rel.Rel, rel.Fields...)
	}

	// CALCULATIONS

	//calc, ok := query.RetrieveCalcOpt(p)
	//if ok {
	//	q = q.Calculate(calc.Op, calc.FldName, calc.Into)
	//}

	return q.Exec(ctx)
}

func (s *Service) delete(ctx context.Context, p *query.Pack) error {
	q := s.storage.NewDelete().Model(p.Model().Pointer())
	PKC, ok := query.PKOpt(p)
	if ok {
		q = q.WherePKs(PKC.Raw())
	}
	return q.Exec(ctx)
}

func (s *Service) update(ctx context.Context, p *query.Pack) error {
	q := s.storage.NewUpdate().Model(p.Model().Pointer())

	// PK

	PKC, ok := query.PKOpt(p)
	if len(PKC) > 1 {
		panic("update queries can't have more than one pk!")
	}
	if ok {
		q = q.WherePK(PKC[0].Raw())
	}

	// FIELDS

	flds, ok := query.RetrieveFieldsOpt(p)
	if ok {
		q = q.Fields(flds...)
	}

	// BULK

	bulkOpt := query.BulkUpdateOpt(p)
	if bulkOpt {
		q = q.Bulk()
	}

	return q.Exec(ctx)
}
