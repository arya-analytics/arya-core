package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("QueryTsRetrieve", func() {
	var (
		series  *storage.ChannelConfig
		sample  *storage.ChannelSample
		samples []*storage.ChannelSample
	)
	BeforeEach(func() {
		series = &storage.ChannelConfig{ID: uuid.New()}
	})
	JustBeforeEach(func() {
		err := engine.NewTSCreate(adapter).Series().Model(series).Exec(ctx)
		Expect(err).To(BeNil())
	})
	Describe("Standard Usage", func() {
		Context("Single sample", func() {
			JustBeforeEach(func() {
				sampleErr := engine.NewTSCreate(adapter).Sample().Model(sample).Exec(ctx)
				Expect(sampleErr).To(BeNil())
			})
			BeforeEach(func() {
				sample = &storage.ChannelSample{
					ChannelConfigID: series.ID,
					Value:           123.2,
					Timestamp:       time.Now().UnixNano(),
				}
			})
			Describe("Retrieving the most recent sample", func() {
				It("Should retrieve the correct item", func() {
					var resSample = &storage.ChannelSample{}
					err := engine.NewTSRetrieve(adapter).Model(resSample).WherePK(
						series.ID).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(sample).To(Equal(resSample))
				})
			})
		})
		Context("Multiple Samples", func() {
			JustBeforeEach(func() {
				err := engine.NewTSCreate(adapter).Sample().Model(&samples).Exec(ctx)
				Expect(err).To(BeNil())
			})
			Describe("Retrieving all samples", func() {
				BeforeEach(func() {
					samples = []*storage.ChannelSample{
						{
							ChannelConfigID: series.ID,
							Value:           47.3,
							Timestamp:       time.Now().UnixNano(),
						},
						{
							ChannelConfigID: series.ID,
							Timestamp:       time.Now().Add(1 * time.Second).UnixNano(),
						},
						{
							ChannelConfigID: series.ID,
							Timestamp:       time.Now().Add(2 * time.Second).UnixNano(),
						},
					}

				})
				It("Should retrieve the correct items", func() {
					var resSamples []*storage.ChannelSample
					err := engine.NewTSRetrieve(adapter).Model(&resSamples).WherePK(
						series.ID).AllTimeRange().Exec(ctx)
					Expect(err).To(BeNil())
					Expect(resSamples).To(HaveLen(3))
				})
			})
			Describe("Retrieving samples from multiple pks", func() {
				var seriesTwo *storage.ChannelConfig
				BeforeEach(func() {
					series = &storage.ChannelConfig{ID: uuid.New()}
					seriesTwo = &storage.ChannelConfig{
						Name: "SG_03",
						ID:   uuid.New(),
					}
					err := engine.NewTSCreate(adapter).Series().Model(seriesTwo).Exec(ctx)
					Expect(err).To(BeNil())
					samples = []*storage.ChannelSample{
						{
							ChannelConfigID: series.ID,
							Value:           47.3,
							Timestamp:       time.Now().UnixNano(),
						},
						{
							ChannelConfigID: seriesTwo.ID,
							Value:           96.7,
							Timestamp:       time.Now().Add(1 * time.Second).UnixNano(),
						},
					}

				})
				It("Should retrieve the correct items", func() {
					var resSamples []*storage.ChannelSample
					err := engine.NewTSRetrieve(adapter).Model(&resSamples).WherePKs(
						[]uuid.UUID{seriesTwo.ID, series.ID}).AllTimeRange().Exec(ctx)
					Expect(err).To(BeNil())
					Expect(samples).To(HaveLen(2))
				})
			})
			Describe("Retrieve samples across a time range", func() {
				var err error
				BeforeEach(func() {
					samples = []*storage.ChannelSample{
						{
							ChannelConfigID: series.ID,
							Timestamp:       time.Now().UnixNano(),
							Value:           1251.3,
						},
						{
							ChannelConfigID: series.ID,
							Timestamp:       time.Now().Add(-12 * time.Second).UnixNano(),
							Value:           432.3,
						},
						{
							ChannelConfigID: series.ID,
							Timestamp:       time.Now().Add(-30 * time.Second).UnixNano(),
							Value:           322.3,
						},
					}

				})
				It("Should retrieve without error", func() {
					var resSamples []*storage.ChannelSample
					toTS := time.Now().Add(3 * time.Second).UnixNano()
					fromTS := time.Now().Add(-15 * time.Second).UnixNano()
					err = engine.NewTSRetrieve(adapter).Model(&resSamples).WherePK(
						series.ID).WhereTimeRange(fromTS, toTS).Exec(ctx)
					Expect(err).To(BeNil())
					Expect(resSamples).To(HaveLen(2))
				})
			})
		})
		Describe("Checking if a series exists", func() {
			BeforeEach(func() { series = &storage.ChannelConfig{ID: uuid.New()} })
			Context("The series does not exist", func() {
				It("Should return false", func() {
					e, err := engine.NewTSRetrieve(adapter).SeriesExists(ctx, uuid.New())
					Expect(e).To(BeFalse())
					Expect(err).To(BeNil())
				})
			})
		})
	})
	Describe("Edge cases + errors", func() {
		BeforeEach(func() {
			samples = []*storage.ChannelSample{{
				ChannelConfigID: series.ID,
				Value:           432.1,
				Timestamp:       time.Now().UnixNano(),
			}}
		})
		JustBeforeEach(func() {
			err := engine.NewTSCreate(adapter).Sample().Model(&samples).Exec(
				ctx)
			Expect(err).To(BeNil())
		})
		Context("Retrieving a sample", func() {
			s := &storage.ChannelSample{}
			Context("No PK provided", func() {
				It("Should return the correct storage error", func() {
					err := engine.NewTSRetrieve(adapter).Model(s).Exec(ctx)
					Expect(err).ToNot(BeNil())
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeInvalidArgs))
				})
			})
			Context("Invalid PK provided", func() {
				It("Should return the correct storage error", func() {
					err := engine.NewTSRetrieve(adapter).WherePK(uuid.New()).Model(s).
						Exec(ctx)
					Expect(err).ToNot(BeNil())
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeItemNotFound))
				})
			})
		})
	})
})
