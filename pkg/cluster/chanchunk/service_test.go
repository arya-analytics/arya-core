package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanchunk/mock"
	"github.com/arya-analytics/aryacore/pkg/cluster/internal"
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
	"strconv"
	"strings"
)

var _ = Describe("Service", func() {
	var (
		remoteSvc     chanchunk.ServiceRemote
		localSvc      chanchunk.ServiceLocal
		svc           *chanchunk.Service
		pool          rpc.Pool
		server        *mock.Server
		grpcServer    *grpc.Server
		lis           net.Listener
		node          = &storage.Node{ID: 1}
		channelConfig = &storage.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
		rangeX        = &storage.Range{
			ID: uuid.New(),
		}
		rangeReplica = &storage.RangeReplica{
			ID:      uuid.New(),
			RangeID: rangeX.ID,
			NodeID:  node.ID,
		}
		channelChunk = &storage.ChannelChunk{
			ID:              uuid.New(),
			RangeID:         rangeX.ID,
			ChannelConfigID: channelConfig.ID,
		}
		channelChunkReplica = &storage.ChannelChunkReplica{
			RangeReplicaID: rangeReplica.ID,
			ChannelChunkID: channelChunk.ID,
			Telem:          telem.NewBulk([]byte{}),
		}
		items = []interface{}{
			node,
			channelConfig,
			rangeX,
			rangeReplica,
		}
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
	Describe("Channel Chunk", func() {
		DescribeTable("Should Create + Retrieve + Delete correctly",
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

				deleteQR := internal.NewQueryRequest(
					internal.QueryVariantDelete,
					resRfl,
				)
				internal.NewPKQueryOpt(deleteQR, rfl.PKChain().Raw())

				By("Being able to handle the delete query")
				Expect(svc.CanHandle(deleteQR)).To(BeTrue())

				By("Executing the delete query")
				Expect(svc.Exec(ctx, deleteQR)).To(BeNil())

				resCCTwo := &storage.ChannelChunk{}
				retrieveQRTwo := internal.NewQueryRequest(
					internal.QueryVariantRetrieve,
					model.NewReflect(resCCTwo),
				)
				internal.NewPKQueryOpt(retrieveQRTwo, rfl.PKChain()[0].Raw())

				By("Not being able to be re-retrieved")
				rTwoErr := svc.Exec(ctx, retrieveQR)
				if rTwoErr != nil {
					Expect(rTwoErr.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
				} else {
					Expect(model.NewPK(resCCTwo.ID).IsZero()).To(BeTrue())
				}
			},
			Entry("Single Chunk", channelChunk, &storage.ChannelChunk{}),
			Entry("Slice of Chunks", &[]*storage.ChannelChunk{channelChunk}, &[]*storage.ChannelChunk{}),
		)
	})
	Describe("Channel Chunk Replica", func() {
		JustBeforeEach(func() {
			chunkCreateQR := internal.NewQueryRequest(
				internal.QueryVariantCreate,
				model.NewReflect(channelChunk),
			)
			Expect(svc.CanHandle(chunkCreateQR)).To(BeTrue())
			Expect(svc.Exec(ctx, chunkCreateQR)).To(BeNil())
		})
		JustAfterEach(func() {
			deleteQR := internal.NewQueryRequest(
				internal.QueryVariantDelete,
				model.NewReflect(channelChunk),
			)
			internal.NewPKQueryOpt(deleteQR, channelChunk.ID)
			Expect(svc.CanHandle(deleteQR)).To(BeTrue())
			err := svc.Exec(ctx, deleteQR)
			Expect(err).To(BeNil())
		})
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

				deleteQR := internal.NewQueryRequest(
					internal.QueryVariantDelete,
					resRfl,
				)
				internal.NewPKQueryOpt(deleteQR, rfl.PKChain().Raw())

				By("Being able to handle the delete query")
				Expect(svc.CanHandle(deleteQR)).To(BeTrue())

				By("Executing the delete query")
				Expect(svc.Exec(ctx, deleteQR)).To(BeNil())

				resCCTwo := &storage.ChannelChunkReplica{}
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
			Entry("Single Chunk Replica", channelChunkReplica, &storage.ChannelChunkReplica{}),
			Entry("Multiple Chunk Replicas", &[]*storage.ChannelChunkReplica{channelChunkReplica}, &[]*storage.ChannelChunkReplica{}),
		)
	})
})
