package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/google/uuid"
)

type Service struct {
	cluster cluster.Cluster
	rngSVC  *rng.Service
}

func NewService(clust cluster.Cluster, rngSVC *rng.Service) *Service {
	return &Service{cluster: clust, rngSVC: rngSVC}
}

func (s *Service) CreateStream(ctx context.Context, cfg *models.ChannelConfig) (chan *models.ChannelChunkReplica, chan error) {
	stream, errChan := make(chan *models.ChannelChunkReplica), make(chan error)
	go func() {
		for repl := range stream {
			c := errutil.NewContextCatcher(ctx)
			alloc := s.rngSVC.NewAllocate()
			chunk := &models.ChannelChunk{ID: uuid.New(), ChannelConfigID: cfg.ID}
			repl.ChannelChunkID = chunk.ID
			c.Exec(alloc.Chunk(cfg.NodeID, chunk).Exec)
			c.Exec(s.cluster.NewCreate().Model(chunk).Exec)
			c.Exec(alloc.ChunkReplica(repl).Exec)
			c.Exec(s.cluster.NewCreate().Model(repl).Exec)
			if c.Error() != nil {
				errChan <- c.Error()
			}
		}
	}()
	return stream, errChan
}

type RetrieveOpts struct {
	startTS int64
	endTS   int64
}

//func (s *Service) RetrieveStream(ctx context.Context, cfg *models.ChannelConfig, opts RetrieveOpts) (chan *models.ChannelChunkReplica, chan error) {
//}
