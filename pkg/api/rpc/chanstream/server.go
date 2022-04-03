package chanstream

import (
	"context"
	api "github.com/arya-analytics/aryacore/pkg/cluster/gen/proto/go/chanstream/v1"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/query"
	qcc "github.com/arya-analytics/aryacore/pkg/query/chanstream"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
	"google.golang.org/grpc"
)

// |||| SERVER ||||

type Server struct {
	api.UnimplementedChannelStreamServiceServer
	svc *chanstream.Service
}

func NewServer(svc *chanstream.Service) *Server {
	return &Server{svc: svc}
}

func (s *Server) BindTo(srv *grpc.Server) {
	api.RegisterChannelStreamServiceServer(srv, s)
}

func (s *Server) RetrieveStream(rpcStream api.ChannelStreamService_RetrieveServer) error {
	return qcc.RetrieveStream(s.svc, &RPCRetrieveProtocol{rpcStream})
}

func (s *Server) CreateStream(rpcStream api.ChannelStreamService_CreateServer) error {
	return qcc.CreateStream(s.svc, &RPCCreateProtocol{rpcStream})
}

// |||| RETRIEVE PROTOCOL ||||

type RPCRetrieveProtocol struct {
	rpcStream api.ChannelStreamService_RetrieveServer
}

func (r *RPCRetrieveProtocol) Context() context.Context {
	return r.rpcStream.Context()
}

func (r *RPCRetrieveProtocol) Receive() (qcc.RetrieveRequest, error) {
	req, err := r.rpcStream.Recv()
	pkc, pkcErr := query.ParsePKC(req.PKC)
	if pkcErr != nil {
		return qcc.RetrieveRequest{}, pkcErr
	}
	return qcc.RetrieveRequest{PKC: pkc}, err
}

func (r *RPCRetrieveProtocol) Send(response qcc.RetrieveResponse) error {
	resp := &api.RetrieveResponse{}
	rpc.NewModelExchange(response, resp).ToDest()
	return r.rpcStream.Send(resp)
}

// |||| CREATE PROTOCOL ||||

type RPCCreateProtocol struct {
	rpcStream api.ChannelStreamService_CreateServer
}

func (r *RPCCreateProtocol) Context() context.Context {
	return r.rpcStream.Context()
}

func (r *RPCCreateProtocol) Receive() (qcc.CreateRequest, error) {
	rpcReq, err := r.rpcStream.Recv()
	req := qcc.CreateRequest{Sample: &models.ChannelSample{}}
	rpc.NewModelExchange(req.Sample, rpcReq.Sample).ToDest()
	return req, err
}

func (r *RPCCreateProtocol) Send(resp qcc.CreateResponse) error {
	rpcResp := &api.CreateResponse{}
	rpc.NewModelExchange(resp, rpcResp).ToDest()
	return r.rpcStream.Send(rpcResp)
}
