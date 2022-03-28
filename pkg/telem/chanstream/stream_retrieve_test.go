package chanstream_test

import (
	"github.com/arya-analytics/aryacore/pkg/cluster"
	clusterchanstream "github.com/arya-analytics/aryacore/pkg/cluster/chanstream"
	chanstreammock "github.com/arya-analytics/aryacore/pkg/cluster/chanstream/mock"
	clustermock "github.com/arya-analytics/aryacore/pkg/cluster/mock"
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	querymock "github.com/arya-analytics/aryacore/pkg/util/query/mock"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
)

var _ = Describe("StreamRetrieve", func() {
	var (
		clust   cluster.Cluster
		persist *chanstreammock.Persist
		svc     *chanstream.Service
		nodeOne = &models.Node{ID: 1, IsHost: true, Address: "localhost:26257"}
		nodeTwo = &models.Node{ID: 2, IsHost: false, Address: "localhost:26257"}

		lis        net.Listener
		grpcServer *grpc.Server

		ccOne       = &models.ChannelConfig{NodeID: nodeOne.ID, ID: uuid.New()}
		ccTwo       = &models.ChannelConfig{NodeID: nodeTwo.ID, ID: uuid.New()}
		items       = []interface{}{nodeOne, nodeTwo, ccOne, ccTwo}
		sampleCount = 1000
		samples     []*models.ChannelSample
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
		server := clusterchanstream.NewServerRPC(persist.Exec)
		grpcServer = grpc.NewServer()
		server.BindTo(grpcServer)
		cSvc := clusterchanstream.NewService(persist.Exec, clusterchanstream.NewRemoteRPC(pool))
		clust.BindService(cSvc)
		clust.BindService(&clustermock.Persist{DataSourceMem: ds})
		svc = chanstream.NewService(clust.Exec)
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
	It("Should retrieve a stream of samples correctly", func() {
		pkc := model.NewPKChain([]uuid.UUID{ccOne.ID})
		stream := svc.NewStreamRetrieve().WherePKC(pkc)
		c := stream.Start(ctx)
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
		Expect(len(resSamples)).To(BeNumerically(">", 7))
	})
	FIt("Should serve multiple concurrent streams correctly", func() {
		pkc := model.NewPKChain([]uuid.UUID{ccOne.ID})
		pkc2 := model.NewPKChain([]uuid.UUID{ccTwo.ID})
		stream := svc.NewStreamRetrieve().WherePKC(pkc)
		stream2 := svc.NewStreamRetrieve().WherePKC(pkc2)
		c2 := stream2.Start(ctx)
		c1 := stream.Start(ctx)

		var (
			resSamples  []*models.ChannelSample
			resSamples2 []*models.ChannelSample
		)
		t := time.NewTimer(100 * time.Millisecond)
	o:
		for {
			select {
			case s := <-c1:
				resSamples = append(resSamples, s)
			case s := <-c2:
				resSamples2 = append(resSamples2, s)
			case <-t.C:
				break o
			}
		}
		log.Info(len(resSamples), len(resSamples2))
		Expect(len(resSamples)).To(BeNumerically(">", 8))
		Expect(len(resSamples2)).To(BeNumerically(">", 8))
	})
})
