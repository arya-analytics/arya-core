package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("QueryTsRetrieve", func() {
	var (
		node          *models.Node
		channelConfig *models.ChannelConfig
		sample        *models.ChannelSample
		samples       []*models.ChannelSample
	)
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
	})
	JustBeforeEach(func() {
		nErr := store.NewCreate().Model(node).Exec(ctx)
		Expect(nErr).To(BeNil())
		cErr := store.NewCreate().Model(channelConfig).Exec(ctx)
		Expect(cErr).To(BeNil())
		sErr := store.NewTSCreate().Series().Model(channelConfig).Exec(ctx)
		Expect(sErr).To(BeNil())
	})
	JustAfterEach(func() {
		cErr := store.NewDelete().Model(channelConfig).WherePK(channelConfig.ID).Exec(ctx)
		Expect(cErr).To(BeNil())
		nErr := store.NewDelete().Model(node).WherePK(node.ID).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	Describe("Standard usage", func() {
		Describe("Retrieving a sample", func() {
			BeforeEach(func() {
				sample = &models.ChannelSample{
					ChannelConfigID: channelConfig.ID,
					Value:           124.7,
					Timestamp:       telem.NewTimeStamp(time.Now()),
				}
			})
			JustBeforeEach(func() {
				err := store.NewTSCreate().Sample().Model(sample).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should Retrieve the correct sample", func() {
				resSample := &models.ChannelSample{}
				err := store.NewTSRetrieve().Model(resSample).WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resSample.ChannelConfigID).To(Equal(channelConfig.ID))
			})
			It("Should Retrieve the correct sample", func() {
			})
		})
		Describe("Retrieving a sample by time rng", func() {
			BeforeEach(func() {
				samples = []*models.ChannelSample{
					{
						ChannelConfigID: channelConfig.ID,
						Value:           12412.3,
						Timestamp:       telem.NewTimeStamp(time.Now()),
					},
					{
						ChannelConfigID: channelConfig.ID,
						Value:           417.5,
						Timestamp:       telem.NewTimeStamp(time.Now().Add(-600 * time.Millisecond)),
					},
					{
						ChannelConfigID: channelConfig.ID,
						Value:           482.5,
						Timestamp:       telem.NewTimeStamp(time.Now().Add(-1600 * time.Second)),
					},
				}
			})
			JustBeforeEach(func() {
				err := store.NewTSCreate().Sample().Model(&samples).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should Retrieve the samples correctly", func() {
				var resSamples []*models.ChannelSample
				sampleTime := samples[0].Timestamp.ToTime()
				fromTs := sampleTime.Add(-800 * time.Millisecond)
				toTs := sampleTime.Add(500 * time.Millisecond)
				err := store.NewTSRetrieve().Model(&resSamples).WherePK(
					channelConfig.ID).WhereTimeRange(fromTs.UnixMicro(), toTs.UnixMicro()).
					Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resSamples).To(HaveLen(2))
				for _, s := range resSamples {
					Expect(s.Timestamp).Should(BeNumerically("<", telem.NewTimeStamp(toTs)))
					Expect(s.Timestamp).Should(BeNumerically(">", telem.NewTimeStamp(fromTs)))
				}
			})
		})
	})
})
