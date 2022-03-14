package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func PopulateOverallocatedRange(
	ctx context.Context,
	asm query.Assemble,
) (uuid.UUID, []*models.RangeReplica, []*models.ChannelChunk, []*models.ChannelChunkReplica) {
	r := &models.Range{
		ID:     uuid.New(),
		Status: models.RangeStatusOpen,
	}
	var rangeReplicas []*models.RangeReplica
	for i := 0; i < 3; i++ {
		rangeReplicas = append(rangeReplicas, &models.RangeReplica{
			ID:      uuid.New(),
			RangeID: r.ID,
			NodeID:  i + 1,
		})
	}
	lease := &models.RangeLease{
		RangeID:        r.ID,
		RangeReplicaID: rangeReplicas[0].ID,
		RangeReplica:   rangeReplicas[0],
	}

	r.RangeLease = lease

	rand.Seed(time.Now().UnixNano())

	var size int64 = 0
	var chunks []*models.ChannelChunk
	var chunkReplicas []*models.ChannelChunkReplica
	for float64(size) < float64(models.MaxRangeSize)*1.25 {
		chunkSize := rand.Int63n(models.MaxChunkSize)
		size += chunkSize
		chunk := &models.ChannelChunk{
			ID:      uuid.New(),
			RangeID: r.ID,
			Size:    chunkSize,
		}
		for i := 0; i < 3; i++ {
			chunkReplicas = append(chunkReplicas,
				&models.ChannelChunkReplica{
					ID:             uuid.New(),
					ChannelChunk:   chunk,
					ChannelChunkID: chunk.ID,
					RangeReplica:   rangeReplicas[i],
					RangeReplicaID: rangeReplicas[i].ID},
			)
		}
		chunks = append(chunks, chunk)
	}
	c := errutil.NewCatchContext(ctx)

	c.Exec(asm.NewCreate().Model(r).Exec)
	c.Exec(asm.NewCreate().Model(&rangeReplicas).Exec)
	c.Exec(asm.NewCreate().Model(&chunks).Exec)
	c.Exec(asm.NewCreate().Model(&chunkReplicas).Exec)
	c.Exec(asm.NewCreate().Model(lease).Exec)

	return r.ID, rangeReplicas, chunks, chunkReplicas
}
