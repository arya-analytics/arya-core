package bulktelem

import (
	"context"
	bulktelemv1 "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/bulktelem/v1"
	qcc "github.com/arya-analytics/aryacore/pkg/query/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	"google.golang.org/grpc"
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

func (s *Server) RetrieveStream(req *bulktelemv1.RetrieveStreamRequest, server bulktelemv1.BulkTelemService_RetrieveStreamServer) error {
	return qcc.RetrieveStream(
		s.svc,
		&RPCStreamRetrieveProtocol{rpcStream: server},
		qcc.StreamRetrieveRequest{ChannelConfigID: parsePK(req.ChannelConfigId)},
	)
}

func (s *Server) CreateStream(server bulktelemv1.BulkTelemService_CreateStreamServer) error {
	return qcc.CreateStream(s.svc, &RPCStreamCreateProtocol{rpcStream: server})
}

func parsePK(pkStr string) uuid.UUID {
	pk, _ := model.NewPK(uuid.UUID{}).NewFromString(pkStr)
	return pk.Raw().(uuid.UUID)
}

// ||||| RETRIEVE PROTOCOL |||||

type RPCStreamRetrieveProtocol struct {
	rpcStream bulktelemv1.BulkTelemService_RetrieveStreamServer
}

func (r *RPCStreamRetrieveProtocol) Context() context.Context {
	return r.rpcStream.Context()
}

func (r *RPCStreamRetrieveProtocol) Send(resp qcc.StreamRetrieveResponse) error {
	return r.rpcStream.Send(&bulktelemv1.RetrieveStreamResponse{
		StartTs:  int64(resp.StartTS),
		DataType: int64(resp.DataType),
		DataRate: float32(resp.DataRate),
		Data:     resp.Data.Bytes(),
	})
}

// ||||| CREATE PROTOCOL |||||

type RPCStreamCreateProtocol struct {
	rpcStream bulktelemv1.BulkTelemService_CreateStreamServer
}

func (c *RPCStreamCreateProtocol) Context() context.Context {
	return c.rpcStream.Context()
}

func (c *RPCStreamCreateProtocol) Send(resp qcc.StreamCreateResponse) error {
	return c.rpcStream.Send(&bulktelemv1.CreateStreamResponse{
		Error: &bulktelemv1.Error{Message: resp.Error.Error()},
	})
}

func (c *RPCStreamCreateProtocol) Receive() (qcc.StreamCreateRequest, error) {
	req, err := c.rpcStream.Recv()
	if err != nil {
		return qcc.StreamCreateRequest{}, err
	}
	cd := telem.NewChunkData(make([]byte, len(req.Data)))
	if _, err := cd.Write(req.Data); err != nil {
		return qcc.StreamCreateRequest{}, err
	}
	return qcc.StreamCreateRequest{
		ConfigPK:  parsePK(req.ChannelConfigId),
		ChunkData: cd,
		StartTS:   telem.TimeStamp(req.StartTs),
	}, nil
}
