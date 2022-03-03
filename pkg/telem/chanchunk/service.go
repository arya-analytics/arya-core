package chanchunk

import (
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/query"
)

type Service struct {
	obs    Observe
	exec   query.Execute
	rngSVC *rng.Service
}

func NewService(exec query.Execute, obs Observe, rngSVC *rng.Service) *Service {
	return &Service{exec: exec, obs: obs, rngSVC: rngSVC}
}

func (s *Service) NewStreamCreate() *QueryStreamCreate {
	return newStreamCreate(s.exec, s.obs, s.rngSVC)
}

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
//			WhereFields(query.WhereFields{"ChannelChunk.startTS": model.FieldInRange(opts.startTS, opts.EndTS)}).
//			exec(ctx); err != nil {
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
