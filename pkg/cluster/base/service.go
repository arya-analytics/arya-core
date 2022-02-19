package base

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type Service struct {
	storage storage.Storage
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
	PKC, ok := internal.PKQueryOpt(qr)
	if ok {
		q.WherePKs(PKC.Raw())
	}
	//flds, ok := internal.FieldsQueryOpt()
	//if ok {
	//	q.WhereFields(flds)
	//}
	return q.Exec(ctx)
}

func (s *Service) delete(ctx context.Context, qr *internal.QueryRequest) error {
	q := s.storage.NewDelete().Model(qr.Model.Pointer())
	return q.Exec(ctx)
}

func (s *Service) update(ctx context.Context, qr *internal.QueryRequest) error {
	q := s.storage.NewUpdate().Model(qr.Model.Pointer())
	return q.Exec(ctx)
}

func catalog() model.Catalog {
	return model.Catalog{
		&storage.Node{},
		&storage.Range{},
		&storage.RangeLease{},
		&storage.ChannelConfig{},
	}
}
