package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type LocalReplicaRetrieveOpts struct {
	PKC       model.PKChain
	OmitBulk  bool
	Relations bool
}

type LocalReplicaDeleteOpts struct {
	PKC model.PKChain
}

type LocalReplicaUpdateOpts struct {
	PK     model.PK
	Fields []string
	Bulk   bool
}

type LocalRangeReplicaRetrieveOpts struct {
	PKC model.PKChain
}

type ServiceLocal interface {
	// |||| REPLICA ||||

	CreateReplica(ctx context.Context, chunkReplica interface{}) error

	RetrieveReplica(ctx context.Context, chunkReplica interface{}, opts LocalReplicaRetrieveOpts) error

	UpdateReplica(ctx context.Context, chunkReplica interface{}, opts LocalReplicaUpdateOpts) error

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
		q = q.Relation("RangeReplica", "ID").
			Relation("RangeReplica.Node", "ID", "Address", "IsHost")
	}
	if opts.OmitBulk {
		q = q.Fields("ID", "ChannelChunkID", "RangeReplicaID")
	}
	return q.Exec(ctx)
}

func (s *ServiceLocalStorage) DeleteReplica(ctx context.Context, opts LocalReplicaDeleteOpts) error {
	q := s.storage.NewDelete().Model(&models.ChannelChunkReplica{})
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	return q.Exec(ctx)
}

func (s *ServiceLocalStorage) UpdateReplica(ctx context.Context, chunkReplica interface{}, opts LocalReplicaUpdateOpts) error {
	q := s.storage.NewUpdate().Model(chunkReplica)
	if len(opts.Fields) > 0 {
		q.Fields(opts.Fields...)
	}
	if opts.Bulk {
		q.Bulk()
	}
	if opts.PK.Raw() != nil {
		q = q.WherePK(opts.PK.Raw())
	}
	return q.Exec(ctx)
}

// |||| RANGE REPLICA ||||

func (s *ServiceLocalStorage) RetrieveRangeReplica(ctx context.Context, rangeReplica interface{}, opts LocalRangeReplicaRetrieveOpts) error {
	q := s.storage.NewRetrieve().
		Model(rangeReplica).
		Relation("Node", "ID", "Address", "IsHost", "RPCPort")
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	return q.Exec(ctx)
}
