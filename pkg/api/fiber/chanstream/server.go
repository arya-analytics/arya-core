package chanstream

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/query"
	qcc "github.com/arya-analytics/aryacore/pkg/query/chanstream"
	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	log "github.com/sirupsen/logrus"
	"github.com/vmihailenco/msgpack/v5"
	"io"
)

// |||| SERVER ||||

type Server struct {
	svc *chanstream.Service
}

func NewServer(svc *chanstream.Service) *Server {
	return &Server{svc: svc}
}

const (
	groupEndpoint    = "/telemstream"
	retrieveEndpoint = "/retrieve"
	createEndpoint   = "/create"
)

func (s *Server) BindTo(router fiber.Router) {
	r := router.Group(groupEndpoint)
	r.Get(retrieveEndpoint, websocket.New(s.retrieveStream))
	r.Get(createEndpoint, websocket.New(s.createStream))
}

func (s *Server) retrieveStream(c *websocket.Conn) {
	defer func() {
		if err := c.Close(); err != nil {
			log.Warn(err)
		}
	}()
	p := &FiberRetrieveProtocol{wsStream: c, ctx: context.Background()}
	if err := qcc.RetrieveStream(s.svc, p); err != nil {
		_ = p.Send(qcc.RetrieveResponse{Error: err})
		log.Warn(err)
	}
}

func (s *Server) createStream(c *websocket.Conn) {
	defer func() {
		if err := c.Close(); err != nil {
			log.Warn(err)
		}
	}()
	p := &FiberCreateProtocol{wsStream: c}
	if err := qcc.CreateStream(s.svc, p); err != nil {
		_ = p.Send(qcc.CreateResponse{Error: err})
		log.Warn(err)
	}
}

// |||| RETRIEVE PROTOCOL ||||

type FiberRetrieveProtocol struct {
	wsStream *websocket.Conn
	ctx      context.Context
}

type RetrieveRequest struct {
	PKC []string
}

func (r *FiberRetrieveProtocol) Context() context.Context {
	return r.ctx
}

func (r *FiberRetrieveProtocol) Receive() (qcc.RetrieveRequest, error) {
	var (
		msg   []byte
		wsReq RetrieveRequest
		req   qcc.RetrieveRequest
	)
	c := errutil.NewCatchSimple()
	c.Exec(func() (err error) { _, msg, err = r.wsStream.ReadMessage(); return err })
	c.Exec(func() error { return msgpack.Unmarshal(msg, &wsReq) })
	c.Exec(func() (err error) { req.PKC, err = query.ParsePKC(wsReq.PKC); return err })
	if c.Error() != nil && websocket.IsCloseError(c.Error(), websocket.CloseNormalClosure) {
		return req, io.EOF
	}
	return req, c.Error()
}

func (r *FiberRetrieveProtocol) Send(res qcc.RetrieveResponse) error {
	msg, err := msgpack.Marshal(res)
	if err != nil {
		return err
	}
	return r.wsStream.WriteMessage(websocket.BinaryMessage, msg)
}

// |||| CREATE PROTOCOL ||||

type FiberCreateProtocol struct {
	wsStream *websocket.Conn
	ctx      context.Context
}

func (c *FiberCreateProtocol) Context() context.Context {
	return c.ctx
}

func (c *FiberCreateProtocol) Receive() (qcc.CreateRequest, error) {
	var (
		msg []byte
		req qcc.CreateRequest
	)
	ca := errutil.NewCatchSimple()
	ca.Exec(func() (err error) { _, msg, err = c.wsStream.ReadMessage(); return err })
	ca.Exec(func() error { return msgpack.Unmarshal(msg, &req) })
	if ca.Error() != nil && websocket.IsCloseError(ca.Error(), websocket.CloseNormalClosure) {
		return req, io.EOF
	}
	return req, ca.Error()
}

func (c *FiberCreateProtocol) Send(res qcc.CreateResponse) error {
	msg, err := msgpack.Marshal(res)
	if err != nil {
		return err
	}
	return c.wsStream.WriteMessage(websocket.BinaryMessage, msg)
}
