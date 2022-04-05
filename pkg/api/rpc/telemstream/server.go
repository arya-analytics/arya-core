package telemstream

import (
	"context"
	api "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/telemstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/query"
	qcc "github.com/arya-analytics/aryacore/pkg/query/chanstream"
	"github.com/arya-analytics/aryacore/pkg/rpc"
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
	return qcc.RetrieveStream(s.svc, &RPCRetrieveProtocol{rpcStream})
}

func (s *Server) Create(rpcStream api.TelemStreamService_CreateServer) error {
	return qcc.CreateStream(s.svc, &RPCCreateProtocol{rpcStream})
}

// |||| RETRIEVE PROTOCOL ||||

type RPCRetrieveProtocol struct {
	rpcStream api.TelemStreamService_RetrieveServer
}

func (r *RPCRetrieveProtocol) Context() context.Context {
	return r.rpcStream.Context()
}

func (r *RPCRetrieveProtocol) Receive() (qcc.RetrieveRequest, error) {
	req, err := r.rpcStream.Recv()
	if err != nil {
		return qcc.RetrieveRequest{}, err
	}
	pkc, pkcErr := query.ParsePKC(req.PKC)
	return qcc.RetrieveRequest{PKC: pkc}, pkcErr
}

func (r *RPCRetrieveProtocol) Send(resp qcc.RetrieveResponse) error {
	rpcResp := &api.RetrieveResponse{Sample: &api.TelemSample{}, Error: &api.Error{}}
	rpc.NewModelExchange(resp.Sample, rpcResp.Sample).ToDest()
	if resp.Error != nil {
		rpcResp.Error.Message = resp.Error.Error()
	}
	return r.rpcStream.Send(rpcResp)
}

// |||| CREATE PROTOCOL ||||

type RPCCreateProtocol struct {
	rpcStream api.TelemStreamService_CreateServer
}

func (r *RPCCreateProtocol) Context() context.Context {
	return r.rpcStream.Context()
}

func (r *RPCCreateProtocol) Receive() (qcc.CreateRequest, error) {
	rpcReq, err := r.rpcStream.Recv()
	if err != nil {
		return qcc.CreateRequest{}, err
	}
	req := qcc.CreateRequest{Sample: &models.ChannelSample{}}
	rpc.NewModelExchange(req.Sample, rpcReq.Sample).ToSource()
	return req, err
}

func (r *RPCCreateProtocol) Send(resp qcc.CreateResponse) error {
	rpcResp := &api.CreateResponse{Error: &api.Error{}}
	if resp.Error != nil {
		rpcResp.Error.Message = resp.Error.Error()
	}
	return r.rpcStream.Send(rpcResp)
}
