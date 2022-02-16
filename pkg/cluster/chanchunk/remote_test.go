package chanchunk_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk/mock"
	"github.com/arya-analytics/aryacore/pkg/rpc"
	rpcmock "github.com/arya-analytics/aryacore/pkg/rpc/mock"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"net"
)

var _ = Describe("ServiceRemoteRPC", func() {
	var (
		ctx        context.Context
		pool       rpc.Pool
		svc        *chanchunk.ServiceRemoteRPC
		server     *mock.Server
		grpcServer *grpc.Server
		addr       net.Addr
		serverErr  error
	)
	BeforeEach(func() {
		ctx = context.Background()
		pool = rpcmock.NewPool()
		svc = chanchunk.NewServiceRemoteRPC(pool)
		server = mock.NewServer()
		grpcServer = grpc.NewServer()
		server.BindTo(grpcServer)
	})
	JustBeforeEach(func() {
		lis, err := net.Listen("tcp", "localhost:0")
		addr = lis.Addr()
		if err != nil {
			serverErr = err
		}
		go func() {
			if err := grpcServer.Serve(lis); err != nil {
				serverErr = err
			}
		}()
		Expect(serverErr).To(BeNil())
	})
	It("Should create the channel chunks correctly", func() {
		cErr := svc.CreateReplicas(ctx, []chanchunk.RemoteReplicaCreateParams{
			{
				Addr: addr.String(),
				Model: model.NewReflect(&[]*storage.ChannelChunkReplica{
					{
						ID:    uuid.New(),
						Telem: telem.NewBulk([]byte{1, 2, 3}),
					},
				}),
			},
		})
		Expect(cErr).To(BeNil())
		Expect(server.CreatedChunks.ChainValue().Len()).To(Equal(1))
	})
})
