package mock

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Pool struct {
	port  int
	conns map[string]*grpc.ClientConn
}

func NewPool(port int) *Pool {
	return &Pool{port: port}
}

func (p *Pool) buildAddr(addr string) string {
	if p.port != 0 {
		return fmt.Sprintf("%s:%v", addr, p.port)
	}
	return addr
}

func (p *Pool) Retrieve(addr string) *grpc.ClientConn {
	conn, ok := p.conns[p.buildAddr(addr)]
	if !ok {
		conn = p.newConn(p.buildAddr(addr))
	}
	return conn
}

func (p *Pool) newConn(addr string) *grpc.ClientConn {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}
