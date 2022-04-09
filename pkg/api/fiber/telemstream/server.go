package telemstream

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
	p := &retrieveProtocol{baseProtocol{conn: c, ctx: context.Background()}}
	c.WriteAndWarn(qcc.RetrieveStream(s.svc, p))
}

func (s *Server) createStream(c *ws.Conn) {
	defer c.Close()
	p := &createProtocol{baseProtocol{conn: c}}
	c.WriteAndWarn(qcc.CreateStream(s.svc, p))
}

// |||| RETRIEVE PROTOCOL ||||

type baseProtocol struct {
	conn *ws.Conn
	ctx  context.Context
}

func (p *baseProtocol) Context() context.Context {
	return p.ctx
}

type retrieveProtocol struct {
	baseProtocol
}

type retrieveRequest struct {
	PKC []string
}

func (r *retrieveProtocol) Receive() (qcc.RetrieveRequest, error) {
	var req = retrieveRequest{}
	if err := r.conn.ReadInto(req); err != nil {
		return qcc.RetrieveRequest{}, err
	}
	pkc, err := query.ParsePKC(req.PKC)
	return qcc.RetrieveRequest{PKC: pkc}, err
}

func (r *retrieveProtocol) Send(res qcc.RetrieveResponse) error {
	return r.conn.Write(res)
}

// |||| CREATE PROTOCOL ||||

type createProtocol struct {
	baseProtocol
}

func (c *createProtocol) Receive() (req qcc.CreateRequest, err error) {
	return req, c.conn.ReadInto(req)
}

func (c *createProtocol) Send(res qcc.CreateResponse) error {
	return c.conn.Write(res)
}
