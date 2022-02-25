package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type LocalReplicaRetrieveOpts struct {
	PKC         model.PKChain
	Fields      []string
	WhereFields model.WhereFields
	OmitBulk    bool
	Relations   bool
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

type Local interface {
	CreateReplica(ctx context.Context, chunkReplica interface{}) error
	RetrieveReplica(ctx context.Context, chunkReplica interface{}, opts LocalReplicaRetrieveOpts) error
	UpdateReplica(ctx context.Context, chunkReplica interface{}, opts LocalReplicaUpdateOpts) error
	DeleteReplica(ctx context.Context, opts LocalReplicaDeleteOpts) error
	RetrieveRangeReplica(ctx context.Context, rangeReplica interface{}, opts LocalRangeReplicaRetrieveOpts) error
}

// |||| LOCAL STORAGE IMPLEMENTATION ||||

type LocalStorage struct {
	storage storage.Storage
}

func NewServiceLocalStorage(storage storage.Storage) Local {
	return &LocalStorage{storage: storage}
}

// |||| REPLICA ||||

func (ls *LocalStorage) CreateReplica(ctx context.Context, chunkReplica interface{}) error {
	return ls.storage.NewCreate().Model(chunkReplica).Exec(ctx)
}

func (ls *LocalStorage) RetrieveReplica(ctx context.Context, chunkReplica interface{}, opts LocalReplicaRetrieveOpts) error {
	q := ls.storage.NewRetrieve().Model(chunkReplica)
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	if opts.Relations {
		q = q.Relation("RangeReplica", "ID").
			Relation("RangeReplica.Node", "ID", "Address", "IsHost")
	}
	if opts.WhereFields != nil {
		q = q.WhereFields(opts.WhereFields)
	}
	if opts.OmitBulk {
		q = q.Fields("ID", "ChannelChunkID", "RangeReplicaID")
	}
	return q.Exec(ctx)
}

func (ls *LocalStorage) DeleteReplica(ctx context.Context, opts LocalReplicaDeleteOpts) error {
	q := ls.storage.NewDelete().Model(&models.ChannelChunkReplica{})
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	return q.Exec(ctx)
}

func (ls *LocalStorage) UpdateReplica(ctx context.Context, chunkReplica interface{}, opts LocalReplicaUpdateOpts) error {
	q := ls.storage.NewUpdate().Model(chunkReplica)
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

func (ls *LocalStorage) RetrieveRangeReplica(ctx context.Context, rangeReplica interface{}, opts LocalRangeReplicaRetrieveOpts) error {
	q := ls.storage.NewRetrieve().
		Model(rangeReplica).
		Relation("Node", "ID", "Address", "IsHost", "RPCPort")
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	return q.Exec(ctx)
}
