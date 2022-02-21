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
		serverOne, serverTwo         *mock.Server
		grpcServerOne, grpcServerTwo *grpc.Server
		nodeOne, nodeTwo             *models.Node
		serverErr                    error
	)
	BeforeEach(func() {
		pool = clustermock.NewNodeRPCPool()
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

		nodeOne = &models.Node{
			ID:       1,
			GRPCPort: lisPort(lisOne),
			Address:  lisOne.Addr().String(),
		}
		nodeTwo = &models.Node{
			ID:       1,
			GRPCPort: lisPort(lisTwo),
			Address:  lisTwo.Addr().String(),
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
		cErr := svc.CreateReplica(ctx, []chanchunk.RemoterReplicaCreateOpts{
			{
				Node: nodeOne,
				ChunkReplica: &[]*models.ChannelChunkReplica{
					{
						ID:    idOne,
						Telem: telem.NewBulk([]byte{1, 2, 3}),
					},
				},
			},
			{
				Node: nodeTwo,
				ChunkReplica: &[]*models.ChannelChunkReplica{{
					ID:    idTwo,
					Telem: telem.NewBulk([]byte{3, 4, 5}),
				},
				},
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
		cErr := svc.DeleteReplica(ctx, []chanchunk.RemoteReplicaDeleteOpts{
			{
				Node: nodeOne,
				PKC:  pkC,
			},
		})
		Expect(cErr).To(BeNil())
		Expect(serverOne.DeletedChunkPKChain.Raw()).To(Equal(pkC.Raw()))
	})
	It("Should retrieve the replicas correctly", func() {
		id := uuid.New()
		var ccr []*models.ChannelChunkReplica
		cErr := svc.RetrieveReplica(ctx, &ccr, []chanchunk.RemoteReplicaRetrieveOpts{{Node: nodeTwo, PKC: model.NewPKChain([]uuid.UUID{id})}})
		Expect(cErr).To(BeNil())
		Expect(ccr).To(HaveLen(1))
		Expect(ccr[0].ID).To(Equal(id))
	})
})
