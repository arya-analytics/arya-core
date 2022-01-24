package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("QueryUpdate", func() {
	BeforeEach(createMockModel)
	AfterEach(deleteMockModel)
	Describe("Update an item", func() {
		var err error
		var um *storage.ChannelConfig
		BeforeEach(func() {
			um = &storage.ChannelConfig{
				ID:     mockModel.ID,
				Name:   "Cool New Named Name",
				NodeID: 1,
			}
			err = mockEngine.NewUpdate(mockADapter).Model(um).WherePK(mockModel.
				ID).Exec(mockCtx)
			Expect(err).To(BeNil())
		})
		It("Should update it without error", func() {
			Expect(err).To(BeNil())
		})
		It("Should reflect updates when retrieved", func() {
			m := &storage.ChannelConfig{}
			if err := mockEngine.NewRetrieve(mockADapter).Model(m).WherePK(mockModel.
				ID).Exec(mockCtx); err != nil {
				log.Fatalln(err)
			}
			Expect(m.ID).To(Equal(um.ID))
			Expect(m.Name).To(Equal(um.Name))
		})
	})

})
