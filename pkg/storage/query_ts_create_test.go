package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"time"
)

var _ = Describe("QueryTsCreate", func() {
	BeforeEach(createMockChannelCfg)
	AfterEach(deleteMockChannelCfg)
	Describe("Standard Usage", func() {
		Describe("Create a new series", func() {
			var err error
			BeforeEach(func() {
				err = mockStorage.NewTSCreate().Series().Model(mockChannelCfg).Exec(
					mockCtx)
			})
			It("Should create the series without error", func() {
				Expect(err).To(BeNil())
			})
			It("Should exist after creation", func() {
				exists, rErr := mockStorage.NewTSRetrieve().SeriesExists(mockCtx,
					mockChannelCfg.ID)
				Expect(rErr).To(BeNil())
				Expect(exists).To(BeTrue())
			})
		})
		Describe("Create a new sample", func() {
			BeforeEach(createMockSeries)
			Context("Single Sample", func() {
				var err error
				var sample *storage.ChannelSample
				BeforeEach(func() {
					sample = &storage.ChannelSample{
						ChannelConfigID: mockChannelCfg.ID,
						Value:           32.7,
						Timestamp:       time.Now().UnixNano(),
					}
					err = mockStorage.NewTSCreate().Sample().Model(sample).Exec(mockCtx)
				})
				It("Should create the sample without error", func() {
					Expect(err).To(BeNil())
				})
				It("Should exist after creation", func() {
					resSample := &storage.ChannelSample{}
					rErr := mockStorage.NewTSRetrieve().Model(resSample).WherePK(
						mockChannelCfg.ID).Exec(mockCtx)
					Expect(rErr).To(BeNil())
					Expect(resSample.Timestamp).To(Equal(sample.Timestamp))
				})
			})
			Context("Multiple Samples", func() {
				var qty = 4
				var err error
				var samples []*storage.ChannelSample
				var channelCfgChain []*storage.ChannelConfig
				log.SetReportCaller(true)
				BeforeEach(func() {
					samples = []*storage.ChannelSample{}
					channelCfgChain = []*storage.ChannelConfig{}
					for i := 0; i < qty; i++ {
						mockCC := &storage.ChannelConfig{
							ID:     uuid.New(),
							Name:   "SG 43",
							NodeID: mockNode.ID,
						}
						channelCfgChain = append(channelCfgChain, mockCC)
						if err := mockStorage.NewCreate().Model(mockCC).Exec(
							mockCtx); err != nil {
							log.Fatalln(err)
						}
						if err := mockStorage.NewTSCreate().Series().Model(mockCC).Exec(
							mockCtx); err != nil {
							log.Fatalln(err)
						}
						samples = append(samples,
							&storage.ChannelSample{ChannelConfigID: mockCC.ID, Value: 126.8,
								Timestamp: time.Now().Add(1 * time.Second).UnixNano(),
							})
					}
					err = mockStorage.NewTSCreate().Sample().Model(&samples).Exec(
						mockCtx)
				})
				It("Should create the samples without error", func() {
					Expect(err).To(BeNil())
				})
				It("The samples should be able to be retrieved after creation", func() {
					var resSamples []*storage.ChannelSample
					var pks = model.NewReflect(&channelCfgChain).PKChain().Interface()
					pks = append(pks, mockChannelCfg.ID)
					rErr := mockStorage.NewTSRetrieve().Model(&resSamples).WherePKs(
						pks).AllTimeRange().Exec(mockCtx)
					Expect(rErr).To(BeNil())
					Expect(resSamples).To(HaveLen(4))
				})
			})
		})
	})
})
