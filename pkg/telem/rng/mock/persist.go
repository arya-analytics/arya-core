package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Persist struct {
	Ranges        []*models.Range
	RangeReplicas []*models.RangeReplica
	RangeLeases   []*models.RangeLease
	Chunks        []*models.ChannelChunk
	ChunkReplicas []*models.ChannelChunkReplica
}

func NewBlankPersist() *Persist {
	return &Persist{
		Ranges:        []*models.Range{},
		RangeReplicas: []*models.RangeReplica{},
		RangeLeases:   []*models.RangeLease{},
		Chunks:        []*models.ChannelChunk{},
		ChunkReplicas: []*models.ChannelChunkReplica{},
	}
}

func (p *Persist) NewRange(ctx context.Context, nodeID int) (*models.Range, error) {
	id := uuid.New()
	rr := &models.RangeReplica{
		ID:      uuid.New(),
		RangeID: id,
		NodeID:  nodeID,
	}
	lease := &models.RangeLease{
		ID:             uuid.New(),
		RangeID:        id,
		RangeReplica:   rr,
		RangeReplicaID: rr.ID,
	}
	r := &models.Range{
		ID:         id,
		Status:     models.RangeStatusOpen,
		RangeLease: lease,
	}
	p.Ranges = append(p.Ranges, r)
	p.RangeReplicas = append(p.RangeReplicas, rr)
	p.RangeLeases = append(p.RangeLeases, lease)
	return r, nil
}

func (p *Persist) CreateRange(ctx context.Context, rng interface{}) error {
	model.NewReflect(rng).ForEach(func(rfl *model.Reflect, i int) {
		p.Ranges = append(p.Ranges, rfl.Pointer().(*models.Range))
	})
	return nil
}

func (p *Persist) CreateRangeLease(ctx context.Context, rngLease interface{}) error {
	model.NewReflect(rngLease).ForEach(func(rfl *model.Reflect, i int) {
		p.RangeLeases = append(p.RangeLeases, rfl.Pointer().(*models.RangeLease))
	})
	return nil
}

func (p *Persist) CreateRangeReplica(ctx context.Context, rngReplica interface{}) error {
	model.NewReflect(rngReplica).ForEach(func(rfl *model.Reflect, i int) {
		p.RangeReplicas = append(p.RangeReplicas, rfl.Pointer().(*models.RangeReplica))
	})
	return nil
}

func (p *Persist) RetrieveRange(ctx context.Context, ID uuid.UUID) (*models.Range, error) {
	for _, rng := range p.Ranges {
		if rng.ID == ID {
			return rng, nil
		}
	}
	return nil, storage.Error{Type: storage.ErrorTypeItemNotFound}
}

func (p *Persist) RetrieveRangeSize(ctx context.Context, ID uuid.UUID) (int64, error) {
	var size int64
	for _, cc := range p.Chunks {
		if cc.RangeID == ID {
			size += cc.Size
		}
	}
	return size, nil
}

func (p *Persist) RetrieveRangeChunks(ctx context.Context, rangeID uuid.UUID) ([]*models.ChannelChunk, error) {
	var chunks []*models.ChannelChunk
	for _, cc := range p.Chunks {
		if cc.RangeID == rangeID {
			chunks = append(chunks, cc)
		}
	}
	return chunks, nil
}

func (p *Persist) RetrieveRangeChunkReplicas(ctx context.Context, rangeID uuid.UUID) ([]*models.ChannelChunkReplica, error) {
	var chunkReplicas []*models.ChannelChunkReplica
	for _, ccr := range p.ChunkReplicas {
		if ccr.ChannelChunk.RangeID == rangeID {
			chunkReplicas = append(chunkReplicas, ccr)
		}
	}
	return chunkReplicas, nil
}

func (p *Persist) ReallocateChunks(ctx context.Context, pks interface{}, newRangeID uuid.UUID) error {
	rng, err := p.RetrieveRange(ctx, newRangeID)
	if err != nil {
		return err
	}
	for _, PK := range model.NewPKChain(pks) {
		for _, cc := range p.Chunks {
			if model.NewPK(cc.ID).Equals(PK) {
				cc.RangeID = newRangeID
				cc.Range = rng
			}
		}
	}
	return nil
}

func (p *Persist) ReallocateChunkReplicas(ctx context.Context, pks interface{}, newRangeReplicaID uuid.UUID) error {
	var replica *models.RangeReplica
	for _, rr := range p.RangeReplicas {
		if rr.ID == newRangeReplicaID {
			replica = rr
		}
	}
	if replica == nil {
		return storage.Error{Type: storage.ErrorTypeItemNotFound}
	}
	for _, PK := range model.NewPKChain(pks) {
		for _, ccr := range p.ChunkReplicas {
			if model.NewPK(ccr.ID).Equals(PK) {
				ccr.RangeReplicaID = newRangeReplicaID
				ccr.RangeReplica = replica
			}
		}
	}
	return nil
}

func (p *Persist) RetrieveRangeReplicas(ctx context.Context, rngID uuid.UUID) ([]*models.RangeReplica, error) {
	var rangeReplicas []*models.RangeReplica
	for _, rr := range p.RangeReplicas {
		if rr.RangeID == rngID {
			rangeReplicas = append(rangeReplicas, rr)
		}
	}
	return rangeReplicas, nil
}

func NewPersistOverallocatedRange() (uuid.UUID, *Persist) {
	rng := &models.Range{
		ID: uuid.New(),
	}
	var rangeReplicas []*models.RangeReplica
	for i := 0; i < 3; i++ {
		rangeReplicas = append(rangeReplicas, &models.RangeReplica{
			ID:      uuid.New(),
			RangeID: rng.ID,
			NodeID:  i + 1,
		})
	}
	lease := &models.RangeLease{
		RangeID:        rng.ID,
		RangeReplicaID: rangeReplicas[0].ID,
		RangeReplica:   rangeReplicas[0],
	}

	rng.RangeLease = lease

	rand.Seed(time.Now().UnixNano())

	var size int64 = 0
	var chunks []*models.ChannelChunk
	var chunkReplicas []*models.ChannelChunkReplica
	for float64(size) < float64(models.MaxRangeSize)*1.25 {
		chunkSize := rand.Int63n(int64(float64(models.MaxRangeSize) * 0.02))
		size += chunkSize
		chunk := &models.ChannelChunk{
			ID:      uuid.New(),
			RangeID: rng.ID,
			Size:    size,
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
	p := &Persist{
		Ranges:        []*models.Range{rng},
		RangeReplicas: rangeReplicas,
		Chunks:        chunks,
		ChunkReplicas: chunkReplicas,
		RangeLeases:   []*models.RangeLease{lease},
	}
	return rng.ID, p
}
