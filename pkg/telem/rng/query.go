package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
)

type QueryAssemble struct {
	query.AssembleBase
}

// || CREATE ||

func NewQueryAssemble(exec query.Execute) *QueryAssemble {
	return &QueryAssemble{AssembleBase: query.NewAssemble(exec)}
}

func (qa *QueryAssemble) CreateRange(ctx context.Context, nodePK int) (*models.Range, error) {
	r := &models.Range{ID: uuid.New(), Status: models.RangeStatusOpen}
	c := errutil.NewCatchContext(ctx)
	c.Exec(qa.NewCreate().Model(r).Exec)
	rr := &models.RangeReplica{ID: uuid.New(), RangeID: r.ID, NodeID: nodePK}
	c.Exec(qa.NewCreate().Model(rr).Exec)
	lease := &models.RangeLease{ID: uuid.New(), RangeID: r.ID, RangeReplicaID: rr.ID, RangeReplica: rr}
	c.Exec(qa.NewCreate().Model(lease).Exec)
	r.RangeLease = lease
	return r, c.Error()
}

func (qa *QueryAssemble) CreateRangeReplica(ctx context.Context, rngPK uuid.UUID, nodePK int) (*models.RangeReplica, error) {
	rr := &models.RangeReplica{RangeID: rngPK, NodeID: nodePK}
	return rr, qa.NewCreate().Model(rr).Exec(ctx)
}

// || RETRIEVE ||

func (qa *QueryAssemble) RetrieveRangeQuery(rng *models.Range, pk uuid.UUID) *query.Retrieve {
	return qa.RetrieveRangeLeaseReplicaRelationQuery(rng).WherePK(pk)
}

func (qa *QueryAssemble) RetrieveRangeReplicasQuery(rr []*models.RangeReplica, rngPK uuid.UUID) *query.Retrieve {
	return qa.NewRetrieve().Model(&rr).WhereFields(query.WhereFields{"RangeID": rngPK})
}

func (qa *QueryAssemble) RetrieveOpenRangesQuery(rng []*models.Range) *query.Retrieve {
	return qa.RetrieveRangeLeaseReplicaRelationQuery(&rng).WhereFields(query.WhereFields{"Status": models.RangeStatusOpen})
}

func (qa *QueryAssemble) RetrieveRangeLeaseReplicaRelationQuery(rng interface{}) *query.Retrieve {
	return qa.NewRetrieve().
		Model(rng).
		Relation("RangeLease", "ID", "RangeReplicaID").
		Relation("RangeLease.RangeReplica", "ID", "NodeID")
}

func (qa *QueryAssemble) RetrieveRangeChunksQuery(chunks interface{}, pk uuid.UUID) *query.Retrieve {
	return qa.NewRetrieve().Model(chunks).WhereFields(query.WhereFields{"RangeID": pk})
}

func (qa *QueryAssemble) RetrieveRangeSizeQuery(pk uuid.UUID, into *int64) *query.Retrieve {
	return qa.RetrieveRangeChunksQuery(&models.ChannelChunk{}, pk).Calc(query.CalcSum, "Size", into)
}

func (qa *QueryAssemble) UpdateRangeStatusQuery(pk uuid.UUID, status models.RangeStatus) *query.Update {
	return qa.NewUpdate().Model(&models.Range{ID: pk, Status: status}).Fields("Status").WherePK(pk)
}

func (qa *QueryAssemble) RetrieveRangeChunkReplicasQuery(ccr []*models.ChannelChunkReplica, rngPK uuid.UUID) *query.Retrieve {
	return qa.NewRetrieve().
		Model(&ccr).
		WhereFields(query.WhereFields{"ChannelChunk.RangeID": rngPK}).
		Fields("ID", "RangeReplicaID")
}

// || RE-ALLOCATE ||

func (qa *QueryAssemble) ReallocateChunkReplicasQuery(pks []uuid.UUID, rrPK uuid.UUID) *query.Update {
	var ccr []*models.ChannelChunkReplica
	for _, pk := range pks {
		ccr = append(ccr, &models.ChannelChunkReplica{ID: pk, RangeReplicaID: rrPK})
	}
	return qa.NewUpdate().Model(&ccr).Fields("RangeReplicaID").Bulk()
}

func (qa *QueryAssemble) ReallocateChunksQuery(pks []uuid.UUID, rngPK uuid.UUID) *query.Update {
	var cc []*models.ChannelChunk
	for _, pk := range pks {
		cc = append(cc, &models.ChannelChunk{ID: pk, RangeID: rngPK})
	}
	return query.NewUpdate().Model(&cc).Fields("RangeID").Bulk()
}
