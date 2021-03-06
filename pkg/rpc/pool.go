package rpc

import (
	"google.golang.org/grpc"
	"sync"
)

type Pool struct {
	dialOpts []grpc.DialOption
	mu       sync.RWMutex
	conns    map[string]*grpc.ClientConn
}

func NewPool(dialOpts ...grpc.DialOption) *Pool {
	return &Pool{dialOpts: dialOpts, conns: map[string]*grpc.ClientConn{}}
}

func (p *Pool) Retrieve(target string) (conn *grpc.ClientConn, err error) {
	var ok bool
	conn, ok = p.conns[target]
	if !ok {
		conn, err = p.newConn(target)
		p.addConn(target, conn)
	}
	return conn, err
}

func (p *Pool) newConn(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, p.dialOpts...)
}

func (p *Pool) addConn(target string, conn *grpc.ClientConn) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.conns[target] = conn
}
