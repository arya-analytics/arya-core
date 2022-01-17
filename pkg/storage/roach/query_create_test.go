package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

var _ = Describe("Create", func() {
	BeforeEach(migrate)
	AfterEach(deleteDummyModel)
	Describe("Create a new Channel Config", func() {
		It("Should createQuery it without error", func() {
			err := dummyEngine.NewCreate(dummyAdapter).Model(dummyModel).Exec(dummyCtx)
			Expect(err).To(BeNil())
		})
		It("Should be able to be re-queried after creation", func() {
			if err := dummyEngine.NewCreate(dummyAdapter).Model(dummyModel).Exec(
				dummyCtx); err != nil {
				log.Fatalln(err)
			}
			m := &storage.ChannelConfig{}
			err := dummyEngine.NewRetrieve(dummyAdapter).Model(m).WherePK(dummyModel.
				ID).Exec(dummyCtx)
			Expect(err).To(BeNil())
			Expect(m.Name).To(Equal(dummyModel.Name))
			Expect(m.ID).To(Equal(dummyModel.ID))
		})
	})
	Describe("Edge cases + errors", func() {
		Context("No PK provided", func() {
			Context("Single Model", func() {
				It("Should return the correct error type", func() {
					m := &storage.ChannelConfig{
						Name: "Hello",
					}
					err := dummyEngine.NewCreate(dummyAdapter).Model(m).Exec(dummyCtx)
					Expect(err).ToNot(BeNil())
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeNoPK))
				})
			})
			Context("Chain of models", func() {
				It("Should return the correct error type", func() {
					m := []*storage.ChannelConfig{
						&storage.ChannelConfig{
							ID: 12, Name: "Hello",
						},
						&storage.ChannelConfig{},
					}
					err := dummyEngine.NewCreate(dummyAdapter).Model(&m).Exec(dummyCtx)
					Expect(err).ToNot(BeNil())
					Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeNoPK))
				})
			})
		})
	})
})
