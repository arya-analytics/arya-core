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
		node          = &models.Node{IsHost: true}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
		items         = []interface{}{node, channelConfig}
	)
	BeforeEach(func() {
		clust = cluster.New()
		var lisErr error
		lis, lisErr = net.Listen("tcp", "localhost:0")
		Expect(lisErr).To(BeNil())
		port, pErr := strconv.Atoi(strings.Split(lis.Addr().String(), ":")[1])
		Expect(pErr).To(BeNil())
		node.RPCPort = port
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
	Describe("Node Is Local", func() {
		It("Should create a stream of samples correctly", func() {
			var samples = []*models.ChannelSample{
				{
					ChannelConfigID: channelConfig.ID,
					Timestamp:       telem.NewTimeStamp(time.Now()),
					Value:           1.0,
				},
				{
					ChannelConfigID: channelConfig.ID,
					Timestamp:       telem.NewTimeStamp(time.Now().Add(1 * time.Second)),
					Value:           2.0,
				},
			}
			c := make(chan *models.ChannelSample)
			sRfl := model.NewReflect(&c)
			errors := make(chan error)
			go func() {
				defer GinkgoRecover()
				Expect(<-errors).To(BeNil())
			}()
			tsquery.NewCreate().Model(sRfl).BindExec(clust.Exec).GoExec(ctx, errors)
			for _, s := range samples {
				c <- s
			}
			time.Sleep(20 * time.Millisecond)
			var resSamples []*models.ChannelSample
			Expect(persist.NewRetrieve().Model(&resSamples).Exec(ctx)).To(BeNil())
			Expect(samples).To(HaveLen(2))

		})
	})
})
