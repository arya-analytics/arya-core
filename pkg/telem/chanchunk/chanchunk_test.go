package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	mockTlm "github.com/arya-analytics/aryacore/pkg/util/telem/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"time"
)

var _ = Describe("Service", func() {
	var (
		svc           *chanchunk.Service
		rngSVC        *rng.Service
		node          *models.Node
		channelConfig *models.ChannelConfig
		items         []interface{}
	)
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID}
		obs := rng.NewObserveMem([]rng.ObservedRange{})
		p := &rng.PersistCluster{Cluster: clust}
		rngSVC = rng.NewService(obs, p)
		svc = chanchunk.NewService(clust, rngSVC)
		items = []interface{}{node, channelConfig}
	})
	JustBeforeEach(func() {
		for _, item := range items {
			Expect(clust.NewCreate().Model(item).Exec(ctx)).To(BeNil())
		}
	})
	JustAfterEach(func() {
		for _, item := range items {
			Expect(clust.NewDelete().Model(item).WherePK(model.NewReflect(item).PK().Raw()).Exec(ctx)).To(BeNil())
		}
	})
	Describe("Standard Usage", func() {
		It("Should allow the create + retrieve of telemetry", func() {
			By("Creating telemetry")
			createStream, errChan := svc.CreateStream(ctx, channelConfig)
			startTs := time.Now().Add(-10 * time.Second).UnixMicro()
			for i := 0; i < 5; i++ {
				tlm := telem.NewBulk([]byte{})
				mockTlm.TelemBulkPopulateRandomFloat64(tlm, 100)
				ccr := &chanchunk.TelemChunk{StartTS: time.Now().UnixMicro(), Data: tlm}
				createStream <- ccr
			}
			close(createStream)
			Expect(<-errChan).To(Equal(io.EOF))
			endTs := time.Now().Add(10 * time.Second).UnixMicro()

			By("Retrieving telemetry")
			var resCCR []*chanchunk.TelemChunk
			resStream, resErrCHan := svc.RetrieveStream(ctx, channelConfig, chanchunk.RetrieveOpts{StartTS: startTs, EndTS: endTs})
			for ccr := range resStream {
				resCCR = append(resCCR, ccr)
			}
			Expect(<-resErrCHan).To(Equal(io.EOF))
			Expect(resCCR).To(HaveLen(5))
		})
	})
})
