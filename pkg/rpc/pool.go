package rpc

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Pool interface {
	Retrieve(addr string) *grpc.ClientConn
}

type PoolImpl struct {
	conns map[string]*grpc.ClientConn
}

func (p *PoolImpl) Retrieve(addr string) *grpc.ClientConn {
	conn, ok := p.conns[addr]
	if !ok {
		conn = p.newConn(addr)
	}
	return conn
}

func (p *PoolImpl) newConn(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr)
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}
