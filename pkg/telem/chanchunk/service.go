package chanchunk

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
)

type Service struct {
	cluster cluster.Cluster
	rngSVC  *rng.Service
}

func NewService(clust cluster.Cluster, rngSVC *rng.Service) *Service {
	return &Service{cluster: clust, rngSVC: rngSVC}
}

func (s *Service) NewStreamCreate() *QueryStreamCreate {
	return newStreamCreate(s.cluster, s.rngSVC)
}

//
//func (s *Service) CreateStream(ctx context.Context, cfg *models.ChannelConfig) (chan *TelemChunkWrapper, chan error) {
//	stream, errChan := make(chan *TelemChunkWrapper), make(chan error)
//	go func() {
//		for tc := range stream {
//			c := errutil.NewContextCatcher(ctx)
//			alloc := s.rngSVC.NewAllocate()
//			chunk := &models.ChannelChunk{
//				ID:              uuid.New(),
//				ChannelConfigID: cfg.ID,
//				startTS:         tc.startTS,
//				Size:            tc.Telem.Size(),
//			}
//			repl := &models.ChannelChunkReplica{
//				ID:             uuid.New(),
//				ChannelChunkID: chunk.ID,
//				Telem:          tc.Telem,
//			}
//			c.Exec(alloc.ChunkData(cfg.NodeID, chunk).Exec)
//			c.Exec(s.cluster.NewCreate().Model(chunk).Exec)
//			c.Exec(alloc.ChunkReplica(repl).Exec)
//			c.Exec(s.cluster.NewCreate().Model(repl).Exec)
//			if c.Error() != nil {
//				errChan <- c.Error()
//			}
//		}
//		errChan <- io.EOF
//	}()
//	return stream, errChan
//}
//
//type RetrieveOpts struct {
//	startTS int64
//	EndTS   int64
//}
//
//func (s *Service) RetrieveStream(ctx context.Context, cfg *models.ChannelConfig, opts RetrieveOpts) (chan *TelemChunkWrapper, chan error) {
//	stream, errChan := make(chan *TelemChunkWrapper), make(chan error)
//	go func() {
//		var replicas []*models.ChannelChunkReplica
//		if err := s.cluster.NewRetrieve().
//			Model(&replicas).
//			Relation("ChannelChunk", "startTS").
//			WhereFields(model.WhereFields{"ChannelChunk.startTS": model.FieldInRange(opts.startTS, opts.EndTS)}).
//			Exec(ctx); err != nil {
//			errChan <- err
//		}
//		for _, ccr := range replicas {
//			stream <- &TelemChunkWrapper{startTS: ccr.ChannelChunk.startTS, Telem: ccr.Telem}
//		}
//		close(stream)
//		errChan <- io.EOF
//	}()
//	return stream, errChan
//}
