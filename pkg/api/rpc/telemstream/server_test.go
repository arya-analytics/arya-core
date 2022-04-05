package telemstream_test

import (
	api "github.com/arya-analytics/aryacore/pkg/api/rpc/gen/proto/go/telemstream/v1"
	"github.com/arya-analytics/aryacore/pkg/api/rpc/telemstream"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	clusterchanstream "github.com/arya-analytics/aryacore/pkg/cluster/chanstream"
	chanstreammock "github.com/arya-analytics/aryacore/pkg/cluster/chanstream/mock"
	clustermock "github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	querymock "github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/route"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"strings"
	"time"
)

var _ = Describe("Server", func() {
	var (
		clust   cluster.Cluster
		persist *chanstreammock.Persist
		svc     *chanstream.Service
		nodeOne = &models.Node{ID: 1, Address: "localhost:26257"}
		nodeTwo = &models.Node{ID: 2, Address: "localhost:26257"}

		lis        net.Listener
		grpcServer *grpc.Server

		ccOne       = &models.ChannelConfig{NodeID: nodeOne.ID, ID: uuid.New()}
		ccTwo       = &models.ChannelConfig{NodeID: nodeTwo.ID, ID: uuid.New()}
		items       = []interface{}{nodeOne, nodeTwo, ccOne, ccTwo}
		sampleCount = 1000
		samples     []*models.ChannelSample

		cl api.TelemStreamServiceClient
	)
	BeforeEach(func() {
		clust = cluster.New()
		var lisErr error
		lis, lisErr = net.Listen("tcp", "localhost:0")
		Expect(lisErr).To(BeNil())
		port, pErr := strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])
		Expect(pErr).To(BeNil())
		nodeTwo.RPCPort = port
		pool := clustermock.NewNodeRPCPool()
		ds := querymock.NewDataSourceMem()
		persist = &chanstreammock.Persist{DataSourceMem: ds}
		svc = chanstream.NewService(clust.Exec)
		clusterServer := clusterchanstream.NewServerRPC(persist.Exec)

		cSvc := clusterchanstream.NewService(persist.Exec, clusterchanstream.NewRemoteRPC(pool))
		clust.BindService(cSvc)
		clust.BindService(&clustermock.Persist{DataSourceMem: ds})
		grpcServer = grpc.NewServer()
		clusterServer.BindTo(grpcServer)

		apiServer := telemstream.NewServer(svc)

		apiServer.BindTo(grpcServer)

		samples = []*models.ChannelSample{}
		for i := 0; i < sampleCount; i++ {
			samples = append(samples, &models.ChannelSample{
				ChannelConfigID: ccOne.ID,
				Timestamp:       telem.NewTimeStamp(time.Now().Add(time.Duration(i) * time.Second)),
				Value:           float64(i),
			})
			samples = append(samples, &models.ChannelSample{
				ChannelConfigID: ccTwo.ID,
				Timestamp:       telem.NewTimeStamp(time.Now().Add(time.Duration(i) * time.Second)),
				Value:           float64(i),
			})
		}

	})
	JustBeforeEach(func() {
		go func() {
			defer GinkgoRecover()
			Expect(grpcServer.Serve(lis)).To(BeNil())
		}()
		conn, err := grpc.Dial(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		Expect(err).To(BeNil())
		cl = api.NewTelemStreamServiceClient(conn)
	})
	JustBeforeEach(func() {
		for _, m := range items {
			Expect(persist.NewCreate().Model(m).Exec(ctx)).To(BeNil())
		}
	})
	JustAfterEach(func() {
		for _, m := range items {
			Expect(persist.NewDelete().Model(m).WherePK(model.NewReflect(m).PK()).Exec(ctx)).To(BeNil())
		}
	})
	BeforeEach(func() { persist.AddQueryHook(clustermock.HostInterceptQueryHook(1)) })
	AfterEach(func() { persist.ClearQueryHooks() })
	Describe("Retrieving Samples", func() {
		Describe("Basic Usage", func() {
			It("Should retrieve a stream of samples correctly", func() {
				stream, err := cl.Retrieve(ctx)
				Expect(err).To(BeNil())
				stream.Send(&api.RetrieveRequest{
					PKC: model.NewPKChain([]uuid.UUID{ccOne.ID}).Strings(),
				})
				t := time.NewTimer(110 * time.Millisecond)

				var resSamples []*api.TelemSample

			o:
				for {
					select {
					case <-t.C:
						break o
					default:
						res, err := stream.Recv()
						Expect(err).To(BeNil())
						Expect(res.Error.Message).To(BeZero())
						resSamples = append(resSamples, res.Sample)
					}
				}
				Expect(stream.CloseSend()).To(Succeed())
				Expect(len(resSamples)).To(BeNumerically(">=", 10))
			})
		})
		Describe("Updating The Request", func() {
			It("Should allow the caller to update the request and receive new data", func() {
				stream, err := cl.Retrieve(ctx)
				Expect(err).To(BeNil())
				Expect(stream.Send(&api.RetrieveRequest{
					PKC: model.NewPKChain([]uuid.UUID{ccOne.ID}).Strings(),
				})).To(Succeed())

				t := time.NewTimer(110 * time.Millisecond)

				var resSamples []*api.TelemSample

			o:
				for {
					select {
					case <-t.C:
						break o
					default:
						res, err := stream.Recv()
						Expect(err).To(BeNil())
						Expect(res.Error.Message).To(BeZero())
						resSamples = append(resSamples, res.Sample)
					}
				}

				Expect(len(resSamples)).To(BeNumerically(">=", 10))

				Expect(stream.Send(&api.RetrieveRequest{
					PKC: model.NewPKChain([]uuid.UUID{ccOne.ID, ccTwo.ID}).Strings(),
				})).To(Succeed())

				time.Sleep(5 * time.Millisecond)

				var roundTwoResSamples []*api.TelemSample

				t2 := time.NewTimer(110 * time.Millisecond)

			o2:
				for {
					select {
					case <-t2.C:
						break o2
					default:
						res, err := stream.Recv()
						Expect(err).To(BeNil())
						Expect(res.Error.Message).To(BeZero())
						roundTwoResSamples = append(roundTwoResSamples, res.Sample)
					}
				}

				Expect(len(roundTwoResSamples)).To(BeNumerically(">=", 20))

				batched := route.BatchModel[string](model.NewReflect(&roundTwoResSamples), "ChannelConfigID")
				Expect(batched[ccOne.ID.String()].ChainValue().Len()).To(BeNumerically(">=", 10))
				Expect(batched[ccTwo.ID.String()].ChainValue().Len()).To(BeNumerically(">=", 10))
			})
		})
	})
	Describe("Creating Samples", func() {
		Describe("Basic Usage", func() {
			It("Should create a stream of samples correctly", func() {
				stream, err := cl.Create(ctx)
				Expect(err).To(BeNil())

				go func() {
					defer GinkgoRecover()
					for {
						res, err := stream.Recv()
						Expect(err).To(BeNil())
						Expect(res.Error.Message).To(BeZero())
					}
				}()

				for _, s := range samples {
					apiS := &api.TelemSample{
						ChannelConfigID: s.ChannelConfigID.String(),
						Value:           s.Value,
						Timestamp:       int64(s.Timestamp),
					}
					Expect(stream.Send(&api.CreateRequest{Sample: apiS})).To(Succeed())
				}

				Expect(stream.CloseSend()).To(Succeed())
				time.Sleep(100 * time.Millisecond)

				var resSamples []*models.ChannelSample
				Expect(persist.NewRetrieve().Model(&resSamples).WherePK(ccOne.ID).Exec(ctx)).To(Succeed())
				Expect(resSamples).To(HaveLen(1000))
			})
		})
	})
})
