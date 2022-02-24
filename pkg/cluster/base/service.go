package base

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type Service struct {
	storage storage.Storage
}

func NewService(s storage.Storage) *Service {
	return &Service{storage: s}
}

func (s *Service) CanHandle(qr *internal.QueryRequest) bool {
	return catalog().Contains(qr.Model.Pointer())
}

func (s *Service) Exec(ctx context.Context, qr *internal.QueryRequest) error {
	return internal.SwitchQueryRequestVariant(
		ctx,
		qr,
		internal.QueryRequestVariantOperations{
			internal.QueryVariantCreate:   s.create,
			internal.QueryVariantRetrieve: s.retrieve,
			internal.QueryVariantUpdate:   s.update,
			internal.QueryVariantDelete:   s.delete,
		})
}

func (s *Service) create(ctx context.Context, qr *internal.QueryRequest) error {
	q := s.storage.NewCreate().Model(qr.Model.Pointer())
	return q.Exec(ctx)
}

func (s *Service) retrieve(ctx context.Context, qr *internal.QueryRequest) error {
	q := s.storage.NewRetrieve().Model(qr.Model.Pointer())

	// PK

	PKC, ok := internal.PKQueryOpt(qr)
	if ok {
		q = q.WherePKs(PKC.Raw())
	}

	// WHERE FIELDS

	wFlds, ok := internal.WhereFieldsQueryOpt(qr)
	if ok {
		q = q.WhereFields(wFlds)
	}

	// FIELDS

	flds, ok := internal.RetrieveFieldsQueryOpt(qr)
	if ok {
		q = q.Fields(flds...)
	}

	// RELATIONS

	for _, rel := range internal.RelationQueryOpts(qr) {
		q = q.Relation(rel.Rel, rel.Fields...)
	}

	// CALCULATIONS

	calc, ok := internal.RetrieveCalculateQueryOpt(qr)
	if ok {
		q = q.Calculate(calc.C, calc.FldName, calc.Into)
	}

	return q.Exec(ctx)
}

func (s *Service) delete(ctx context.Context, qr *internal.QueryRequest) error {
	q := s.storage.NewDelete().Model(qr.Model.Pointer())
	PKC, ok := internal.PKQueryOpt(qr)
	if ok {
		q = q.WherePKs(PKC.Raw())
	}
	return q.Exec(ctx)
}

func (s *Service) update(ctx context.Context, qr *internal.QueryRequest) error {
	q := s.storage.NewUpdate().Model(qr.Model.Pointer())

	// PK

	PKC, ok := internal.PKQueryOpt(qr)
	if len(PKC) > 1 {
		panic("update queries can't have more than one pk!")
	}
	if ok {
		q = q.WherePK(PKC[0].Raw())
	}

	// BULK

	bulkOpt := internal.BulkUpdateQueryOpt(qr)
	if bulkOpt {
		q = q.Bulk()
	}

	return q.Exec(ctx)
}

func catalog() model.Catalog {
	return model.Catalog{
		&models.Node{},
		&models.Range{},
		&models.RangeReplica{},
		&models.RangeLease{},
		&models.ChannelConfig{},
		&models.ChannelChunk{},
	}
}
