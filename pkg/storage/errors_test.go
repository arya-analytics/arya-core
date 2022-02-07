package storage_test

import (
	"fmt"
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
	Context("Error Handler", func() {
		Context("Converter chain handles error", func() {
			converterChain := storage.ErrorTypeConverterChain{
				func(err error) (storage.ErrorType, bool) {
					return storage.ErrTypeRelationshipViolation, true
				},
			}
			converterDefault := func(err error) (storage.ErrorType, bool) {
				return storage.ErrTypeUnknown, true
			}
			handler := storage.NewErrorHandler(converterChain, converterDefault)
			It("Should return a relationship violation error", func() {
				err := handler.Exec(fmt.Errorf("random error"))
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeRelationshipViolation))
			})
		})
		Context("Default handler handles error", func() {
			converterChain := storage.ErrorTypeConverterChain{
				func(err error) (storage.ErrorType, bool) {
					return storage.ErrTypeUnknown, false
				},
			}
			converterDefault := func(err error) (storage.ErrorType, bool) {
				return storage.ErrTypeRelationshipViolation, true
			}
			handler := storage.NewErrorHandler(converterChain, converterDefault)
			It("Should return a relationship violation error", func() {
				err := handler.Exec(fmt.Errorf("random error"))
				Expect(err.(storage.Error).Type).To(Equal(storage.ErrTypeRelationshipViolation))
			})
		})
		Context("Neither handler handles the error", func() {
			converterChain := storage.ErrorTypeConverterChain{
				func(err error) (storage.ErrorType, bool) {
					return storage.ErrTypeItemNotFound, false
				},
			}
			converterDefault := func(err error) (storage.ErrorType, bool) {
				return storage.ErrTypeRelationshipViolation, false
			}
			handler := storage.NewErrorHandler(converterChain, converterDefault)
			It("Should return an unknown error", func() {
				err := handler.Exec(fmt.Errorf("random error"))
				sErr := err.(storage.Error)
				Expect(sErr.Type).To(Equal(storage.ErrTypeUnknown))
				Expect(sErr.Message).To(Equal("storage - unknown error"))

			})
		})
	})
})
