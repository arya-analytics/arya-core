package roach_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BaseEngine", func() {
	Describe("Adapter", func() {
		var a storage.Adapter
		BeforeEach(func() {
			a = dummyEngine.NewAdapter()
		})
		Describe("NewModel Adapter", func() {
			It("Should createQuery a new adapter without error", func() {
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
					e := &mock.MDEngine{}
					ba := e.NewAdapter()
					Expect(dummyEngine.IsAdapter(ba)).To(BeFalse())
				})
			})
		})
	})
})
