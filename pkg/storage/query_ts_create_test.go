package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("QueryTSCreate", func() {
	var (
		node          *models.Node
		channelConfig *models.ChannelConfig
	)
	BeforeEach(func() {
		node = &models.Node{ID: 1}
		channelConfig = &models.ChannelConfig{NodeID: 1, ID: uuid.New()}
	})
	JustBeforeEach(func() {
		nErr := store.NewCreate().Model(node).Exec(ctx)
		Expect(nErr).To(BeNil())
		cErr := store.NewCreate().Model(channelConfig).Exec(ctx)
		Expect(cErr).To(BeNil())
	})
	JustAfterEach(func() {
		cErr := store.NewDelete().Model(channelConfig).WherePK(channelConfig.ID).
			Exec(ctx)
		Expect(cErr).To(BeNil())
		nErr := store.NewDelete().Model(node).WherePK(node.ID).Exec(ctx)
		Expect(nErr).To(BeNil())
	})
	Describe("Standard Usage", func() {
		Describe("Create a new series", func() {
			It("Should create the series correctly", func() {
				err := store.NewTSCreate().Series().Model(channelConfig).Exec(ctx)
				Expect(err).To(BeNil())
				exists, rErr := store.NewTSRetrieve().SeriesExists(ctx, channelConfig.ID)
				Expect(rErr).To(BeNil())
				Expect(exists).To(BeTrue())
			})
		})
		Describe("Create a new sample", func() {
			JustBeforeEach(func() {
				sErr := store.NewTSCreate().Series().Model(channelConfig).Exec(ctx)
				if sErr != nil {
					Expect(sErr.(query.Error).Type).To(Equal(query.ErrorTypeUniqueViolation))
				} else {
					Expect(sErr).To(BeNil())
				}
			})
			Context("Single Sample", func() {
				var sample *models.ChannelSample
				BeforeEach(func() {
					sample = &models.ChannelSample{
						ChannelConfigID: channelConfig.ID,
						Value:           31.6,
						Timestamp:       telem.NewTimeStamp(time.Now()),
					}
				})
				It("Should create the sample correctly", func() {
					err := store.NewTSCreate().Sample().Model(sample).Exec(ctx)
					Expect(err).To(BeNil())
					resSample := &models.ChannelSample{}
					rErr := store.NewTSRetrieve().Model(resSample).WherePK(channelConfig.ID).Exec(ctx)
					Expect(rErr).To(BeNil())
					Expect(resSample.Timestamp).To(Equal(sample.Timestamp))
				})
			})
			Context("Multiple Samples", func() {
				var samples []*models.ChannelSample
				var channelConfigTwo *models.ChannelConfig
				BeforeEach(func() {
					channelConfigTwo = &models.ChannelConfig{
						ID:     uuid.New(),
						Name:   "SG 43",
						NodeID: node.ID,
					}
					samples = []*models.ChannelSample{
						{
							ChannelConfigID: channelConfig.ID,
							Value:           3124.4,
							Timestamp:       telem.NewTimeStamp(time.Now()),
						},
						{
							ChannelConfigID: channelConfigTwo.ID,
							Value:           3124.4,
							Timestamp:       telem.NewTimeStamp(time.Now().Add(1 * time.Second)),
						},
						{
							ChannelConfigID: channelConfig.ID,
							Value:           3124.4,
							Timestamp:       telem.NewTimeStamp(time.Now().Add(1 * time.Second)),
						},
					}
				})
				JustBeforeEach(func() {
					cErr := store.NewCreate().Model(channelConfigTwo).Exec(ctx)
					Expect(cErr).To(BeNil())
					seriesErr := store.NewTSCreate().Series().Model(channelConfigTwo).Exec(ctx)
					Expect(seriesErr).To(BeNil())
				})
				It("The samples should be able to be retrieved after creation", func() {
					var resSamples []*models.ChannelSample
					cErr := store.NewTSCreate().Sample().Model(&samples).Exec(ctx)
					Expect(cErr).To(BeNil())
					rErr := store.NewTSRetrieve().Model(&resSamples).WherePKs([]uuid.
						UUID{channelConfigTwo.ID, channelConfig.ID}).AllTimeRange().Exec(ctx)
					Expect(rErr).To(BeNil())
					Expect(resSamples).To(HaveLen(3))
				})
			})
		})
	})
})
