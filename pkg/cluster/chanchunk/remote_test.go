package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk/mock"
	clustermock "github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

func lisPort(lis net.Listener) int {
	port, pErr := strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])
	Expect(pErr).To(BeNil())
	return port
}

var _ = Describe("ServiceRemoteRPC", func() {
	var (
		pool                         *cluster.NodeRPCPool
		svc                          chanchunk.ServiceRemote
		persistOne, persistTwo       *mock.ServerRPCPersist
		grpcServerOne, grpcServerTwo *grpc.Server
		nodeOne, nodeTwo             *models.Node
		serverErr                    error
	)
	BeforeEach(func() {
		pool = clustermock.NewNodeRPCPool()
		svc = chanchunk.NewServiceRemoteRPC(pool)
		persistOne, persistTwo = &mock.ServerRPCPersist{}, &mock.ServerRPCPersist{}
		serverOne, serverTwo := chanchunk.NewServerRPC(persistOne), chanchunk.NewServerRPC(persistTwo)
		grpcServerOne, grpcServerTwo = grpc.NewServer(), grpc.NewServer()
		serverOne.BindTo(grpcServerOne)
		serverTwo.BindTo(grpcServerTwo)
	})
	JustBeforeEach(func() {
		lisOne, err := net.Listen("tcp", "localhost:0")
		Expect(err).To(BeNil())
		lisTwo, err := net.Listen("tcp", "localhost:0")

		nodeOne = &models.Node{
			ID:      1,
			RPCPort: lisPort(lisOne),
			Address: lisOne.Addr().String(),
		}
		nodeTwo = &models.Node{
			ID:      1,
			RPCPort: lisPort(lisTwo),
			Address: lisTwo.Addr().String(),
		}

		Expect(err).To(BeNil())
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
		cErr := svc.Create(ctx, []chanchunk.RemoteCreateOpts{
			{
				Node: nodeOne,
				ChunkReplica: &[]*models.ChannelChunkReplica{
					{
						ID:    idOne,
						Telem: telem.NewChunkData([]byte{1, 2, 3}),
					},
				},
			},
			{
				Node: nodeTwo,
				ChunkReplica: &[]*models.ChannelChunkReplica{{
					ID:    idTwo,
					Telem: telem.NewChunkData([]byte{3, 4, 5}),
				},
				},
			},
		})
		Expect(cErr).To(BeNil())
		Expect(persistOne.ChunkReplicas).To(HaveLen(1))
		Expect(persistOne.ChunkReplicas[0].ID).To(Equal(idOne))
		Expect(persistTwo.ChunkReplicas).To(HaveLen(1))
		Expect(persistTwo.ChunkReplicas[0].ID).To(Equal(idTwo))
	})
	It("Should delete the replicas correctly", func() {
		idOne := uuid.New()
		cErr := svc.Create(ctx, []chanchunk.RemoteCreateOpts{
			{
				Node: nodeOne,
				ChunkReplica: &[]*models.ChannelChunkReplica{
					{
						ID:    idOne,
						Telem: telem.NewChunkData([]byte{1, 2, 3}),
					},
				},
			},
		})

		Expect(cErr).To(BeNil())
		dErr := svc.Delete(ctx, []chanchunk.RemoteDeleteOpts{
			{
				Node: nodeOne,
				PKC:  model.NewPKChain([]uuid.UUID{idOne}),
			},
		})
		Expect(dErr).To(BeNil())
		Expect(persistOne.ChunkReplicas).To(HaveLen(0))
	})
	It("Should retrieve the replicas correctly", func() {
		idOne := uuid.New()
		cErr := svc.Create(ctx, []chanchunk.RemoteCreateOpts{
			{
				Node: nodeOne,
				ChunkReplica: &[]*models.ChannelChunkReplica{
					{
						ID:    idOne,
						Telem: telem.NewChunkData([]byte{1, 2, 3}),
					},
				},
			},
		})
		Expect(cErr).To(BeNil())
		var ccr []*models.ChannelChunkReplica
		rErr := svc.Retrieve(ctx, &ccr, []chanchunk.RemoteRetrieveOpts{{Node: nodeOne, PKC: model.NewPKChain([]uuid.UUID{idOne})}})
		Expect(rErr).To(BeNil())
		Expect(ccr).To(HaveLen(1))
		Expect(ccr[0].ID).To(Equal(idOne))
		Expect(ccr[0].Telem.Bytes()).To(Equal([]byte{1, 2, 3}))
	})
})
