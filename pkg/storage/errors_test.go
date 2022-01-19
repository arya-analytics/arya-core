package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Errors", func() {
	Context("Error string", func() {
		It("Should return the correct string", func() {
			err := storage.NewError(storage.ErrTypeInvalidField)
			Expect(err.Error()).To(Equal("storage: ErrTypeInvalidField"))
		})
	})
})
