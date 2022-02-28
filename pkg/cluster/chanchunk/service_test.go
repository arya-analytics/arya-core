package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk/mock"
	clustermock "github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	mockTlm "github.com/arya-analytics/aryacore/pkg/util/telem/mock"
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
		persist       *mock.ServerRPCPersist
		server        *chanchunk.ServerRPC
		grpcServer    *grpc.Server
		lis           net.Listener
		ccr           *models.ChannelChunkReplica
		node          = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
		rng           = &models.Range{ID: uuid.New()}
		rr            = &models.RangeReplica{ID: uuid.New(), RangeID: rng.ID, NodeID: node.ID}
		cc            = &models.ChannelChunk{ID: uuid.New(), RangeID: rng.ID, ChannelConfigID: channelConfig.ID}
		items         = []interface{}{node, channelConfig, rng, rr, cc}
	)
	BeforeEach(func() {
		store.AddQueryHook(mock.HostInterceptQueryHook(2))
		clust = cluster.New()
		var lisErr error
		lis, lisErr = net.Listen("tcp", "localhost:0")
		Expect(lisErr).To(BeNil())
		port, pErr := strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])
		Expect(pErr).To(BeNil())
		node.RPCPort = port
		pool = clustermock.NewNodeRPCPool()
		remoteSvc = chanchunk.NewServiceRemoteRPC(pool)
		persist = &mock.ServerRPCPersist{}
		server = chanchunk.NewServerRPC(persist)
		grpcServer = grpc.NewServer()
		server.BindTo(grpcServer)
		localSvc = chanchunk.NewServiceLocalStorage(store)
		svc = chanchunk.NewService(localSvc, remoteSvc)
		clust.BindService(svc)
		ccr = &models.ChannelChunkReplica{
			ID:             uuid.New(),
			RangeReplicaID: rr.ID,
			ChannelChunkID: cc.ID,
			Telem:          telem.NewChunkData([]byte{}),
		}
		mockTlm.TelemBulkPopulateRandomFloat64(ccr.Telem, 500)
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
	Describe("Node Is Remote", func() {
		It("Should create the chunk replica correctly", func() {
			err := clust.NewCreate().Model(ccr).Exec(ctx)
			Expect(err).To(BeNil())
			Expect(persist.ChunkReplicas).To(HaveLen(1))
			Expect(persist.ChunkReplicas[0].ID).To(Equal(ccr.ID))
		})
		It("Should retrieve the chunk replica correctly", func() {
			cErr := clust.NewCreate().Model(ccr).Exec(ctx)
			ccr.Telem = telem.NewChunkData([]byte{})
			sCErr := store.NewCreate().Model(ccr).Exec(ctx)
			Expect(sCErr).To(BeNil())
			Expect(cErr).To(BeNil())
			resCCR := &models.ChannelChunkReplica{}
			rErr := clust.NewRetrieve().Model(resCCR).WherePK(ccr.ID).Exec(ctx)
			Expect(rErr).To(BeNil())
			Expect(resCCR.ID).To(Equal(ccr.ID))
		})
		It("Should delete the chunk replica correctly", func() {
			cErr := clust.NewCreate().Model(ccr).Exec(ctx)
			ccr.Telem = telem.NewChunkData([]byte{})
			sCErr := store.NewCreate().Model(ccr).Exec(ctx)
			Expect(sCErr).To(BeNil())
			Expect(cErr).To(BeNil())
			dErr := clust.NewDelete().Model(&models.ChannelChunkReplica{}).WherePK(ccr.ID).Exec(ctx)
			Expect(dErr).To(BeNil())
			Expect(persist.ChunkReplicas).To(HaveLen(0))
		})
	})
	Describe("Node Is Local", func() {
		BeforeEach(func() {
			store.AddQueryHook(mock.HostInterceptQueryHook(1))
		})
		It("Should create the chunk replica correctly", func() {
			err := clust.NewCreate().Model(ccr).Exec(ctx)
			Expect(err).To(BeNil())
		})
		It("Should retrieve the chunk replica correctly", func() {
			cErr := clust.NewCreate().Model(ccr).Exec(ctx)
			Expect(cErr).To(BeNil())
			resCCR := &models.ChannelChunkReplica{}
			rErr := clust.NewRetrieve().Model(resCCR).WherePK(ccr.ID).Exec(ctx)
			Expect(rErr).To(BeNil())
			Expect(resCCR.ID).To(Equal(ccr.ID))
		})
		It("Should retrieve only meta data when bulk telem field not specified", func() {
			cErr := clust.NewCreate().Model(ccr).Exec(ctx)
			Expect(cErr).To(BeNil())
			resCCR := &models.ChannelChunkReplica{}
			rErr := clust.NewRetrieve().Model(resCCR).WhereFields(model.WhereFields{"ID": ccr.ID}).Fields("ID").Exec(ctx)
			Expect(rErr).To(BeNil())
			Expect(resCCR.ID).To(Equal(ccr.ID))
			Expect(resCCR.Telem).To(BeNil())
		})
		It("Should delete the chunk replica correctly", func() {
			cErr := clust.NewCreate().Model(ccr).Exec(ctx)
			Expect(cErr).To(BeNil())
			dErr := clust.NewDelete().Model(&models.ChannelChunkReplica{}).WherePK(ccr.ID).Exec(ctx)
			Expect(dErr).To(BeNil())
			Expect(persist.ChunkReplicas).To(HaveLen(0))
		})
		It("Should update the replica correctly", func() {
			cErr := clust.NewCreate().Model(ccr).Exec(ctx)
			Expect(cErr).To(BeNil())
			rrTwo := &models.RangeReplica{ID: uuid.New(), RangeID: rng.ID, NodeID: node.ID}
			rrErr := store.NewCreate().Model(rrTwo).Exec(ctx)
			Expect(rrErr).To(BeNil())
			ccr.RangeReplicaID = rrTwo.ID
			uErr := clust.NewUpdate().Model(ccr).WherePK(ccr.ID).Fields("RangeReplicaID").Exec(ctx)
			Expect(uErr).To(BeNil())
		})
	})
})
