package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("retrieveQuery", func() {
	Describe("Standard usage", func() {
		Context("Meta Data Only", func() {
			BeforeEach(createMockChannelCfg)
			AfterEach(deleteMockChannelCfg)
			Context("Single item", func() {
				Describe("Retrieve a channel config", func() {
					It("Should retrieve without errutil", func() {
						m := &storage.ChannelConfig{}
						err := mockStorage.NewRetrieve().Model(m).WherePK(mockChannelCfg.ID).Exec(mockCtx)
						Expect(err).To(BeNil())
					})
					It("Should retrieve the correct item", func() {
						m := &storage.ChannelConfig{}
						err := mockStorage.NewRetrieve().Model(m).WherePK(mockChannelCfg.ID).Exec(mockCtx)
						Expect(err).To(BeNil())
						Expect(m.ID).To(Equal(mockChannelCfg.ID))
						Expect(m.Name).To(Equal(mockChannelCfg.Name))
					})
				})
			})
		})
		Context("Object Data + Meta Data", func() {
			Context("Single item", func() {
				BeforeEach(createMockChannelChunk)
				AfterEach(deleteMockChannelChunk)
				Describe("Retrieve a channel chunk", func() {
					var retrievedModel = &storage.ChannelChunk{}
					var err error
					BeforeEach(func() {
						err = mockStorage.NewRetrieve().Model(retrievedModel).WherePK(
							mockChannelChunk.ID).Exec(mockCtx)
					})
					It("Should retrieve it without errutil", func() {
						Expect(err).To(BeNil())
					})
					It("Should retrieve the correct item", func() {
						Expect(retrievedModel.ID).To(Equal(mockChannelChunk.ID))
						Expect(retrievedModel.Data).ToNot(BeNil())
						b := make([]byte, retrievedModel.Data.Size())
						_, err = retrievedModel.Data.Read(b)
						Expect(err.Error()).To(Equal("EOF"))
						Expect(b).To(Equal(mockBytes))
					})
				})
			})
		})
	})
})
