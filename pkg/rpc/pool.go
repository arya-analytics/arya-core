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

func (p *Pool) Retrieve(addr string) (*grpc.ClientConn, error) {
	conn, ok := p.conns[addr]
	if !ok {
		var err error
		conn, err = p.newConn(addr)
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

func (p *Pool) newConn(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr, p.dialOpts...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
