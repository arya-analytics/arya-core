package chanstream_test

import (
	"context"
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
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
)

var _ = Describe("streamRetrieve", func() {
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
		server := clusterchanstream.NewServerRPC(persist.Exec)

		cSvc := clusterchanstream.NewService(persist.Exec, clusterchanstream.NewRemoteRPC(pool))
		clust.BindService(cSvc)
		clust.BindService(&clustermock.Persist{DataSourceMem: ds})
		grpcServer = grpc.NewServer()
		server.BindTo(grpcServer)
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
	BeforeEach(func() { persist.AddQueryHook(clustermock.HostInterceptQueryHook(1)) })
	AfterEach(func() { persist.ClearQueryHooks() })
	It("Should retrieve a stream of samples correctly", func() {
		pkc := model.NewPKChain([]uuid.UUID{ccOne.ID})
		c := make(chan *models.ChannelSample)
		aCtx, cancel := context.WithCancel(ctx)
		stream, err := svc.
			NewTSRetrieve().
			Model(&c).
			WherePKs(pkc).
			Stream(aCtx)
		Expect(err).To(BeNil())
		var resSamples []*models.ChannelSample
		t := time.NewTimer(100 * time.Millisecond)
		go func() {
			defer GinkgoRecover()
			Expect(<-stream.Errors).To(BeNil())
		}()
	o:
		for {
			select {
			case s := <-c:
				resSamples = append(resSamples, s)
			case <-t.C:
				cancel()
				break o
			}
		}
		Expect(len(resSamples)).To(BeNumerically(">", 7))
	})
	It("Should serve multiple concurrent streams correctly", func() {
		pkc := model.NewPKChain([]uuid.UUID{ccOne.ID})
		pkc2 := model.NewPKChain([]uuid.UUID{ccTwo.ID})
		c1 := make(chan *models.ChannelSample)
		c2 := make(chan *models.ChannelSample)
		aCtx, cancel := context.WithCancel(ctx)
		stream, err := svc.NewTSRetrieve().Model(&c1).WherePKs(pkc).Stream(aCtx)
		Expect(err).To(BeNil())
		stream2, err := svc.NewTSRetrieve().Model(&c2).WherePKs(pkc2).Stream(aCtx)
		Expect(err).To(BeNil())

		go func() {
			defer GinkgoRecover()
			var err error
			select {
			case err = <-stream.Errors:
			case err = <-stream2.Errors:
			}
			Expect(err).To(BeNil())
		}()

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
				cancel()
				break o
			}
		}

		Expect(len(resSamples2)).To(BeNumerically(">", 8))
		Expect(len(resSamples2)).To(BeNumerically("<", 12))
	})
})
