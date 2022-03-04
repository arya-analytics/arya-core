package bulktelem

import (
	"context"
	bulktelemv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/bulktelem/v1"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"io"
	"sync"
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
	start := true
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		relayErrors(stream, server)
		wg.Done()
	}()
	for {
		req, err := server.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			stream.Errors() <- err
		}

		if start {
			if sErr := startStream(server.Context(), stream, req); sErr != nil {
				return sErr
			}
			start = false
		}

		sendData(stream, req)
	}
	stream.Close()
	wg.Wait()
	return nil
}

func startStream(ctx context.Context, stream *chanchunk.QueryStreamCreate, req *bulktelemv1.CreateStreamRequest) error {
	pk, err := model.NewPK(uuid.UUID{}).NewFromString(req.ChannelConfigId)
	if err != nil {
		return err
	}
	return stream.Start(ctx, pk.Raw().(uuid.UUID))
}

func sendData(stream *chanchunk.QueryStreamCreate, req *bulktelemv1.CreateStreamRequest) {
	cd := telem.NewChunkData(make([]byte, len(req.Data)))
	cd.Write(req.Data)
	stream.Send(telem.TimeStamp(req.StartTs), cd)
}

func relayErrors(stream *chanchunk.QueryStreamCreate, server bulktelemv1.BulkTelemService_CreateStreamServer) {
	for err := range stream.Errors() {
		server.Send(&bulktelemv1.CreateStreamResponse{Error: &bulktelemv1.Error{Message: err.Error()}})
	}
}
