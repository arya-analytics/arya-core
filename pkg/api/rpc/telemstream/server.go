package telemstream

import (
	"context"
	api "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/telemstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/query"
	qcs "github.com/arya-analytics/aryacore/pkg/query/chanstream"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	errorv1 "github.com/arya-analytics/aryacore/pkg/rpc/gen/proto/go/error/v1"
	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
	"google.golang.org/grpc"
)

// |||| SERVER ||||

type Server struct {
	api.UnimplementedTelemStreamServiceServer
	svc *chanstream.Service
}

func NewServer(svc *chanstream.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) BindTo(srv *grpc.Server) {
	api.RegisterTelemStreamServiceServer(srv, s)
}

func (s *Server) Retrieve(rpcStream api.TelemStreamService_RetrieveServer) error {
	return qcs.RetrieveStream(s.svc, &retrieveProtocol{rpcStream})
}

func (s *Server) Create(rpcStream api.TelemStreamService_CreateServer) error {
	return qcs.CreateStream(s.svc, &createProtocol{rpcStream})
}

// |||| RETRIEVE PROTOCOL ||||

type retrieveProtocol struct {
	conn api.TelemStreamService_RetrieveServer
}

func (r *retrieveProtocol) Context() context.Context {
	return r.conn.Context()
}

func (r *retrieveProtocol) Receive() (qcs.RetrieveRequest, error) {
	req, err := r.conn.Recv()
	if err != nil {
		return qcs.RetrieveRequest{}, err
	}
	pkc, pkcErr := query.ParsePKC(req.Pkc)
	return qcs.RetrieveRequest{PKC: pkc}, pkcErr
}

func (r *retrieveProtocol) Send(resp qcs.RetrieveResponse) error {
	rpcResp := &api.RetrieveResponse{TelemSample: &api.TelemSample{}, Error: &errorv1.Error{}}
	rpc.NewModelExchange(resp.Sample, rpcResp.TelemSample).ToDest()
	if resp.Error != nil {
		rpcResp.Error.Message = resp.Error.Error()
	}
	return r.conn.Send(rpcResp)
}

// |||| CREATE PROTOCOL ||||

type createProtocol struct {
	conn api.TelemStreamService_CreateServer
}

func (r *createProtocol) Context() context.Context {
	return r.conn.Context()
}

func (r *createProtocol) Receive() (qcs.CreateRequest, error) {
	rpcReq, err := r.conn.Recv()
	if err != nil {
		return qcs.CreateRequest{}, err
	}
	req := qcs.CreateRequest{Sample: &models.ChannelSample{}}
	rpc.NewModelExchange(req.Sample, rpcReq.TelemSample).ToSource()
	return req, err
}

func (r *createProtocol) Send(resp qcs.CreateResponse) error {
	rpcResp := &api.CreateResponse{Error: &errorv1.Error{}}
	if resp.Error != nil {
		rpcResp.Error.Message = resp.Error.Error()
	}
	return r.conn.Send(rpcResp)
}
