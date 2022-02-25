package mock

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	"github.com/google/uuid"
	"time"
)

type Local struct {
	Nodes          []*models.Node
	ChannelConfigs []*models.ChannelConfig
	Chunks         []*models.ChannelChunk
	ChunkReplicas  []*models.ChannelChunkReplica
	RangeReplicas  []*models.RangeReplica
}

func NewPrepopulatedLocal() *Local {
	nodes := []*models.Node{{ID: 1, IsHost: true}, {ID: 2, IsHost: false}}
	rng := &models.Range{ID: uuid.New()}
	channelConfig := &models.ChannelConfig{
		ID:     uuid.New(),
		Node:   nodes[0],
		NodeID: nodes[0].ID,
	}
	rr := []*models.RangeReplica{
		{
			ID:      uuid.New(),
			RangeID: rng.ID,
			Range:   rng,
			NodeID:  nodes[0].ID,
			Node:    nodes[0],
		},
		{
			ID:      uuid.New(),
			RangeID: rng.ID,
			NodeID:  nodes[1].ID,
			Node:    nodes[1],
		},
	}
	lease := &models.RangeLease{
		ID:             uuid.New(),
		RangeReplica:   rr[0],
		RangeReplicaID: rr[0].ID,
	}
	rng.RangeLease = lease

	var (
		chunks        []*models.ChannelChunk
		chunkReplicas []*models.ChannelChunkReplica
	)
	for i := 0; i < 30; i++ {
		baseTlm := telem.NewBulk([]byte{})
		mock.TelemBulkPopulateRandomFloat64(baseTlm, 100)
		c := &models.ChannelChunk{
			ID:              uuid.New(),
			ChannelConfigID: channelConfig.ID,
			ChannelConfig:   channelConfig,
			Size:            baseTlm.Size(),
			StartTS:         time.Now().UnixMilli(),
		}
		chunks = append(chunks, c)
		for i := 0; i < 2; i++ {
			tlm := telem.NewBulk([]byte{})
			mock.TelemBulkPopulateRandomFloat64(tlm, 100)
			ccr := &models.ChannelChunkReplica{
				ID:             uuid.New(),
				ChannelChunkID: c.ID,
				ChannelChunk:   c,
				Telem:          tlm,
				RangeReplica:   rr[i],
				RangeReplicaID: rr[i].ID,
			}
			chunkReplicas = append(chunkReplicas, ccr)
		}
	}

	return &Local{
		Nodes:          nodes,
		RangeReplicas:  rr,
		ChannelConfigs: []*models.ChannelConfig{channelConfig},
		ChunkReplicas:  chunkReplicas,
		Chunks:         chunks,
	}
}

func (s *Local) CreateReplica(ctx context.Context, chunkReplica interface{}) error {
	rfl := model.NewReflect(chunkReplica)
	rfl.ForEach(func(rfl *model.Reflect, i int) {
		s.ChunkReplicas = append(s.ChunkReplicas, rfl.Pointer().(*models.ChannelChunkReplica))
	})
	return nil
}

func (s *Local) RetrieveReplica(ctx context.Context, chunkReplica interface{}, opts chanchunk.LocalReplicaRetrieveOpts) error {
	var chunkReplicas []*models.ChannelChunkReplica
	if opts.PKC != nil {
		chunkReplicas = s.retrieveReplicaByPK(opts.PKC)
	}
	if opts.WhereFields != nil {
		chunkReplicas = s.retrieveReplicaByWhereFields(opts.WhereFields)
	}
	sourceRfl := model.NewReflect(chunkReplica)
	if sourceRfl.IsChain() {
		model.NewExchange(chunkReplica, &chunkReplicas).ToSource()
	} else {
		if len(chunkReplicas) == 0 {
			return storage.Error{Type: storage.ErrorTypeItemNotFound}
		}
		model.NewExchange(chunkReplica, chunkReplicas[0]).ToSource()
	}
	return nil
}

func (s *Local) DeleteReplica(ctx context.Context, opts chanchunk.LocalReplicaDeleteOpts) error {
	if opts.PKC == nil {
		panic("delete queries require a primary key")
	}
	for _, PK := range opts.PKC {
		for i, ccr := range s.ChunkReplicas {
			if model.NewPK(ccr.ID).Equals(PK) {
				s.ChunkReplicas = append(s.ChunkReplicas[:i], s.ChunkReplicas[i+1:]...)
			}
		}
	}
	return nil
}

func (s *Local) UpdateReplica(ctx context.Context, chunkReplica interface{}, opts chanchunk.LocalReplicaUpdateOpts) error {
	uRfl := model.NewReflect(chunkReplica)
	uRfl.ForEach(func(updateRFL *model.Reflect, i int) {
		ccrs := s.retrieveReplicaByPK(uRfl.PKChain())
		exc := model.NewExchange(model.NewReflect(&ccrs).ChainValueByIndex(0).Pointer(), updateRFL.Pointer())
		exc.ToSource()
	})
	return nil
}

func (s *Local) RetrieveRangeReplica(ctx context.Context, rangeReplica interface{}, opts chanchunk.LocalRangeReplicaRetrieveOpts) error {
	if opts.PKC == nil {
		panic("replica retrieve queries require a primary key")
	}
	var rangeReplicas []*models.RangeReplica
	for _, pk := range opts.PKC {
		for _, rr := range s.RangeReplicas {
			if model.NewPK(rr.ID).Equals(pk) {
				rangeReplicas = append(rangeReplicas, rr)
			}
		}
	}
	sourceRfl := model.NewReflect(rangeReplica)
	if sourceRfl.IsChain() {
		model.NewExchange(rangeReplica, &rangeReplicas).ToSource()
	} else {
		if len(rangeReplicas) == 0 {
			return storage.Error{Type: storage.ErrorTypeItemNotFound}
		}
		model.NewExchange(rangeReplica, rangeReplicas[0]).ToSource()
	}
	return nil
}

func (s *Local) retrieveReplicaByWhereFields(flds model.WhereFields) (replicas []*models.ChannelChunkReplica) {
	rfl := model.NewReflect(&s.ChunkReplicas)
	rfl.ForEach(func(rfl *model.Reflect, i int) {
		allMatch := true
		for k, v := range flds {
			if rfl.StructFieldByName(k).Interface() != v {
				allMatch = false
			}
		}
		if allMatch {
			replicas = append(replicas, rfl.Pointer().(*models.ChannelChunkReplica))
		}
	})
	return replicas
}

func (s *Local) retrieveReplicaByPK(PKC model.PKChain) (replicas []*models.ChannelChunkReplica) {
	for _, PK := range PKC {
		for _, ccr := range s.ChunkReplicas {
			if model.NewPK(ccr.ID).Equals(PK) {
				replicas = append(replicas, ccr)
			}
		}
	}
	return replicas
}
