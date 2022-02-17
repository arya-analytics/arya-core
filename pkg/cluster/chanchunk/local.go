package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

// |||| INTERFACE ||||

type LocalChunkRetrieveOpts struct {
	PKC model.PKChain
}

type LocalChunkDeleteOpts struct {
	PKC model.PKChain
}

type LocalReplicaRetrieveOpts struct {
	PKC       model.PKChain
	OmitBulk  bool
	Relations bool
}

type LocalReplicaDeleteOpts struct {
	PKC model.PKChain
}

type LocalRangeReplicaRetrieveOpts struct {
	PKC model.PKChain
}

type ServiceLocal interface {
	// |||| CHUNK ||||

	CreateChunk(ctx context.Context, chunk interface{}) error

	RetrieveChunk(ctx context.Context, chunk interface{}, opts LocalChunkRetrieveOpts) error

	DeleteChunk(ctx context.Context, opts LocalChunkDeleteOpts) error

	// |||| REPLICA ||||

	CreateReplica(ctx context.Context, chunkReplica interface{}) error

	RetrieveReplica(ctx context.Context, chunkReplica interface{}, opts LocalReplicaRetrieveOpts) error

	DeleteReplica(ctx context.Context, opts LocalReplicaDeleteOpts) error

	// |||| RANGE REPLICA ||||

	RetrieveRangeReplica(ctx context.Context, rangeReplica interface{}, opts LocalRangeReplicaRetrieveOpts) error
}

// |||| LOCAL STORAGE IMPLEMENTATION ||||

type ServiceLocalStorage struct {
	storage storage.Storage
}

func NewServiceLocalStorage(storage storage.Storage) ServiceLocal {
	return &ServiceLocalStorage{storage: storage}
}

// |||| CHUNK ||||

func (s *ServiceLocalStorage) CreateChunk(ctx context.Context, chunk interface{}) error {
	return s.storage.NewCreate().Model(chunk).Exec(ctx)
}

func (s *ServiceLocalStorage) RetrieveChunk(ctx context.Context, chunk interface{}, opts LocalChunkRetrieveOpts) error {
	q := s.storage.NewRetrieve().Model(chunk)
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	return q.Exec(ctx)
}

func (s *ServiceLocalStorage) DeleteChunk(ctx context.Context, opts LocalChunkDeleteOpts) error {
	q := s.storage.NewDelete().Model(&storage.ChannelChunk{})
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	return q.Exec(ctx)
}

// |||| REPLICA ||||

func (s *ServiceLocalStorage) CreateReplica(ctx context.Context, chunkReplica interface{}) error {
	return s.storage.NewCreate().Model(chunkReplica).Exec(ctx)
}

func (s *ServiceLocalStorage) RetrieveReplica(ctx context.Context, chunkReplica interface{}, opts LocalReplicaRetrieveOpts) error {
	q := s.storage.NewRetrieve().Model(chunkReplica)
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	if opts.Relations {
		q = q.Relation("RangeReplica", "id").
			Relation("RangeReplica.Node", "id", "address", "is_host")
	}
	return q.Exec(ctx)
}

func (s *ServiceLocalStorage) DeleteReplica(ctx context.Context, opts LocalReplicaDeleteOpts) error {
	q := s.storage.NewDelete().Model(&storage.ChannelChunkReplica{})
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	return q.Exec(ctx)
}

// |||| RANGE REPLICA ||||

func (s *ServiceLocalStorage) RetrieveRangeReplica(ctx context.Context, rangeReplica interface{}, opts LocalRangeReplicaRetrieveOpts) error {
	q := s.storage.NewRetrieve().
		Model(rangeReplica).
		Relation("Node", "id", "address", "is_host")
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	return q.Exec(ctx)
}