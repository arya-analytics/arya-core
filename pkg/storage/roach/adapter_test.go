package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/stub"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)



var _ = Describe("Engine", func() {
	Describe("Adapter", func() {
		var a storage.Adapter
		BeforeEach(func() {
			a = dummyEngine.NewAdapter()
		})
		Describe("New Adapter", func() {
			It("Should create a new adapter without error", func() {
				Expect(len(a.ID().String())).To(Equal(len(uuid.New().String())))
			})
		})
		Describe("Is Adapter", func() {
			Context("Adapter is the correct type", func() {
				It("Should return true", func() {
					Expect(dummyEngine.IsAdapter(a)).To(BeTrue())
				})
			})
			Context("Adapter is the incorrect type", func() {
				It("Should return false", func() {
					e := &stub.MDEngine{}
					ba := e.NewAdapter()
					Expect(dummyEngine.IsAdapter(ba)).To(BeFalse())
				})
			})
		})
	})
})
