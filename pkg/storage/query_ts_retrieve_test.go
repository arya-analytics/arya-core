package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("QueryTsRetrieve", func() {
	var (
		node          *storage.Node
		channelConfig *storage.ChannelConfig
		sample        *storage.ChannelSample
		samples       []*storage.ChannelSample
	)
	BeforeEach(func() {
		node = &storage.Node{ID: 1}
		channelConfig = &storage.ChannelConfig{NodeID: node.ID, ID: uuid.New()}
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
				sample = &storage.ChannelSample{
					ChannelConfigID: channelConfig.ID,
					Value:           124.7,
					Timestamp:       time.Now().Unix(),
				}
			})
			JustBeforeEach(func() {
				err := store.NewTSCreate().Sample().Model(sample).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve the correct sample", func() {
				resSample := &storage.ChannelSample{}
				err := store.NewTSRetrieve().Model(resSample).WherePK(channelConfig.
					ID).Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resSample.ChannelConfigID).To(Equal(channelConfig.ID))
			})
			It("Should retrieve the correct sample", func() {
			})
		})
		Describe("Retrieving a sample by time range", func() {
			BeforeEach(func() {
				samples = []*storage.ChannelSample{
					{
						ChannelConfigID: channelConfig.ID,
						Value:           12412.3,
						Timestamp:       time.Now().UnixNano(),
					},
					{
						ChannelConfigID: channelConfig.ID,
						Value:           417.5,
						Timestamp: time.Now().Add(-600 * time.Millisecond).
							UnixNano(),
					},
					{
						ChannelConfigID: channelConfig.ID,
						Value:           482.5,
						Timestamp:       time.Now().Add(-1600 * time.Second).UnixNano(),
					},
				}
			})
			JustBeforeEach(func() {
				err := store.NewTSCreate().Sample().Model(&samples).Exec(ctx)
				Expect(err).To(BeNil())
			})
			It("Should retrieve the samples correctly", func() {
				var resSamples []*storage.ChannelSample
				sampleTime := time.Unix(0, samples[0].Timestamp)
				fromTs := sampleTime.Add(-800 * time.Millisecond)
				toTs := sampleTime.Add(500 * time.Millisecond)
				err := store.NewTSRetrieve().Model(&resSamples).WherePK(
					channelConfig.ID).WhereTimeRange(fromTs.UnixNano(), toTs.UnixNano()).
					Exec(ctx)
				Expect(err).To(BeNil())
				Expect(resSamples).To(HaveLen(2))
				for _, s := range resSamples {
					Expect(s.Timestamp).Should(BeNumerically("<", toTs.UnixNano()))
					Expect(s.Timestamp).Should(BeNumerically(">", fromTs.UnixNano()))
				}
			})
		})
	})
})
