package timeseries_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/redis/timeseries"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("Timeseries", func() {
	BeforeEach(createMockTS)
	AfterEach(deleteMockTS)
	Describe("Create Key", func() {
		It("Should create the Key", func() {
			exists, eErr := mockClient.Exists(mockCtx, mockTSKey).Result()
			Expect(eErr).To(BeNil())
			Expect(exists != 0).To(BeTrue())
		})
	})
	Describe("Add Sample", func() {
		Context("A single sample", func() {
			It("Should add the sample without error", func() {
				err := mockClient.TSCreateSamples(mockCtx, timeseries.Sample{
					Key:       mockTSKey,
					Value:     123.2,
					Timestamp: time.Now().UnixNano(),
				}).Err()
				Expect(err).To(BeNil())
				_, rErr := mockClient.TSGet(mockCtx, mockTSKey).Result()
				Expect(rErr).To(BeNil())
			})
		})
		Context("Multiple samples", func() {
			It("Should add the samples without error", func() {
				err := mockClient.TSCreateSamples(mockCtx, timeseries.Sample{
					Key:       mockTSKey,
					Value:     123.2,
					Timestamp: time.Now().UnixNano(),
				},
					timeseries.Sample{
						Key:       mockTSKey,
						Value:     123.5,
						Timestamp: time.Unix(0, 0).Unix(),
					},
				).Err()
				Expect(err).To(BeNil())
				samples, rErr := mockClient.TSGetAll(mockCtx, mockTSKey).Result()
				Expect(rErr).To(BeNil())
				Expect(samples).To(HaveLen(2))
			})
		})
	})
})
