package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("QueryTsCreate", func() {
	var (
		series *models.ChannelConfig
		sample *models.ChannelSample
	)
	BeforeEach(func() {
		series = &models.ChannelConfig{ID: uuid.New()}
	})
	JustBeforeEach(func() {
		err := engine.NewTSCreate(adapter).Series().Model(series).Exec(ctx)
		Expect(err).To(BeNil())
	})
	Describe("Standard Usage", func() {
		Describe("Create a new series", func() {
			It("Should exist after creation", func() {
				exists, err := engine.NewTSRetrieve(adapter).SeriesExists(ctx, series.ID)
				Expect(err).To(BeNil())
				Expect(exists).To(BeTrue())
			})
		})
		Describe("Create a new sample", func() {
			Context("Single sample", func() {
				JustBeforeEach(func() {
					err := engine.NewTSCreate(adapter).Sample().Model(sample).Exec(ctx)
					Expect(err).To(BeNil())
				})
				BeforeEach(func() {
					sample = &models.ChannelSample{
						Timestamp:       telem.NewTimeStamp(time.Now()),
						Value:           123.2,
						ChannelConfigID: series.ID,
					}
				})
				It("Should be able to re-retrieve the sample after creation", func() {
					var resSamples []*models.ChannelSample
					rErr := engine.NewTSRetrieve(adapter).Model(&resSamples).
						WherePK(series.ID).Exec(
						ctx)
					Expect(rErr).To(BeNil())
					Expect(resSamples).To(HaveLen(1))
					Expect(resSamples[0].ChannelConfigID).To(Equal(series.ID))
					Expect(resSamples[0].Value).To(Equal(sample.Value))
					Expect(resSamples[0].Timestamp).To(Equal(sample.Timestamp))
				})
			})
		})
	})
	Describe("Edge cases + errors", func() {
		Describe("Not selecting a variant", func() {
			It("Should return the correct storage error", func() {
				err := engine.NewTSCreate(adapter).Model(series).Exec(ctx)
				Expect(err).ToNot(BeNil())
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeInvalidArgs))
			})
		})
	})
})
