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

})
