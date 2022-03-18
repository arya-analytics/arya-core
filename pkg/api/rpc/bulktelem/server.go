package bulktelem

import (
	"context"
	bulktelemv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/bulktelem/v1"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"io"
)

type Server struct {
	bulktelemv1.UnimplementedBulkTelemServiceServer
	svc *chanchunk.Service
}

func NewServer(svc *chanchunk.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) BindTo(srv *grpc.Server) {
	bulktelemv1.RegisterBulkTelemServiceServer(srv, s)
}

func (s *Server) CreateStream(server bulktelemv1.BulkTelemService_CreateStreamServer) error {
	stream := s.svc.NewStreamCreate()
	wg := errgroup.Group{}
	wg.Go(func() error { return relayErrors(stream, server) })
	wg.Go(func() error {
		defer stream.Close()
		start := true
		for {
			req, err := server.Recv()
			if err == io.EOF {
				return nil
			}
			if err != nil {
				return err
			}
			if start {
				if sErr := startStream(server.Context(), stream, req); sErr != nil {
					return sErr
				}
				start = false
			}
			if dErr := sendData(stream, req); dErr != nil {
				return dErr
			}
		}
	})
	return wg.Wait()
}

func startStream(ctx context.Context, stream *chanchunk.StreamCreate, req *bulktelemv1.CreateStreamRequest) error {
	pk, err := model.NewPK(uuid.UUID{}).NewFromString(req.ChannelConfigId)
	if err != nil {
		return err
	}
	return stream.Start(ctx, pk.Raw().(uuid.UUID))
}

func sendData(stream *chanchunk.StreamCreate, req *bulktelemv1.CreateStreamRequest) error {
	cd := telem.NewChunkData(make([]byte, len(req.Data)))
	if _, err := cd.Write(req.Data); err != nil {
		return err
	}
	stream.Send(telem.TimeStamp(req.StartTs), cd)
	return nil
}

func relayErrors(stream *chanchunk.StreamCreate, server bulktelemv1.BulkTelemService_CreateStreamServer) error {
	for err := range stream.Errors() {
		if sErr := server.Send(&bulktelemv1.CreateStreamResponse{Error: &bulktelemv1.Error{Message: err.Error()}}); sErr != nil {
			return sErr
		}
	}
	return nil
}
