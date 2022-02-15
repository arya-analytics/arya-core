package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type ServiceLocal struct {
	storage storage.Storage
}

func NewServiceLocal(storage storage.Storage) *ServiceLocal {
	return &ServiceLocal{storage: storage}
}

// |||| CHUNK ||||

func (s *ServiceLocal) Create(ctx context.Context, cc *model.Reflect) error {
	return s.storage.NewCreate().Model(cc.Pointer()).Exec(ctx)
}

func (s *ServiceLocal) Retrieve(ctx context.Context, cc *model.Reflect, pkC model.PKChain, omitBulk bool) error {
	return s.storage.NewRetrieve().Model(cc.Pointer()).WherePKs(pkC.Raw()).Exec(ctx)
}

func (s *ServiceLocal) Delete(ctx context.Context, ccPKC model.PKChain) error {
	return s.storage.NewDelete().Model(&storage.ChannelChunk{}).WherePK(ccPKC.Raw()).Exec(ctx)
}

// |||| REPLICA ||||

func (s *ServiceLocal) CreateReplicas(ctx context.Context, ccr *model.Reflect) error {
	return s.storage.NewCreate().Model(ccr.Pointer()).Exec(ctx)
}

func (s *ServiceLocal) RetrieveReplicas(ctx context.Context, ccr *model.Reflect, ccrPKC model.PKChain, omitBulk bool) error {
	return s.storage.NewRetrieve().Model(ccr.Pointer()).WherePKs(ccrPKC.Raw()).Exec(ctx)
}

func (s *ServiceLocal) DeleteReplicas(ctx context.Context, ccr *model.Reflect, ccrPKC model.PKChain) error {
	return s.storage.NewDelete().Model(&storage.ChannelChunkReplica{}).WherePKs(ccrPKC.Raw()).Exec(ctx)
}

// |||| RANGE REPLICA ||||

func (s *ServiceLocal) RetrieveRangeReplicas(ctx context.Context, rr *model.Reflect, rrPKC model.PKChain) error {
	return s.storage.NewRetrieve().
		Model(rr.Pointer()).
		WherePKs(rrPKC.Raw()).
		Relation("Node", "id", "address", "is_host").
		Exec(ctx)
}
