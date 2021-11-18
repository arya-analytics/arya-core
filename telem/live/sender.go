package live

import (
	"github.com/arya-analytics/aryacore/telem"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/vmihailenco/msgpack/v5"
)

type SenderConfig struct {
	ID       uuid.UUID
	ChanCfgs []int32
}

type Sender interface {
	send(slc telem.Slice)
	receive() (slc telem.Slice)
	start()
	id() uuid.UUID
}

type WSSender struct {
	ID   uuid.UUID
	rel  *Relay
	conn *websocket.Conn
	p    chan telem.Slice
}

func (s WSSender) send(slc telem.Slice) {
	s.p <- slc
}

func (s WSSender) receive() (slc telem.Slice) {
	return <-s.p
}

func (s WSSender) id() uuid.UUID {
	return s.ID
}

func (s WSSender) encode(slc telem.Slice) []byte {
	b, err := msgpack.Marshal(slc)
	if err != nil {
		panic(err)
	}
	return b
}

func (s WSSender) decode(b []byte) (cfg SenderConfig) {
	cfg.ID = s.id()
	if err := msgpack.Unmarshal(b, &cfg); err != nil {
		panic(err)
	}
	return cfg
}

func (s WSSender) listen() {
	for {
		_, msg, err := s.conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		cfg := s.decode(msg)
		s.rel.updateSenderConfig <- cfg
	}
}

func (s WSSender) start() () {
	go s.listen()
	s.rel.addSender <- s
	defer func() {
		s.rel.removeSender <- s
	}()
	for {
		slc := s.receive()
		data, err := msgpack.Marshal(slc)
		if err != nil {
			panic(err)
		}
		if err := s.conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
			panic(err)
		}
	}
}

func NewSender(rel *Relay, conn *websocket.Conn) Sender {
	p := make(chan telem.Slice)
	id := uuid.New()
	return &WSSender{id, rel, conn, p}
}
