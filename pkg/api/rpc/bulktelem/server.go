package bulktelem

import (
	"errors"
	bulktelemv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/bulktelem/v1"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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
	ctx := server.Context()
	stream := s.svc.NewStreamCreate()
	baseReq, err := server.Recv()
	if err != nil {
		return err
	}
	pk, err := model.NewPK(uuid.UUID{}).NewFromString(baseReq.ChannelConfigID)
	if err != nil {
		return errors.New("Invalid channel config pk specified")
	}
	if err := stream.Start(ctx, pk.Raw().(uuid.UUID)); err != nil {
		return err
	}

	data := telem.NewChunkData(make([]byte, len(baseReq.Data)))
	if _, err := data.Write(baseReq.Data); err != nil {
		return err
	}

	stream.Send(telem.TimeStamp(baseReq.StartTS), data)

	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		for {
			req, err := server.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
				stream.Errors() <- err
			}
			data := telem.NewChunkData(make([]byte, len(req.Data)))
			if _, err := data.Write(req.Data); err != nil {
				log.Fatalln(err)
				stream.Errors() <- err
			}
			stream.Send(telem.TimeStamp(req.StartTS), data)
		}
		stream.Close()
		wg.Done()
	}()
	wg.Wait()
	for err := range stream.Errors() {
		log.Fatalln(err)
		server.Send(&bulktelemv1.BulkTelemServiceCreateStreamResponse{
			Error: &bulktelemv1.BulkTelemServiceError{Message: err.Error()},
		})
	}
	wg.Wait()
	return nil
}
