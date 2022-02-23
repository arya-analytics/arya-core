package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk/mock"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
	clustermock "github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
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
		remoteSvc     chanchunk.ServiceRemote
		localSvc      chanchunk.ServiceLocal
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
			channelChunk,
		}
	)
	BeforeEach(func() {
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
		DescribeTable("Should Create + Retrieve + Delete the chunk replica correctly",
			func(cc interface{}, resCC interface{}) {
				rfl, resRfl := model.NewReflect(cc), model.NewReflect(resCC)
				createQR := &internal.QueryRequest{
					Variant: internal.QueryVariantCreate,
					Model:   rfl,
				}
				By("Being able to handle the create query")
				Expect(svc.CanHandle(createQR)).To(BeTrue())
				By("Being able to execute the create query")
				Expect(svc.Exec(ctx, createQR)).To(BeNil())

				retrieveQR := internal.NewQueryRequest(
					internal.QueryVariantRetrieve,
					model.NewReflect(resCC),
				)
				internal.NewPKQueryOpt(retrieveQR, rfl.PKChain().Raw())
				By("Being able to handle the retrieve query")
				Expect(svc.CanHandle(retrieveQR)).To(BeTrue())

				By("Executing the retrieve query")
				Expect(svc.Exec(ctx, retrieveQR)).To(BeNil())

				By("Retrieving the correct item")
				Expect(resRfl.PKChain()).To(Equal(rfl.PKChain()))
				resRflItem, ok := resRfl.ValueByPK(resRfl.PKChain()[0])
				Expect(ok).To(BeTrue())
				Expect(resRflItem.Pointer().(*models.ChannelChunkReplica).Telem.Bytes()).To(Equal([]byte("randomdata")))

				deleteQR := internal.NewQueryRequest(
					internal.QueryVariantDelete,
					resRfl,
				)
				internal.NewPKQueryOpt(deleteQR, rfl.PKChain().Raw())

				By("Being able to handle the delete query")
				Expect(svc.CanHandle(deleteQR)).To(BeTrue())

				By("Executing the delete query")
				Expect(svc.Exec(ctx, deleteQR)).To(BeNil())

				resCCTwo := &models.ChannelChunkReplica{}
				retrieveQRTwo := internal.NewQueryRequest(
					internal.QueryVariantRetrieve,
					model.NewReflect(resCCTwo),
				)
				internal.NewPKQueryOpt(retrieveQRTwo, rfl.PKChain()[0].Raw())

				By("Not being able to be re-retrieved")
				rTwoErr := svc.Exec(ctx, retrieveQRTwo)
				if rTwoErr != nil {
					Expect(rTwoErr.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
				} else {
					Expect(model.NewPK(resCCTwo.ID).IsZero()).To(BeTrue())
				}
			},
			Entry("Single Chunk Replica", channelChunkReplica, &models.ChannelChunkReplica{}),
			Entry("Multiple Chunk Replicas", &[]*models.ChannelChunkReplica{channelChunkReplica}, &[]*models.ChannelChunkReplica{}),
		)
	})
})
