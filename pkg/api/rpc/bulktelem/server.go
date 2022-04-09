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
		&streamRetrieveProtocol{conn: server},
		qcc.StreamRetrieveRequest{
			ChannelConfigID: parsePK(req.ChannelConfigId),
			TimeRange:       telem.NewTimeRange(telem.TimeStamp(req.StartTs), telem.TimeStamp(req.EndTs)),
		},
	)
}

func (s *Server) CreateStream(server bulktelemv1.BulkTelemService_CreateStreamServer) error {
	err := qcc.CreateStream(s.svc, &streamCreateProtocol{conn: server})
	return err
}

func parsePK(pkStr string) uuid.UUID {
	pk, _ := model.NewPK(uuid.UUID{}).NewFromString(pkStr)
	return pk.Raw().(uuid.UUID)
}

// ||||| RETRIEVE PROTOCOL |||||

type streamRetrieveProtocol struct {
	conn bulktelemv1.BulkTelemService_RetrieveStreamServer
}

func (r *streamRetrieveProtocol) Context() context.Context {
	return r.conn.Context()
}

func (r *streamRetrieveProtocol) Send(resp qcc.StreamRetrieveResponse) error {
	return r.conn.Send(&bulktelemv1.RetrieveStreamResponse{
		StartTs:  int64(resp.StartTS),
		DataType: int64(resp.DataType),
		DataRate: float32(resp.DataRate),
		Data:     resp.Data.Bytes(),
	})
}

// ||||| CREATE PROTOCOL |||||

type streamCreateProtocol struct {
	conn bulktelemv1.BulkTelemService_CreateStreamServer
}

func (c *streamCreateProtocol) Context() context.Context {
	return c.conn.Context()
}

func (c *streamCreateProtocol) Send(resp qcc.StreamCreateResponse) error {
	return c.conn.Send(&bulktelemv1.CreateStreamResponse{
		Error: &bulktelemv1.Error{Message: resp.Error.Error()},
	})
}

func (c *streamCreateProtocol) Receive() (qcc.StreamCreateRequest, error) {
	req, err := c.conn.Recv()
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
