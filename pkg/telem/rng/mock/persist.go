package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/google/uuid"
	"math/rand"
	"reflect"
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

func (p *Persist) Exec(ctx context.Context, qp *query.Pack) error {
	return query.Switch(ctx, qp, query.Ops{
		Create: p.create,
	})
}

func (p *Persist) create(ctx context.Context, qp *query.Pack) error {
	switch qp.Model().Type() {
	case reflect.TypeOf(models.Range{}):
		return p.createRange(ctx, qp)
	case reflect.TypeOf(models.RangeReplica{}):
		return p.createRangeReplica(ctx, qp)
	case reflect.TypeOf(models.RangeLease{}):
		return p.createRangeLease(ctx, qp)
	default:
		panic("mock create received unknown model")
	}
}

func (p *Persist) createRange(ctx context.Context, qp *query.Pack) error {
	model.NewReflect(p.Ranges).ChainAppendEach(qp.Model())
	return nil
}

func (p *Persist) createRangeReplica(ctx context.Context, qp *query.Pack) error {
	model.NewReflect(p.RangeReplicas).ChainAppendEach(qp.Model())
	return nil
}

func (p *Persist) createRangeLease(ctx context.Context, qp *query.Pack) error {
	model.NewReflect(p.RangeLeases).ChainAppendEach(qp.Model())
	return nil
}

func (p *Persist) retrieveRange(ctx context.Context, qp *query.Pack) error {
	pkc, ok := query.PKOpt(qp)
	if ok {
		return p.retrieveRangesByPK(qp, pkc)
	}
}

func (p *Persist) retrieveRangesByPK(qp *query.Pack, pkc model.PKChain) error {
	exc := model.NewExchange(qp.Model().Pointer(), &[]*models.Range{})
	for _, rng := range p.Ranges {
		for _, pk := range pkc {
			if model.NewPK(rng.ID).Equals(pk) {
				exc.Dest().ChainAppend(model.NewReflect(rng))
			}
		}
	}
	exc.ToSource()
	return nil
}

func retrieveRangesByStatus(ctx context.Context) ([]*models.Range, error) {
	var ranges []*models.Range
	for _, rng := range p.Ranges {
		if rng.Status == models.RangeStatusOpen {
			ranges = append(ranges, rng)
		}
	}
	return ranges, nil
}

func (p *Persist) RetrieveRangeSize(ctx context.Context, PK uuid.UUID) (int64, error) {
	var size int64 = 0
	for _, cc := range p.Chunks {
		if cc.RangeID == PK {
			size += cc.Size
		}
	}
	return size, nil
}

func (p *Persist) RetrieveRangeChunks(ctx context.Context, rngPK uuid.UUID) ([]*models.ChannelChunk, error) {
	var chunks []*models.ChannelChunk
	for _, cc := range p.Chunks {
		if cc.RangeID == rngPK {
			chunks = append(chunks, cc)
		}
	}
	return chunks, nil
}

func (p *Persist) RetrieveRangeChunkReplicas(ctx context.Context, rngPK uuid.UUID) ([]*models.ChannelChunkReplica, error) {
	var chunkReplicas []*models.ChannelChunkReplica
	for _, ccr := range p.ChunkReplicas {
		if ccr.ChannelChunk.RangeID == rngPK {
			chunkReplicas = append(chunkReplicas, ccr)
		}
	}
	return chunkReplicas, nil
}

func (p *Persist) ReallocateChunks(ctx context.Context, pks []uuid.UUID, newRngPK uuid.UUID) error {
	rng, err := p.RetrieveRange(ctx, newRngPK)
	if err != nil {
		return err
	}
	for _, PK := range model.NewPKChain(pks) {
		for _, cc := range p.Chunks {
			if model.NewPK(cc.ID).Equals(PK) {
				cc.RangeID = newRngPK
				cc.Range = rng
			}
		}
	}
	return nil
}

func (p *Persist) ReallocateChunkReplicas(ctx context.Context, pks []uuid.UUID, newRRPK uuid.UUID) error {
	var replica *models.RangeReplica
	for _, rr := range p.RangeReplicas {
		if rr.ID == newRRPK {
			replica = rr
		}
	}
	if replica == nil {
		return query.Error{Type: query.ErrorTypeItemNotFound}
	}
	for _, PK := range model.NewPKChain(pks) {
		for _, ccr := range p.ChunkReplicas {
			if model.NewPK(ccr.ID).Equals(PK) {
				ccr.RangeReplicaID = newRRPK
				ccr.RangeReplica = replica
			}
		}
	}
	return nil
}

func (p *Persist) RetrieveRangeReplicas(ctx context.Context, rngPK uuid.UUID) ([]*models.RangeReplica, error) {
	var rangeReplicas []*models.RangeReplica
	for _, rr := range p.RangeReplicas {
		if rr.RangeID == rngPK {
			rangeReplicas = append(rangeReplicas, rr)
		}
	}
	return rangeReplicas, nil
}

func (p *Persist) UpdateRangeStatus(ctx context.Context, PK uuid.UUID, status models.RangeStatus) error {
	found := false
	for _, r := range p.Ranges {
		if r.ID == PK {
			r.Status = status
			found = true
		}
	}
	if !found {
		return query.Error{Type: query.ErrorTypeItemNotFound}
	}
	return nil
}

func NewPersistOverallocatedRange() (*Persist, uuid.UUID) {
	rng := &models.Range{
		ID:     uuid.New(),
		Status: models.RangeStatusOpen,
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
		chunkSize := rand.Int63n(models.MaxChunkSize)
		size += chunkSize
		chunk := &models.ChannelChunk{
			ID:      uuid.New(),
			RangeID: rng.ID,
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
	p := &Persist{
		Ranges:        []*models.Range{rng},
		RangeReplicas: rangeReplicas,
		Chunks:        chunks,
		ChunkReplicas: chunkReplicas,
		RangeLeases:   []*models.RangeLease{lease},
	}
	return p, rng.ID
}
