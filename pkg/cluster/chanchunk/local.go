package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

// |||| INTERFACE ||||

type ServiceLocal interface {
	// |||| CHUNK ||||

	Create(ctx context.Context, cc *model.Reflect) error
	Retrieve(ctx context.Context, cc *model.Reflect, ccPKC model.PKChain) error
	Delete(ctx context.Context, ccPKC model.PKChain) error

	// |||| REPLICA ||||

	CreateReplicas(ctx context.Context, ccr *model.Reflect) error
	RetrieveReplicas(ctx context.Context, ccr *model.Reflect, ccrPKC model.PKChain, omitBulk bool) error
	DeleteReplicas(ctx context.Context, ccrPKC model.PKChain) error

	// |||| RANGE REPLICA ||||

	RetrieveRangeReplicas(ctx context.Context, rr *model.Reflect, rrPKC model.PKChain) error
}

// |||| LOCAL STORAGE IMPLEMENTATION ||||

type ServiceLocalStorage struct {
	storage storage.Storage
}

func NewServiceLocalStorage(storage storage.Storage) ServiceLocal {
	return &ServiceLocalStorage{storage: storage}
}

// |||| CHUNK ||||

func (s *ServiceLocalStorage) Create(ctx context.Context, cc *model.Reflect) error {
	return s.storage.NewCreate().Model(cc.Pointer()).Exec(ctx)
}

func (s *ServiceLocalStorage) Retrieve(ctx context.Context, cc *model.Reflect, pkC model.PKChain) error {
	return s.storage.NewRetrieve().Model(cc.Pointer()).WherePKs(pkC.Raw()).Exec(ctx)
}

func (s *ServiceLocalStorage) Delete(ctx context.Context, ccPKC model.PKChain) error {
	return s.storage.NewDelete().Model(&storage.ChannelChunk{}).WherePKs(ccPKC.Raw()).Exec(ctx)
}

// |||| REPLICA ||||

func (s *ServiceLocalStorage) CreateReplicas(ctx context.Context, ccr *model.Reflect) error {
	return s.storage.NewCreate().Model(ccr.Pointer()).Exec(ctx)
}

func (s *ServiceLocalStorage) RetrieveReplicas(ctx context.Context, ccr *model.Reflect, ccrPKC model.PKChain, omitBulk bool) error {
	return s.storage.NewRetrieve().Model(ccr.Pointer()).WherePKs(ccrPKC.Raw()).Exec(ctx)
}

func (s *ServiceLocalStorage) DeleteReplicas(ctx context.Context, ccrPKC model.PKChain) error {
	return s.storage.NewDelete().Model(&storage.ChannelChunkReplica{}).WherePKs(ccrPKC.Raw()).Exec(ctx)
}

// |||| RANGE REPLICA ||||

func (s *ServiceLocalStorage) RetrieveRangeReplicas(ctx context.Context, rr *model.Reflect, rrPKC model.PKChain) error {
	return s.storage.NewRetrieve().
		Model(rr.Pointer()).
		WherePKs(rrPKC.Raw()).
		Relation("Node", "id", "address", "is_host").
		Exec(ctx)
}
