package live

import (
	"github.com/arya-analytics/aryacore/telem"
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
	"time"
)

type ReceiverConfig []int32

type Receiver interface {
	send(sect telem.Slice)
	receive() (sect telem.Slice)
	Start(cfg ReceiverConfig) ()
}

type WSReceiver struct {
	rel  *Relay
	conn *websocket.Conn
	p    chan telem.Slice
}

func (rcv WSReceiver) send(slc telem.Slice) {
	rcv.p <- slc
}

func (rcv WSReceiver) receive() (slc telem.Slice) {
	slc = telem.Slice{}
	select {
	case slc = <-rcv.p:
		return slc
	default:
		return slc
	}
}

func (rcv WSReceiver) decode(b []byte) (slc telem.Slice) {
	if err := msgpack.Unmarshal(b, slc); err != nil {
		panic(err)
	}
	return slc
}

func (rcv WSReceiver) Start(cfg ReceiverConfig) () {
	rcv.rel.addReceiver <- rcv
	defer func() {
		rcv.rel.removeReceiver <- rcv
	}()
	for {
		_, b, err := rcv.conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		rcv.send(rcv.decode(b))
	}
}

type DummyReceiver struct {
	rel *Relay
	p   chan telem.Slice
}

func (rcv DummyReceiver) send(slc telem.Slice) {
	rcv.p <- slc
}

func (rcv DummyReceiver) receive() (slc telem.Slice) {
	slc = telem.Slice{}
	select {
	case slc = <-rcv.p:
		return slc
	default:
		return slc
	}
}

func (rcv DummyReceiver) Start(cfg ReceiverConfig) () {
	rcv.rel.addReceiver <- rcv
	defer func() {
		rcv.rel.removeReceiver <- rcv
	}()
	t := time.NewTicker(100 * time.Millisecond)

	slc := telem.Slice{
		123: telem.Value{123.2, 123.2},
	}

	for {
		select {
		case <-t.C:
			rcv.send(slc)
		}

	}
}

func NewReceiver(rel *Relay) Receiver {
	p := make(chan telem.Slice)
	return &DummyReceiver{rel, p}
}
