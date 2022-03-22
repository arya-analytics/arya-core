package chanstream_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanstream"
	"github.com/arya-analytics/aryacore/pkg/cluster/chanstream/mock"
	clustermock "github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	modelMock "github.com/arya-analytics/aryacore/pkg/util/model/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	querymock "github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/query/tsquery"
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
		remoteSvc     *chanstream.ServiceRemoteRPC
		server        *chanstream.ServerRPC
		pool          *cluster.NodeRPCPool
		grpcServer    *grpc.Server
		svc           *chanstream.Service
		persist       *mock.Persist
		lis           net.Listener
		nodeOne       = &models.Node{ID: 1, IsHost: true, Address: "localhost:26257"}
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
		remoteSvc = chanstream.NewServiceRemoteRPC(pool)
		svc = chanstream.NewService(persist.Exec, remoteSvc)
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
		//for _, m := range items {
		//	Expect(persist.NewDelete().Model(m).WherePK(model.NewReflect(m).PK()).Exec(ctx)).To(BeNil())
		//}
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
	Describe("Create", func() {
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
		Context("Node Is Local", func() {
			It("Should create a stream of samples correctly", func() {
				c := make(chan *models.ChannelSample)
				sRfl := model.NewReflect(&c)
				errors := make(chan error)
				go func() {
					panic(<-errors)
				}()
				tsquery.NewCreate().Model(sRfl).BindExec(clust.Exec).GoExec(ctx, errors)
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
				nodeOne.IsHost = false
			})
			It("Should create a stream of samples correctly", func() {
				c := make(chan *models.ChannelSample)
				sRfl := model.NewReflect(&c)
				errors := make(chan error)
				go func() {
					panic(<-errors)
				}()
				tsquery.NewCreate().Model(sRfl).BindExec(clust.Exec).GoExec(ctx, errors)
				for _, s := range samples {
					c <- s
				}
				time.Sleep(20 * time.Millisecond)
				var resSamples []*models.ChannelSample
				Expect(persist.NewRetrieve().Model(&resSamples).Exec(ctx)).To(BeNil())
				Expect(samples).To(HaveLen(sampleCount))
			})
		})
	})
})
