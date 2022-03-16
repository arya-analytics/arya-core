package rng

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
)

// QueryAssemble provides definitions for common queries that need to be performed on ranges. Extends query.Assemble
// interface.
//
// Retrieve, Update, and Delete query extensions provide a similar (yet more specific) API to query.Assemble.
// In contrast, Create query extensions return a concrete, created value. This is done to simplify the interface and
// enable the creation of multiple item combinations at once.
//
// Instantiate QueryAssemble by calling NewQueryAssemble.
type QueryAssemble struct {
	query.AssembleBase
}

// |||| CONSTRUCTOR |||

func NewQueryAssemble(exec query.Execute) *QueryAssemble {
	return &QueryAssemble{AssembleBase: query.NewAssemble(exec)}
}

// |||| CREATE ||||

// CreateRange creates the following items:
//
// 1. models.Range - Marked As Open
// 2. models.RangeReplica - The leaseholder replica assigned to the node with the given nodePK.
// 3. models.RangeLease - A lease bound to the range that marks the created replica as the leaseholder.
//
// Returns the created range with the models.RangeLease and models.RangeReplica assigned as nested fields, as well
// as any errors encountered during creation.
func (qa *QueryAssemble) CreateRange(ctx context.Context, nodePK int) (*models.Range, error) {
	c := errutil.NewCatchContext(ctx)
	r := &models.Range{ID: uuid.New(), Status: models.RangeStatusOpen}
	c.Exec(qa.NewCreate().Model(r).Exec)
	rr := &models.RangeReplica{ID: uuid.New(), RangeID: r.ID, NodeID: nodePK}
	c.Exec(qa.NewCreate().Model(rr).Exec)
	lease := &models.RangeLease{ID: uuid.New(), RangeID: r.ID, RangeReplicaID: rr.ID, RangeReplica: rr}
	c.Exec(qa.NewCreate().Model(lease).Exec)
	r.RangeLease = lease
	return r, c.Error()
}

// CreateRangeReplica creates a new models.RangeReplica bound to the range with the provided rngPK and assigned to the
// node with the given nodePK.
//
// Returns the created replica along with any errors encountered during creation.
func (qa *QueryAssemble) CreateRangeReplica(ctx context.Context, rngPK uuid.UUID, nodePK int) (*models.RangeReplica, error) {
	rr := &models.RangeReplica{ID: uuid.New(), RangeID: rngPK, NodeID: nodePK}
	return rr, qa.NewCreate().Model(rr).Exec(ctx)
}

// ||| RETRIEVE |||

// RetrieveRangeQuery retrieves the range with the given primary key.
// Also retrieves the relations specified in retrieveRangeLeaseReplicaRelationQuery.
func (qa *QueryAssemble) RetrieveRangeQuery(rng *models.Range, pk uuid.UUID) *query.Retrieve {
	return qa.retrieveRangeLeaseReplicaRelationQuery(rng).WherePK(pk)
}

// RetrieveRangeReplicasQuery retrieves all range replicas bound to the range with PK rngPK.
func (qa *QueryAssemble) RetrieveRangeReplicasQuery(rr *[]*models.RangeReplica, rngPK uuid.UUID) *query.Retrieve {
	return qa.NewRetrieve().Model(rr).WhereFields(query.WhereFields{"RangeID": rngPK})
}

// RetrieveOpenRangesQuery retrieves all ranges marked as open.
func (qa *QueryAssemble) RetrieveOpenRangesQuery(rng *[]*models.Range) *query.Retrieve {
	return qa.retrieveRangeLeaseReplicaRelationQuery(rng).WhereFields(query.WhereFields{"Status": models.RangeStatusOpen})
}

// RetrieveRangeChunksQuery retrieves all models.ChannelChunk that belong to the range with PK rngPK.
func (qa *QueryAssemble) RetrieveRangeChunksQuery(chunks *[]*models.ChannelChunk, rngPK uuid.UUID) *query.Retrieve {
	return qa.NewRetrieve().Model(chunks).WhereFields(query.WhereFields{"RangeID": rngPK})
}

// RetrieveRangeSizeQuery calculates the size of the models.Range with the given PK. Binds the calculated
// size into the 'into' arg.
func (qa *QueryAssemble) RetrieveRangeSizeQuery(pk uuid.UUID, into *int64) *query.Retrieve {
	return qa.RetrieveRangeChunksQuery(&[]*models.ChannelChunk{}, pk).Calc(query.CalcSum, "Size", into)
}

// RetrieveRangeChunkReplicasQuery retrieves all models.ChannelChunkReplica that belong to the range with PK rngPK.
// It's typical for this query to return large numbers of objects, so only retrieves "ID" and "RangeReplicaID" fields.
func (qa *QueryAssemble) RetrieveRangeChunkReplicasQuery(ccr *[]*models.ChannelChunkReplica, rngPK uuid.UUID) *query.Retrieve {
	return qa.NewRetrieve().
		Model(ccr).
		WhereFields(query.WhereFields{"ChannelChunk.RangeID": rngPK}).
		Fields("ID", "RangeReplicaID")
}

// |||| RE-ALLOCATE ||||

// ReallocateChunkReplicasQuery reassigns all models.ChannelChunkReplica with primary keys in pks, so
// they belong to the models.RangeReplica with PK rrPK.
func (qa *QueryAssemble) ReallocateChunkReplicasQuery(pks []uuid.UUID, rrPK uuid.UUID) *query.Update {
	var ccr []*models.ChannelChunkReplica
	for _, pk := range pks {
		ccr = append(ccr, &models.ChannelChunkReplica{ID: pk, RangeReplicaID: rrPK})
	}
	return qa.NewUpdate().Model(&ccr).Fields("RangeReplicaID").Bulk()
}

// ReallocateChunksQuery reassigns all models.ChannelChunk with primary keys in pks, so they belong to the
// models.Range with PK rngPK.
func (qa *QueryAssemble) ReallocateChunksQuery(pks []uuid.UUID, rngPK uuid.UUID) *query.Update {
	var cc []*models.ChannelChunk
	for _, pk := range pks {
		cc = append(cc, &models.ChannelChunk{ID: pk, RangeID: rngPK})
	}
	return qa.NewUpdate().Model(&cc).Fields("RangeID").Bulk()
}

// |||| UPDATE ||||

// UpdateRangeStatusQuery query updates the status of models.Range with the provided PK to the provided status.
func (qa *QueryAssemble) UpdateRangeStatusQuery(pk uuid.UUID, status models.RangeStatus) *query.Update {
	return qa.NewUpdate().Model(&models.Range{ID: pk, Status: status}).Fields("Status").WherePK(pk)
}

func (qa *QueryAssemble) retrieveRangeLeaseReplicaRelationQuery(rng interface{}) *query.Retrieve {
	return qa.NewRetrieve().
		Model(rng).
		Relation("RangeLease", "ID", "RangeReplicaID").
		Relation("RangeLease.RangeReplica", "ID", "NodeID")
}
