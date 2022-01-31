package redis_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Redis Engine", func() {
	Describe("Adapter", func() {
		var a storage.Adapter
		BeforeEach(func() {
			a = mockEngine.NewAdapter()
		})
		Describe("New Adapter", func() {
			It("Should create a new adapter without error", func() {
				Expect(reflect.TypeOf(a.ID())).To(Equal(reflect.TypeOf(uuid.New())))
			})
		})
		Describe("Is Adapter", func() {
			Context("Adapter is the correct type", func() {
				It("Should return true", func() {
					Expect(mockEngine.IsAdapter(a)).To(BeTrue())
				})
			})
			Context("Adapter is the incorrect type", func() {
				It("Should return false", func() {
					e := &mock.MDEngine{}
					ba := e.NewAdapter()
					Expect(mockEngine.IsAdapter(ba)).To(BeFalse())
				})
			})
		})
	})
	Describe("Catalog", func() {
		Describe("Contains", func() {
			Context("Model in catalog", func() {
				It("Should return true", func() {
					Expect(mockEngine.InCatalog(&storage.ChannelSample{})).To(BeTrue())
				})
			})
			Context("Model not in catalog", func() {
				It("Should return false", func() {
					Expect(mockEngine.InCatalog(&mock.ModelB{})).To(BeFalse())
				})
			})
		})
	})
})
