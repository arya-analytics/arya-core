package ws

import (
	"github.com/arya-analytics/aryacore/pkg/query"
	"github.com/arya-analytics/aryacore/pkg/util/errutil"
	"github.com/gofiber/websocket/v2"
	"github.com/vmihailenco/msgpack/v5"
)

type Conn struct {
	*websocket.Conn
}

func (c *Conn) Send(msg interface{}) error {
	b, err := marshal(msg)
	if err != nil {
		return err
	}
	err = c.WriteMessage(websocket.BinaryMessage, b)
	if Closed(err) {
		return query.StreamCloseError()
	}
	return err
}

func (c *Conn) Receive(msg interface{}) error {
	_, b, err := c.ReadMessage()
	if Closed(err) {
		return query.StreamCloseError()
	}
	if err != nil {
		return err
	}
	return unMarshal(b, msg)
}

func (c *Conn) SendAndWarn(err error) {
	if err != nil {
		errutil.Warn(c.SendError(err))
		errutil.Warn(err)
	}
}

type ErrorMsg struct {
	Error string
}

func (c *Conn) SendError(err error) error {
	return c.Send(ErrorMsg{Error: err.Error()})
}

func (c *Conn) Close() {
	errutil.Warn(c.SendError(query.StreamCloseError()))
	errutil.Warn(c.Conn.Close())
}

func Closed(err error) bool {
	return websocket.IsCloseError(err, websocket.CloseNormalClosure)
}

func marshal(msg interface{}) ([]byte, error) {
	return msgpack.Marshal(msg)
}

func unMarshal(b []byte, msg interface{}) error {
	return msgpack.Unmarshal(b, msg)
}
