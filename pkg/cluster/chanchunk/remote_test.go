package chanchunk_test

import (
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
		pool                         rpc.Pool
		svc                          chanchunk.ServiceRemote
		serverOne, serverTwo         *mock.Server
		grpcServerOne, grpcServerTwo *grpc.Server
		addrOne, addrTwo             net.Addr
		serverErr                    error
	)
	BeforeEach(func() {
		pool = rpcmock.NewPool(0)
		svc = chanchunk.NewServiceRemoteRPC(pool)
		serverOne, serverTwo = mock.NewServer(), mock.NewServer()
		grpcServerOne, grpcServerTwo = grpc.NewServer(), grpc.NewServer()
		serverOne.BindTo(grpcServerOne)
		serverTwo.BindTo(grpcServerTwo)
	})
	JustBeforeEach(func() {
		lisOne, err := net.Listen("tcp", "localhost:0")
		Expect(err).To(BeNil())
		lisTwo, err := net.Listen("tcp", "localhost:0")
		Expect(err).To(BeNil())
		addrOne, addrTwo = lisOne.Addr(), lisTwo.Addr()
		go func() {
			if err := grpcServerOne.Serve(lisOne); err != nil {
				serverErr = err
			}
		}()
		Expect(serverErr).To(BeNil())
		go func() {
			if err := grpcServerTwo.Serve(lisTwo); err != nil {
				serverErr = err
			}
		}()
		Expect(serverErr).To(BeNil())
	})
	It("Should create the replicas correctly", func() {
		idOne, idTwo := uuid.New(), uuid.New()
		cErr := svc.CreateReplicas(ctx, []chanchunk.RemoteReplicaCreateParams{
			{
				Addr: addrOne.String(),
				Model: model.NewReflect(&[]*storage.ChannelChunkReplica{
					{
						ID:    idOne,
						Telem: telem.NewBulk([]byte{1, 2, 3}),
					},
				}),
			},
			{
				Addr: addrTwo.String(),
				Model: model.NewReflect(&[]*storage.ChannelChunkReplica{{
					ID:    idTwo,
					Telem: telem.NewBulk([]byte{3, 4, 5}),
				},
				}),
			},
		})
		Expect(cErr).To(BeNil())
		Expect(serverOne.CreatedChunks.ChainValue().Len()).To(Equal(1))
		Expect(serverOne.CreatedChunks.ChainValueByIndex(0).PK().String()).To(Equal(idOne.String()))
		Expect(serverTwo.CreatedChunks.ChainValue().Len()).To(Equal(1))
		Expect(serverTwo.CreatedChunks.ChainValueByIndex(0).PK().String()).To(Equal(idTwo.String()))
	})
	It("Should delete the replicas correctly", func() {
		pkC := model.NewPKChain([]uuid.UUID{uuid.New()})
		cErr := svc.DeleteReplicas(ctx, []chanchunk.RemoteReplicaDeleteParams{
			{
				Addr: addrOne.String(),
				PKC:  pkC,
			},
		})
		Expect(cErr).To(BeNil())
		Expect(serverOne.DeletedChunkPKChain.Raw()).To(Equal(pkC.Raw()))
	})
	It("Should retrieve the replicas correctly", func() {
		id := uuid.New()
		ccR := model.NewReflect(&[]*storage.ChannelChunkReplica{})
		cErr := svc.RetrieveReplicas(ctx, ccR, []chanchunk.RemoteReplicaRetrieveParams{{Addr: addrOne.String(), PKC: model.NewPKChain([]uuid.UUID{id})}})
		Expect(cErr).To(BeNil())
		Expect(ccR.ChainValue().Len()).To(Equal(1))
		Expect(ccR.ChainValueByIndex(0).PK().Raw()).To(Equal(id))
	})
})
