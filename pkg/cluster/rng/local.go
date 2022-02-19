package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type LocalRangeRetrieveOpts struct {
	PKC         model.PKChain
	LeaseNodePK model.PK
}

type LocalLeaseRetrieveOpts struct {
	PK             model.PK
	RangePK        model.PK
	RangeReplicaPK model.PK
}

type ServiceLocal interface {
	// |||| RANGE ||||

	CreateRange(ctx context.Context, rng interface{}) error
	RetrieveRange(ctx context.Context, rng interface{}, opts LocalRangeRetrieveOpts) error

	// |||| RANGE LEASE ||||

	CreateLease(ctx context.Context, lease interface{}) error
	RetrieveLease(ctx context.Context, lease interface{}, opts LocalLeaseRetrieveOpts) error
}

type ServiceLocalStorage struct {
	storage storage.Storage
}

func NewServiceLocalStorage(storage storage.Storage) ServiceLocal {
	return &ServiceLocalStorage{storage: storage}
}

// |||| RANGE ||||

func (s *ServiceLocalStorage) CreateRange(ctx context.Context, rng interface{}) error {
	return s.storage.NewCreate().Model(rng).Exec(ctx)
}

func (s *ServiceLocalStorage) RetrieveRange(ctx context.Context, rng interface{}, opts LocalRangeRetrieveOpts) error {
	q := s.storage.NewRetrieve().Model(rng)
	if !opts.PKC.AllZero() {
		q.WherePKs(opts.PKC.Raw())
	} else if !opts.LeaseNodePK.IsZero() {
		q.WhereFields(storage.Fields{RangeLeaseNodePkFldName: opts.LeaseNodePK.Raw()})
	}
	return q.Exec(ctx)
}

// |||| RANGE LEASE ||||

func (s *ServiceLocalStorage) CreateLease(ctx context.Context, lease interface{}) error {
	return s.storage.NewCreate().Model(lease).Exec(ctx)
}

func (s *ServiceLocalStorage) RetrieveLease(ctx context.Context, lease interface{}, opts LocalLeaseRetrieveOpts) error {
	q := s.storage.NewRetrieve().Model(lease)
	if !opts.RangePK.IsZero() {
		q.WhereFields(storage.Fields{LeaseRangePKFldName: opts.RangePK.Raw()})
	}
	if !opts.RangeReplicaPK.IsZero() {
		q.WhereFields(storage.Fields{LeaseRangeReplicaPKFldName: opts.RangeReplicaPK.Raw()})
	}
	return q.Exec(ctx)
}

// |||| RANGE REPLICA ||||
