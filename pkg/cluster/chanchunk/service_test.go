package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk/mock"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	rpcmock "github.com/arya-analytics/aryacore/pkg/rpc/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

var _ = Describe("Service", func() {
	var (
		remoteSvc  chanchunk.ServiceRemote
		localSvc   chanchunk.ServiceLocal
		svc        *chanchunk.Service
		pool       rpc.Pool
		server     *mock.Server
		grpcServer *grpc.Server
		lis        net.Listener
	)
	BeforeEach(func() {
		var lisErr error
		lis, lisErr = net.Listen("tcp", "localhost:0")
		Expect(lisErr).To(BeNil())
		port, pErr := strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])
		Expect(pErr).To(BeNil())
		pool = rpcmock.NewPool(port)
		remoteSvc = chanchunk.NewServiceRemoteRPC(pool)
		server = mock.NewServer()
		grpcServer = grpc.NewServer()
		server.BindTo(grpcServer)
		localSvc = chanchunk.NewServiceLocalStorage(store)
		svc = chanchunk.NewService(localSvc, remoteSvc)
	})
	JustBeforeEach(func() {
		var serverErr error
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				serverErr = err
			}
		}()
		Expect(serverErr).To(BeNil())
	})
	It("Should do it all correctly", func() {
		Expect(svc).ToNot(BeNil())
	})

})
