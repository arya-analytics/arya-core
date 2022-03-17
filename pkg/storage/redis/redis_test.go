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
	Describe("adapter", func() {
		var a storage.Adapter
		BeforeEach(func() {
			var err error
			a, err = engine.NewAdapter()
			Expect(err).To(BeNil())
		})
		Describe("New adapter", func() {
			It("Should create a new adapter without error", func() {
				Expect(reflect.TypeOf(a.ID())).To(Equal(reflect.TypeOf(uuid.New())))
			})
		})
		Describe("Is adapter", func() {
			Context("adapter is the correct type", func() {
				It("Should return true", func() {
					Expect(engine.IsAdapter(a)).To(BeTrue())
				})
			})
			Context("adapter is the incorrect type", func() {
				It("Should return false", func() {
					e := &mock.MDEngine{}
					ba := e.NewAdapter()
					Expect(engine.IsAdapter(ba)).To(BeFalse())
				})
			})
		})
	})
})
