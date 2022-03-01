package chanchunk_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/telem/chanchunk"
	"github.com/arya-analytics/aryacore/pkg/telem/rng"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("QueryStreamCreate", func() {
	var (
		node   *models.Node
		config *models.ChannelConfig
		svc    *chanchunk.Service
		items  []interface{}
	)
	BeforeEach(func() {
		rngObs := rng.NewObserveMem([]rng.ObservedRange{})
		rngPst := rng.NewPersistCluster(clust)
		rngSvc := rng.NewService(rngObs, rngPst)
		svc = chanchunk.NewService(clust, rngSvc)
		node = &models.Node{ID: 1}
		config = &models.ChannelConfig{
			Name:           "Awesome Channel",
			NodeID:         node.ID,
			DataRate:       telem.DataRate(1),
			DataType:       telem.DataTypeFloat64,
			ConflictPolicy: models.ChannelConflictPolicyDiscard,
		}
		items = []interface{}{node, config}
	})
	JustBeforeEach(func() {
		for _, item := range items {
			Expect(clust.NewCreate().Model(item).Exec(ctx)).To(BeNil())
		}
	})
	Describe("Standard Usage", func() {
		It("Should create a new telemetry chunk correctly", func() {
			By("Creating the stream")
			stream := svc.NewStreamCreate()

			By("Starting the stream")
			go stream.Start(ctx, config.ID)

			var streamError error
			go func() {
				streamError = <-stream.Errors()
			}()

			data := telem.NewChunkData([]byte{})
			Expect(data.WriteData([]float64{1, 2, 3, 4})).To(BeNil())

			By("Sending a new chunk")
			stream.Send(telem.TimeStamp(0), data)

			By("Closing the stream")
			stream.Close()

			By("Being error free")
			Expect(streamError).To(BeNil())

			By("Retrieving the chunk after creation")
			resCC := &models.ChannelChunk{}
			Expect(clust.NewRetrieve().
				Model(resCC).
				WhereFields(query.WhereFields{"StartTS": telem.TimeStamp(0)}).
				Exec(ctx)).To(BeNil())
			Expect(resCC.Size).To(Equal(4))
		})
	})
})
