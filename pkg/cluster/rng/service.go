package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"reflect"
)

type Service struct {
	local ServiceLocal
}

func NewService(local ServiceLocal) *Service {
	return &Service{local: local}
}

func (s *Service) CanHandle(q *internal.QueryRequest) bool {
	return catalog().Contains(q.Model.Pointer())
}

func (s *Service) Exec(ctx context.Context, qr *internal.QueryRequest) error {
	switch qr.Model.Type() {
	case reflect.TypeOf(storage.Range{}):
		return internal.SwitchQueryRequestVariant(ctx, qr, internal.QueryRequestVariantOperations{
			internal.QueryVariantCreate:   s.createRange,
			internal.QueryVariantRetrieve: s.retrieveRange,
		})
	case reflect.TypeOf(storage.RangeLease{}):
		return internal.SwitchQueryRequestVariant(ctx, qr, internal.QueryRequestVariantOperations{})
	case reflect.TypeOf(storage.RangeReplica{}):
		return internal.SwitchQueryRequestVariant(ctx, qr, internal.QueryRequestVariantOperations{})
	default:
		panic("range service received an unknown model type!")
	}
}

// |||| RANGE ||||

func (s *Service) createRange(ctx context.Context, qr *internal.QueryRequest) error {
	return s.local.CreateRange(ctx, qr.Model.Pointer())
}

const (
	RangeLeaseNodePkFldName = "RangeLease.RangeReplica.NodeID"
)

func (s *Service) retrieveRange(ctx context.Context, qr *internal.QueryRequest) error {
	PKC, ok := internal.PKQueryOpt(qr)
	if ok {
		return s.local.RetrieveRange(ctx, qr.Model.Pointer(), LocalRangeRetrieveOpts{PKC: PKC})
	}
	flds := internal.FieldsQueryOpt(qr)
	leaseNodePK, ok := flds.Retrieve(RangeLeaseNodePkFldName)
	if ok {
		return s.local.RetrieveRange(
			ctx,
			qr.Model.Pointer(),
			LocalRangeRetrieveOpts{LeaseNodePK: model.NewPK(leaseNodePK)},
		)
	}
	panic("retrieve range received insufficient query opts")
}

// |||| LEASE ||||

func (s *Service) createLease(ctx context.Context, qr *internal.QueryRequest) error {
	return s.local.CreateLease(ctx, qr.Model.Pointer())
}

const (
	LeaseRangeReplicaPKFldName = "RangeReplica.ID"
	LeaseRangePKFldName        = "RangeID"
)

func (s *Service) retrieveLease(ctx context.Context, qr *internal.QueryRequest) error {
	PKC, ok := internal.PKQueryOpt(qr)
	if ok {
		if len(PKC) > 1 {
			panic("lease retrieve query cannot have more than one primary key!")
		}
		return s.local.RetrieveLease(ctx, qr.Model.Pointer(), LocalLeaseRetrieveOpts{PK: PKC[0]})
	}
	flds := internal.FieldsQueryOpt(qr)
	rangePK, ok := flds.Retrieve(LeaseRangePKFldName)
	if ok {
		return s.local.RetrieveLease(ctx, qr.Model.Pointer(), LocalLeaseRetrieveOpts{RangePK: model.NewPK(rangePK)})
	}
	rangeReplicaPK, ok := flds.Retrieve(LeaseRangeReplicaPKFldName)
	if ok {
		return s.local.RetrieveLease(ctx, qr.Model.Pointer(), LocalLeaseRetrieveOpts{RangeReplicaPK: model.NewPK(rangeReplicaPK)})
	}
	panic("retrieve range lease retrieved insufficient query opts")
}

func catalog() model.Catalog {
	return model.Catalog{
		&storage.Range{},
		&storage.RangeLease{},
		&storage.RangeReplica{},
	}
}
