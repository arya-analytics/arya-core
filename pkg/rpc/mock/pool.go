package mock

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Pool struct {
	conns map[string]*grpc.ClientConn
}

func NewPool() *Pool {
	return &Pool{}
}

func (p *Pool) Retrieve(addr string) *grpc.ClientConn {
	conn, ok := p.conns[addr]
	if !ok {
		conn = p.newConn(addr)
	}
	return conn
}

func (p *Pool) newConn(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}
