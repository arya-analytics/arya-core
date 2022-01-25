package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hook", func() {
	Describe("UUID auto-generation", func() {
		var (
			cc  *storage.ChannelConfig
			err error
		)
		BeforeEach(func() {
			cc = &storage.ChannelConfig{
				Name:   "Auto-generated UUID",
				NodeID: mockNode.ID,
			}
			err = mockEngine.NewCreate(mockAdapter).Model(cc).Exec(mockCtx)
		})
		It("Should generate it without err", func() {
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			var retrievedCC = &storage.ChannelConfig{}
			err := mockEngine.NewRetrieve(mockAdapter).
				Model(retrievedCC).
				Where("NAME = ?", cc.Name).
				Exec(mockCtx)
			Expect(err).To(BeNil())
			Expect(retrievedCC.Name).To(Equal(cc.Name))
		})
	})
})
