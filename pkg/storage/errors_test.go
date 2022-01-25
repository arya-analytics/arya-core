package storage_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Errors", func() {
	Context("Error string", func() {
		It("Should return the correct string", func() {
			err := storage.Error{Type: storage.ErrTypeUnknown, Message: "Unknown Error"}
			Expect(err.Error()).To(Equal("storage: ErrTypeUnknown - Unknown Error"))
		})
	})
})
