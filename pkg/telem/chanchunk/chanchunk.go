package chanchunk

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"io"
)

type Service struct {
	cluster cluster.Cluster
	rngSVC  *rng.Service
}

func NewService(clust cluster.Cluster, rngSVC *rng.Service) *Service {
	return &Service{cluster: clust, rngSVC: rngSVC}
}

type TelemChunk struct {
	StartTS int64
	Data    *telem.Bulk
}

func (s *Service) CreateStream(ctx context.Context, cfg *models.ChannelConfig) (chan *TelemChunk, chan error) {
	stream, errChan := make(chan *TelemChunk), make(chan error)
	go func() {
		for tc := range stream {
			c := errutil.NewContextCatcher(ctx)
			alloc := s.rngSVC.NewAllocate()
			chunk := &models.ChannelChunk{
				ID:              uuid.New(),
				ChannelConfigID: cfg.ID,
				StartTS:         tc.StartTS,
				Size:            tc.Data.Size(),
			}
			repl := &models.ChannelChunkReplica{
				ID:             uuid.New(),
				ChannelChunkID: chunk.ID,
				Telem:          tc.Data,
			}
			c.Exec(alloc.Chunk(cfg.NodeID, chunk).Exec)
			c.Exec(s.cluster.NewCreate().Model(chunk).Exec)
			c.Exec(alloc.ChunkReplica(repl).Exec)
			c.Exec(s.cluster.NewCreate().Model(repl).Exec)
			if c.Error() != nil {
				errChan <- c.Error()
			}
		}
		errChan <- io.EOF
	}()
	return stream, errChan
}

type RetrieveOpts struct {
	StartTS int64
	EndTS   int64
}

func (s *Service) RetrieveStream(ctx context.Context, cfg *models.ChannelConfig, opts RetrieveOpts) (chan *TelemChunk, chan error) {
	stream, errChan := make(chan *TelemChunk), make(chan error)
	go func() {
		var replicas []*models.ChannelChunkReplica
		if err := s.cluster.NewRetrieve().
			Model(&replicas).
			Relation("ChannelChunk", "StartTS").
			WhereFields(model.WhereFields{"ChannelChunk.StartTS": model.FieldInRange(opts.StartTS, opts.EndTS)}).
			Exec(ctx); err != nil {
			errChan <- err
		}
		for _, ccr := range replicas {
			stream <- &TelemChunk{StartTS: ccr.ChannelChunk.StartTS, Data: ccr.Telem}
		}
		close(stream)
		errChan <- io.EOF
	}()
	return stream, errChan
}
