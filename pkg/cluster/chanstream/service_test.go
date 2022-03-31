package chanstream_test

import (
	"context"
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanstream"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanstream/mock"
	clustermock "github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	modelMock "github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	querymock "github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query/streamq"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"time"
)

var _ = Describe("Service", func() {
	var (
		clust         cluster.Cluster
		server        *chanstream.ServerRPC
		pool          *cluster.NodeRPCPool
		grpcServer    *grpc.Server
		svc           *chanstream.Service
		persist       *mock.Persist
		lis           net.Listener
		nodeOne       = &models.Node{ID: 1, Address: "localhost:26257"}
		channelConfig = &models.ChannelConfig{NodeID: nodeOne.ID, ID: uuid.New()}
		items         = []interface{}{nodeOne, channelConfig}
	)
	BeforeEach(func() {
		clust = cluster.New()
		var lisErr error
		lis, lisErr = net.Listen("tcp", "localhost:0")
		Expect(lisErr).To(BeNil())
		port, pErr := strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])
		Expect(pErr).To(BeNil())
		nodeOne.RPCPort = port
		pool = clustermock.NewNodeRPCPool()
		ds := querymock.NewDataSourceMem()
		persist = &mock.Persist{DataSourceMem: ds}
		server = chanstream.NewServerRPC(persist.Exec)
		grpcServer = grpc.NewServer()
		server.BindTo(grpcServer)
		svc = chanstream.NewService(persist.Exec, chanstream.NewRemoteRPC(pool))
		clust.BindService(svc)
		clust.BindService(&clustermock.Persist{DataSourceMem: ds})
	})
	JustBeforeEach(func() {
		go func() {
			defer GinkgoRecover()
			Expect(grpcServer.Serve(lis)).To(BeNil())
		}()
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
	AfterEach(func() {
		persist.ClearQueryHooks()
	})
	Describe("CanHandle", func() {
		It("Should return false for a query it can't handle", func() {
			c := make(chan *modelMock.ModelA)
			p := query.NewRetrieve().Model(&c).Pack()
			Expect(svc.CanHandle(p)).To(BeFalse())
		})
		It("Should return true for a query it can handle", func() {
			c := make(chan *models.ChannelSample)
			p := query.NewRetrieve().Model(&c).Pack()
			Expect(svc.CanHandle(p)).To(BeTrue())
		})
	})
	var (
		sampleCount = 1000
		samples     []*models.ChannelSample
	)
	BeforeEach(func() {
		samples = []*models.ChannelSample{}
		for i := 0; i < sampleCount; i++ {
			samples = append(samples, &models.ChannelSample{
				ChannelConfigID: channelConfig.ID,
				Timestamp:       telem.NewTimeStamp(time.Now().Add(time.Duration(i) * time.Second)),
				Value:           float64(i),
			})
		}
	})
	Describe("TSCreate", func() {
		Context("Node Is Local", func() {
			BeforeEach(func() {
				persist.AddQueryHook(clustermock.HostInterceptQueryHook(1))
			})
			It("Should create a stream of samples correctly", func() {
				c := make(chan *models.ChannelSample)
				sRfl := model.NewReflect(&c)
				stream, err := streamq.NewTSCreate().Model(sRfl).BindExec(clust.Exec).Stream(ctx)
				Expect(err).To(BeNil())
				go func() {
					panic(<-stream.Errors)
				}()
				for _, s := range samples {
					c <- s
				}
				time.Sleep(20 * time.Millisecond)
				var resSamples []*models.ChannelSample
				Expect(persist.NewRetrieve().Model(&resSamples).Exec(ctx)).To(BeNil())
				Expect(samples).To(HaveLen(sampleCount))
			})

		})

		Context("Node Is Remote", func() {
			BeforeEach(func() {
				persist.AddQueryHook(clustermock.HostInterceptQueryHook(2))
			})
			It("Should create a stream of samples correctly", func() {
				c := make(chan *models.ChannelSample)
				sRfl := model.NewReflect(&c)
				stream, err := streamq.NewTSCreate().Model(sRfl).BindExec(clust.Exec).Stream(ctx)
				Expect(err).To(BeNil())
				for _, s := range samples {
					c <- s
				}
				go func() {
					panic(<-stream.Errors)
				}()
				time.Sleep(20 * time.Millisecond)
				var resSamples []*models.ChannelSample
				Expect(persist.NewRetrieve().Model(&resSamples).Exec(ctx)).To(BeNil())
				Expect(resSamples).To(HaveLen(sampleCount))
			})
		})
	})

	Describe("TSRetrieve", func() {
		BeforeEach(func() {
			Expect(persist.NewCreate().Model(&samples).Exec(ctx)).To(BeNil())
		})
		Context("Node Is Local", func() {
			BeforeEach(func() {
				persist.AddQueryHook(clustermock.HostInterceptQueryHook(1))
			})
			It("Should retrieve a stream of samples correctly", func() {
				c := make(chan *models.ChannelSample)
				sRfl := model.NewReflect(&c)
				stream, err := streamq.NewTSRetrieve().Model(sRfl).WherePK(channelConfig.ID).BindExec(clust.Exec).Stream(ctx)
				Expect(err).To(BeNil())
				go func() {
					panic(<-stream.Errors)
				}()
				var resSamples []*models.ChannelSample
				t := time.NewTimer(100 * time.Millisecond)
			o:
				for {
					select {
					case s := <-c:
						resSamples = append(resSamples, s)
					case <-t.C:
						break o
					}
				}
				Expect(len(resSamples)).To(BeNumerically(">", 8))
			})
			It("Should retrieve a streamq of samples form multi channels correctly", func() {
				ccTwo := &models.ChannelConfig{ID: uuid.New(), NodeID: 1}
				Expect(persist.NewCreate().Model(ccTwo).Exec(ctx))
				c := make(chan *models.ChannelSample, 2)
				sRfl := model.NewReflect(&c)
				stream, err := streamq.
					NewTSRetrieve().
					Model(sRfl).
					WherePKs([]uuid.UUID{channelConfig.ID, ccTwo.ID}).
					BindExec(clust.Exec).
					Stream(ctx)
				Expect(err).To(BeNil())
				go func() {
					panic(<-stream.Errors)
				}()
				var resSamples []*models.ChannelSample
				t := time.NewTimer(100 * time.Millisecond)
			o:
				for {
					select {
					case s := <-c:
						resSamples = append(resSamples, s)
					case <-t.C:
						break o
					}
				}
				Expect(len(resSamples)).To(BeNumerically(">", 16))
			})
			It("Should stop retrieving samples after a context is canceled", func() {
				c := make(chan *models.ChannelSample)
				sRfl := model.NewReflect(&c)
				aCtx, cancel := context.WithCancel(ctx)
				stream, err := streamq.NewTSRetrieve().Model(sRfl).WherePK(channelConfig.ID).BindExec(clust.Exec).Stream(aCtx)
				Expect(err).To(BeNil())
				go func() {
					panic(<-stream.Errors)
				}()
				var resSamples []*models.ChannelSample
				go func() {
					for s := range c {
						resSamples = append(resSamples, s)
					}
				}()
				time.Sleep(20 * time.Millisecond)
				cancel()
				time.Sleep(50 * time.Millisecond)
				Expect(len(resSamples)).To(BeNumerically("<", 4))
			})
		})
		Context("Node Is Remote", func() {
			BeforeEach(func() {
				persist.AddQueryHook(clustermock.HostInterceptQueryHook(2))
			})
			It("Should retrieve a stream of samples correctly", func() {
				c := make(chan *models.ChannelSample)
				sRfl := model.NewReflect(&c)
				stream, err := streamq.NewTSRetrieve().Model(sRfl).WherePK(channelConfig.ID).BindExec(clust.Exec).Stream(ctx)
				Expect(err).To(BeNil())
				go func() {
					panic(<-stream.Errors)
				}()
				var resSamples []*models.ChannelSample
				t := time.NewTimer(100 * time.Millisecond)
			o:
				for {
					select {
					case s := <-c:
						resSamples = append(resSamples, s)
					case <-t.C:
						break o
					}
				}
				Expect(len(resSamples)).To(BeNumerically(">", 8))
			})
		})
	})
})
