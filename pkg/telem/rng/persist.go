package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
)

// || CREATE ||

func createRange(ctx context.Context, exec query.Execute, nodePK int) (*models.Range, error) {
	r := &models.Range{ID: uuid.New(), Status: models.RangeStatusOpen}
	c := errutil.NewCatchContext(ctx)
	c.Exec(query.NewCreate().BindExec(exec).Model(r).Exec)
	rr := &models.RangeReplica{ID: uuid.New(), RangeID: r.ID, NodeID: nodePK}
	c.Exec(query.NewCreate().BindExec(exec).Model(rr).Exec)
	lease := &models.RangeLease{ID: uuid.New(), RangeID: r.ID, RangeReplicaID: rr.ID, RangeReplica: rr}
	c.Exec(query.NewCreate().BindExec(exec).Model(lease).Exec)
	r.RangeLease = lease
	return r, c.Error()
}

func createRangeReplica(ctx context.Context, qExec query.Execute, rngPK uuid.UUID, nodePK int) (*models.RangeReplica, error) {
	rr := &models.RangeReplica{RangeID: rngPK, NodeID: nodePK}
	return rr, query.NewCreate().Model(rr).BindExec(qExec).Exec(ctx)
}

// || RETRIEVE ||

func retrieveRangeQuery(qExec query.Execute, rng *models.Range, pk uuid.UUID) *query.Retrieve {
	return rangeLeaseReplicaRelationQuery(qExec, rng).WherePK(pk)
}

func retrieveRangeReplicasQuery(qExec query.Execute, rr []*models.RangeReplica, rngPK uuid.UUID) *query.Retrieve {
	return query.NewRetrieve().BindExec(qExec).Model(&rr).WhereFields(query.WhereFields{"RangeID": rngPK})
}

func retrieveOpenRangesQuery(qExec query.Execute, rng []*models.Range) *query.Retrieve {
	return rangeLeaseReplicaRelationQuery(qExec, &rng).WhereFields(query.WhereFields{"Status": models.RangeStatusOpen})
}

func rangeLeaseReplicaRelationQuery(qExec query.Execute, rng interface{}) *query.Retrieve {
	return query.NewRetrieve().
		BindExec(qExec).
		Model(rng).
		Relation("RangeLease", "ID", "RangeReplicaID").
		Relation("RangeLease.RangeReplica", "ID", "NodeID")
}

func retrieveRangeChunksQuery(qExec query.Execute, chunks interface{}, pk uuid.UUID) *query.Retrieve {
	return query.NewRetrieve().BindExec(qExec).Model(chunks).WhereFields(query.WhereFields{"RangeID": pk})
}

func retrieveRangeSizeQuery(qExec query.Execute, pk uuid.UUID, into *int64) *query.Retrieve {
	return retrieveRangeChunksQuery(qExec, &models.ChannelChunk{}, pk).Calc(query.CalcSum, "Size", &into)
}

func updateRangeStatusQuery(qExec query.Execute, pk uuid.UUID, status models.RangeStatus) *query.Update {
	return query.NewUpdate().BindExec(qExec).Model(&models.Range{ID: pk, Status: status}).Fields("Status").WherePK(pk)
}

func retrieveRangeChunkReplicasQuery(qExec query.Execute, ccr []*models.ChannelChunkReplica, rngPK uuid.UUID) *query.Retrieve {
	return query.NewRetrieve().
		BindExec(qExec).
		Model(&ccr).
		WhereFields(query.WhereFields{"ChannelChunk.RangeID": rngPK}).
		Fields("ID", "RangeReplicaID")
}

// || RE-ALLOCATE ||

func reallocateChunkReplicasQuery(qExec query.Execute, pks []uuid.UUID, rrPK uuid.UUID) *query.Update {
	var ccr []*models.ChannelChunkReplica
	for _, pk := range pks {
		ccr = append(ccr, &models.ChannelChunkReplica{ID: pk, RangeReplicaID: rrPK})
	}
	return query.NewUpdate().BindExec(qExec).Model(&ccr).Fields("RangeReplicaID").Bulk()
}

func reallocateChunksQuery(qExec query.Execute, pks []uuid.UUID, rngPK uuid.UUID) *query.Update {
	var cc []*models.ChannelChunk
	for _, pk := range pks {
		cc = append(cc, &models.ChannelChunk{ID: pk, RangeID: rngPK})
	}
	return query.NewUpdate().BindExec(qExec).Model(&cc).Fields("RangeID").Bulk()
}
