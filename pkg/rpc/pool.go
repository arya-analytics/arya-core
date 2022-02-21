package rpc

import (
	"google.golang.org/grpc"
)

type Pool struct {
	dialOpts []grpc.DialOption
	conns    map[string]*grpc.ClientConn
}

func NewPool(DialOpts ...grpc.DialOption) *Pool {
	return &Pool{dialOpts: DialOpts, conns: map[string]*grpc.ClientConn{}}
}

func (p *Pool) Retrieve(target string) (*grpc.ClientConn, error) {
	conn, ok := p.conns[target]
	if ok {
		return conn, nil
	}
	if !ok {
		var err error
		conn, err = p.newConn(target)
		p.conns[target] = conn
		return conn, err
	}
	return conn, nil
}

func (p *Pool) newConn(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, p.dialOpts...)
}
