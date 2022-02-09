package storage_test

import (
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/storage"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("errChan", func() {
	Context("Error string", func() {
		It("Should return the correct string", func() {
			err := storage.Error{Type: storage.ErrorTypeUnknown, Message: "Unknown Error"}
			Expect(err.Error()).To(Equal("storage: ErrorTypeUnknown - Unknown Error"))
		})
	})
	Context("Error FieldHandlers", func() {
		Context("Converter chain handles error", func() {
			converterNonDefault := func(err error) (storage.ErrorType, bool) {
				return storage.ErrorTypeRelationshipViolation, true
			}
			converterDefault := func(err error) (storage.ErrorType, bool) {
				return storage.ErrorTypeUnknown, true
			}
			handler := storage.NewErrorHandler(converterDefault, converterNonDefault)
			It("Should return a relationship violation error", func() {
				err := handler.Exec(fmt.Errorf("random error"))
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeRelationshipViolation))
			})
		})
		Context("Default handler handles error", func() {
			converterNonDefault := func(err error) (storage.ErrorType, bool) {
				return storage.ErrorTypeUnknown, false
			}
			converterDefault := func(err error) (storage.ErrorType, bool) {
				return storage.ErrorTypeRelationshipViolation, true
			}
			handler := storage.NewErrorHandler(converterDefault, converterNonDefault)
			It("Should return a relationship violation error", func() {
				err := handler.Exec(fmt.Errorf("random error"))
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrorTypeRelationshipViolation))
			})
		})
		Context("Neither handler handles the error", func() {
			converterNonDefault := func(err error) (storage.ErrorType, bool) {
				return storage.ErrorTypeItemNotFound, false
			}
			converterDefault := func(err error) (storage.ErrorType, bool) {
				return storage.ErrorTypeRelationshipViolation, false
			}
			handler := storage.NewErrorHandler(converterDefault, converterNonDefault)
			It("Should return an unknown error", func() {
				err := handler.Exec(fmt.Errorf("random error"))
				sErr := err.(storage.Error)
				Expect(sErr.Type).To(Equal(storage.ErrorTypeUnknown))
				Expect(sErr.Message).To(Equal("storage - unknown error"))
			})
		})
	})
})
