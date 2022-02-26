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

var _ = Describe("Service", func() {
	var (
		clust         cluster.Cluster
		remoteSvc     chanchunk.ServiceRemote
		localSvc      chanchunk.Local
		svc           *chanchunk.Service
		pool          *cluster.NodeRPCPool
		server        *mock.Server
		grpcServer    *grpc.Server
		lis           net.Listener
		node          = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
		rangeX        = &models.Range{
			ID: uuid.New(),
		}
		rangeReplica = &models.RangeReplica{
			ID:      uuid.New(),
			RangeID: rangeX.ID,
			NodeID:  node.ID,
		}
		rangeReplicaTwo = &models.RangeReplica{
			ID:      uuid.New(),
			RangeID: rangeX.ID,
			NodeID:  node.ID,
		}
		channelChunk = &models.ChannelChunk{
			ID:              uuid.New(),
			RangeID:         rangeX.ID,
			ChannelConfigID: channelConfig.ID,
		}
		channelChunkReplica = &models.ChannelChunkReplica{
			RangeReplicaID: rangeReplica.ID,
			ChannelChunkID: channelChunk.ID,
		}
		items = []interface{}{
			node,
			channelConfig,
			rangeX,
			rangeReplica,
			rangeReplicaTwo,
			channelChunk,
		}
	)
	BeforeEach(func() {
		store.AddQueryHook(mock.HostInterceptQueryHook(2))
		clust = cluster.New()
		channelChunkReplica.Telem = telem.NewBulk([]byte("randomdata"))
		var lisErr error
		lis, lisErr = net.Listen("tcp", "localhost:0")
		Expect(lisErr).To(BeNil())
		port, pErr := strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])
		Expect(pErr).To(BeNil())
		node.RPCPort = port
		pool = clustermock.NewNodeRPCPool()
		remoteSvc = chanchunk.NewServiceRemoteRPC(pool)
		server = mock.NewServer()
		grpcServer = grpc.NewServer()
		server.BindTo(grpcServer)
		localSvc = chanchunk.NewServiceLocalStorage(store)
		svc = chanchunk.NewService(localSvc, remoteSvc)
		clust.BindService(svc)
	})
	BeforeEach(func() {
		localSvc = chanchunk.NewServiceLocalStorage(store)

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
	JustBeforeEach(func() {
		for _, m := range items {
			err := store.NewCreate().Model(m).Exec(ctx)
			Expect(err).To(BeNil())
		}
	})
	JustAfterEach(func() {
		for _, m := range items {
			err := store.NewDelete().Model(m).WherePK(model.NewReflect(m).PK().Raw()).Exec(ctx)
			Expect(err).To(BeNil())
		}
	})
	Describe("Channel Chunk Replica", func() {
		It("Should Create the chunk replica correctly", func() {
			id := uuid.New()
			err := clust.NewCreate().Model(&models.ChannelChunkReplica{
				ID:             id,
				RangeReplicaID: rangeReplica.ID,
				ChannelChunkID: channelChunk.ID,
				Telem:          telem.NewBulk([]byte{}),
			}).Exec(ctx)
			Expect(err).To(BeNil())
			Expect(server.CreatedChunks.ChainValue().Len()).To(Equal(1))
			Expect(server.CreatedChunks.ChainValueByIndex(0).PK().Raw()).To(Equal(id.String()))
		})
	})
})
