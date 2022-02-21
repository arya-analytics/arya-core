package rpc_test

import (
	"github.com/arya-analytics/aryacore/pkg/rpc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

var _ = Describe("Pool", func() {
	var (
		lis    net.Listener
		server *grpc.Server
		pool   *rpc.Pool
	)
	BeforeEach(func() {
		var lisErr error
		lis, lisErr = net.Listen("tcp", "localhost:0")
		Expect(lisErr).To(BeNil())
		server = grpc.NewServer()
		pool = rpc.NewPool(grpc.WithTransportCredentials(insecure.NewCredentials()))
	})
	JustBeforeEach(func() {
		var err error
		go func() {
			err = server.Serve(lis)
		}()
		Expect(err).To(BeNil())
	})
	It("Should return the correct connection", func() {
		conn, err := pool.Retrieve(lis.Addr().String())
		Expect(err).To(BeNil())
		Expect(conn.GetState()).To(Equal(connectivity.Idle))
	})
})
