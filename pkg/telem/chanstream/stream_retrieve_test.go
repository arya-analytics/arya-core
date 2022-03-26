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
	"time"

	"github.com/arya-analytics/aryacore/pkg/telem/chanstream"
)

var _ = Describe("StreamRetrieve", func() {
	var (
		clust         cluster.Cluster
		persist       *chanstreammock.Persist
		svc           *chanstream.Service
		nodeOne       = &models.Node{ID: 1, IsHost: true, Address: "localhost:26257"}
		channelConfig = &models.ChannelConfig{NodeID: nodeOne.ID, ID: uuid.New()}
		items         = []interface{}{nodeOne, channelConfig}
		sampleCount   = 1000
		samples       []*models.ChannelSample
	)
	BeforeEach(func() {
		clust = cluster.New()
		pool := clustermock.NewNodeRPCPool()
		ds := querymock.NewDataSourceMem()
		persist = &chanstreammock.Persist{DataSourceMem: ds}
		cSvc := clusterchanstream.NewService(persist.Exec, clusterchanstream.NewRemoteRPC(pool))
		clust.BindService(cSvc)
		clust.BindService(&clustermock.Persist{DataSourceMem: ds})
		svc = chanstream.NewService(clust.Exec)
		samples = []*models.ChannelSample{}
		for i := 0; i < sampleCount; i++ {
			samples = append(samples, &models.ChannelSample{
				ChannelConfigID: channelConfig.ID,
				Timestamp:       telem.NewTimeStamp(time.Now().Add(time.Duration(i) * time.Second)),
				Value:           float64(i),
			})
		}
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
		pkc := model.NewPKChain([]uuid.UUID{channelConfig.ID})
		stream := svc.NewStreamRetrieve().WherePKC(pkc)
		c := stream.Start(ctx)
		var resSamples []*models.ChannelSample
		t := time.NewTimer(200 * time.Millisecond)
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
})
