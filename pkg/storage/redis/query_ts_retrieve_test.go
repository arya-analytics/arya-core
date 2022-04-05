package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/models"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	"github.com/arya-analytics/aryacore/pkg/util/telem"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("QueryTsRetrieve", func() {
	var (
		series  *models.ChannelConfig
		sample  *models.ChannelSample
		samples []*models.ChannelSample
	)
	BeforeEach(func() {
		series = &models.ChannelConfig{ID: uuid.New()}
	})
	JustBeforeEach(func() {
		err := engine.NewTSCreate().Model(series).Exec(ctx)
		Expect(err).To(BeNil())
	})
	Describe("Standard Usage", func() {
		Context("Single sample", func() {
			JustBeforeEach(func() {
				sampleErr := engine.NewTSCreate().Model(sample).Exec(ctx)
				Expect(sampleErr).To(BeNil())
			})
			BeforeEach(func() {
				sample = &models.ChannelSample{
					ChannelConfigID: series.ID,
					Value:           123.2,
					Timestamp:       telem.NewTimeStamp(time.Now()),
				}
			})
			Describe("Retrieving the most recent sample", func() {
				It("Should retrieve the correct item", func() {
					var resSample = &models.ChannelSample{}
					err := engine.NewTSRetrieve().Model(resSample).WherePK(series.ID).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(sample).To(Equal(resSample))
				})
			})
		})
		Context("Multiple Samples", func() {
			JustBeforeEach(func() {
				err := engine.NewTSCreate().Model(&samples).Exec(ctx)
				Expect(err).To(BeNil())
			})
			Describe("Retrieving all samples", func() {
				BeforeEach(func() {
					samples = []*models.ChannelSample{
						{
							ChannelConfigID: series.ID,
							Value:           47.3,
							Timestamp:       telem.NewTimeStamp(time.Now()),
						},
						{
							ChannelConfigID: series.ID,
							Timestamp:       telem.NewTimeStamp(time.Now().Add(1 * time.Second)),
						},
						{
							ChannelConfigID: series.ID,
							Timestamp:       telem.NewTimeStamp(time.Now().Add(2 * time.Second)),
						},
					}

				})
				It("Should retrieve the correct items", func() {
					var resSamples []*models.ChannelSample
					err := engine.NewTSRetrieve().Model(&resSamples).WherePK(series.ID).AllTime().Exec(ctx)
					Expect(err).To(BeNil())
					Expect(resSamples).To(HaveLen(3))
				})
			})
			Describe("Retrieving samples from multiple pks", func() {
				var seriesTwo *models.ChannelConfig
				BeforeEach(func() {
					series = &models.ChannelConfig{ID: uuid.New()}
					seriesTwo = &models.ChannelConfig{
						Name: "SG_03",
						ID:   uuid.New(),
					}
					err := engine.NewTSCreate().Model(seriesTwo).Exec(ctx)
					Expect(err).To(BeNil())
					samples = []*models.ChannelSample{
						{
							ChannelConfigID: series.ID,
							Value:           47.3,
							Timestamp:       telem.NewTimeStamp(time.Now()),
						},
						{
							ChannelConfigID: seriesTwo.ID,
							Value:           96.7,
							Timestamp:       telem.NewTimeStamp(time.Now().Add(1 * time.Second)),
						},
					}

				})
				It("Should retrieve the correct items", func() {
					var resSamples []*models.ChannelSample
					err := engine.NewTSRetrieve().Model(&resSamples).WherePKs([]uuid.UUID{seriesTwo.ID, series.ID}).AllTime().Exec(ctx)
					Expect(err).To(BeNil())
					Expect(samples).To(HaveLen(2))
				})
			})
			Describe("Acquire samples across a time rng", func() {
				var err error
				BeforeEach(func() {
					samples = []*models.ChannelSample{
						{
							ChannelConfigID: series.ID,
							Timestamp:       telem.NewTimeStamp(time.Now()),
							Value:           1251.3,
						},
						{
							ChannelConfigID: series.ID,
							Timestamp:       telem.NewTimeStamp(time.Now().Add(-12 * time.Second)),
							Value:           432.3,
						},
						{
							ChannelConfigID: series.ID,
							Timestamp:       telem.NewTimeStamp(time.Now().Add(-30 * time.Second)),
							Value:           322.3,
						},
					}

				})
				It("Should retrieve without error", func() {
					var resSamples []*models.ChannelSample
					tr := telem.NewTimeRange(
						telem.NewTimeStamp(time.Now().Add(-15*time.Second)),
						telem.NewTimeStamp(time.Now().Add(3*time.Second)),
					)
					err = engine.NewTSRetrieve().Model(&resSamples).WherePK(series.ID).WhereTimeRange(tr).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(resSamples).To(HaveLen(2))
				})
			})
		})
		Describe("Checking if a series exists", func() {
			BeforeEach(func() { series = &models.ChannelConfig{ID: uuid.New()} })
			Context("The series does not exist", func() {
				It("Should return false", func() {
					err := engine.NewTSRetrieve().Model(series).WherePK(uuid.New()).Exec(ctx)
					Expect(err).ToNot(BeNil())
				})
			})
		})
	})
	Describe("Edge cases + errors", func() {
		BeforeEach(func() {
			samples = []*models.ChannelSample{{
				ChannelConfigID: series.ID,
				Value:           432.1,
				Timestamp:       telem.NewTimeStamp(time.Now()),
			}}
		})
		JustBeforeEach(func() {
			err := engine.NewTSCreate().Model(&samples).Exec(ctx)
			Expect(err).To(BeNil())
		})
		Context("Retrieving a sample", func() {
			s := &models.ChannelSample{}
			Context("No PKC provided", func() {
				It("Should panic", func() {
					Expect(func() { engine.NewTSRetrieve().Model(s).Exec(ctx) }).To(Panic())
				})
			})
			Context("Invalid PKC provided", func() {
				It("Should return the correct storage error", func() {
					err := engine.NewTSRetrieve().WherePK(uuid.New()).Model(s).
						Exec(ctx)
					Expect(err).ToNot(BeNil())
					Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeItemNotFound))
				})
			})
		})
	})
})
