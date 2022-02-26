package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
)

type LocalRetrieveOpts struct {
	PKC           model.PKChain
	Fields        []string
	WhereFields   model.WhereFields
	NodeRelations bool
}

type LocalDeleteOpts struct {
	PKC model.PKChain
}

type LocalUpdateOpts struct {
	PK     model.PK
	Fields []string
	Bulk   bool
}

type Local interface {
	Create(ctx context.Context, ccr interface{}) error
	Retrieve(ctx context.Context, ccr interface{}, opts LocalRetrieveOpts) error
	Update(ctx context.Context, ccr interface{}, opts LocalUpdateOpts) error
	Delete(ctx context.Context, opts LocalDeleteOpts) error
	RetrieveRangeReplica(ctx context.Context, rangeReplica interface{}, pkc model.PKChain) error
}

// |||| LOCAL STORAGE IMPLEMENTATION ||||

type LocalStorage struct {
	storage storage.Storage
}

func NewServiceLocalStorage(storage storage.Storage) Local {
	return &LocalStorage{storage: storage}
}

// |||| REPLICA ||||

func (ls *LocalStorage) Create(ctx context.Context, chunkReplica interface{}) error {
	return ls.storage.NewCreate().Model(chunkReplica).Exec(ctx)
}

func (ls *LocalStorage) Retrieve(ctx context.Context, chunkReplica interface{}, opts LocalRetrieveOpts) error {
	q := ls.storage.NewRetrieve().Model(chunkReplica)
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	if opts.NodeRelations {
		q = q.Relation("RangeReplica", "ID").
			Relation("RangeReplica.Node", "ID", "Address", "IsHost", "RPCPort")
	}
	if opts.WhereFields != nil {
		q = q.WhereFields(opts.WhereFields)
	}
	if len(opts.Fields) > 0 {
		q = q.Fields(opts.Fields...)
	}
	return q.Exec(ctx)
}

func (ls *LocalStorage) Delete(ctx context.Context, opts LocalDeleteOpts) error {
	q := ls.storage.NewDelete().Model(&models.ChannelChunkReplica{})
	if opts.PKC != nil {
		q = q.WherePKs(opts.PKC.Raw())
	}
	return q.Exec(ctx)
}

func (ls *LocalStorage) Update(ctx context.Context, chunkReplica interface{}, opts LocalUpdateOpts) error {
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

func (ls *LocalStorage) RetrieveRangeReplica(ctx context.Context, rangeReplica interface{}, pkc model.PKChain) error {
	return ls.storage.NewRetrieve().
		Model(rangeReplica).
		Relation("Node", "ID", "Address", "IsHost", "RPCPort").
		WherePKs(pkc.Raw()).
		Exec(ctx)
}
