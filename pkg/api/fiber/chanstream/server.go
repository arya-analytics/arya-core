package chanstream

import (
	"context"
	cf "github.com/arya-analytics/aryacore/pkg/api/fiber"
	"github.com/arya-analytics/aryacore/pkg/query"
	qcc "github.com/arya-analytics/aryacore/pkg/query/chanstream"
	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
	"github.com/arya-analytics/aryacore/pkg/ws"
	"github.com/gofiber/fiber/v2"
)

// |||| SERVER ||||

type Server struct {
	svc *chanstream.Service
}

func NewServer(svc *chanstream.Service) *Server {
	return &Server{svc: svc}
}

const (
	groupEndpoint    = "/stream"
	retrieveEndpoint = "/retrieve"
	createEndpoint   = "/create"
)

func (s *Server) BindTo(router fiber.Router) {
	r := router.Group(groupEndpoint)
	r.Get(retrieveEndpoint, cf.WebsocketHandler(s.retrieveStream))
	r.Get(createEndpoint, cf.WebsocketHandler(s.createStream))
}

func (s *Server) retrieveStream(c *ws.Conn) {
	defer c.Close()
	p := &fiberRetrieveProtocol{fiberBaseProtocol{conn: c, ctx: context.Background()}}
	c.SendAndWarn(qcc.RetrieveStream(s.svc, p))
}

func (s *Server) createStream(c *ws.Conn) {
	defer c.Close()
	p := &fiberCreateProtocol{fiberBaseProtocol{conn: c}}
	c.SendAndWarn(qcc.CreateStream(s.svc, p))
}

// |||| RETRIEVE PROTOCOL ||||

type fiberBaseProtocol struct {
	conn *ws.Conn
	ctx  context.Context
}

func (p *fiberBaseProtocol) Context() context.Context {
	return p.ctx
}

type fiberRetrieveProtocol struct {
	fiberBaseProtocol
}

type retrieveRequest struct {
	PKC []string
}

func (r *fiberRetrieveProtocol) Receive() (qcc.RetrieveRequest, error) {
	var req = retrieveRequest{}
	if err := r.conn.Receive(req); err != nil {
		return qcc.RetrieveRequest{}, err
	}
	pkc, err := query.ParsePKC(req.PKC)
	return qcc.RetrieveRequest{PKC: pkc}, err
}

func (r *fiberRetrieveProtocol) Send(res qcc.RetrieveResponse) error {
	return r.conn.Send(res)
}

// |||| CREATE PROTOCOL ||||

type fiberCreateProtocol struct {
	fiberBaseProtocol
}

func (c *fiberCreateProtocol) Receive() (req qcc.CreateRequest, err error) {
	return req, c.conn.Receive(req)
}

func (c *fiberCreateProtocol) Send(res qcc.CreateResponse) error {
	return c.conn.Send(res)
}
