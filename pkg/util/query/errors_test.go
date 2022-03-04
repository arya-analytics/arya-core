package query_test

import (
	"errors"
	"fmt"
	"github.com/arya-analytics/aryacore/pkg/util/query"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("errChan", func() {
	Context("Error string", func() {
		It("Should return the correct string", func() {
			err := query.Error{Type: query.ErrorTypeUnknown, Message: "Unknown Error", Base: errors.New("unknown error")}
			Expect(err.Error()).To(Equal("query: ErrorTypeUnknown - Unknown Error - unknown error"))
		})
	})
	Context("Error FieldHandlers", func() {
		Context("Converter chain handles error", func() {
			converterNonDefault := func(err error) (error, bool) {
				return query.NewSimpleError(query.ErrorTypeRelationshipViolation, err), true
			}
			handler := query.NewErrorConvertChain(converterNonDefault)
			It("Should return a relationship violation error", func() {
				err := handler.Exec(fmt.Errorf("random error"))
				Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeRelationshipViolation))
			})
		})
		Context("Default handler handles error", func() {
			converterNonDefault := func(err error) (error, bool) {
				return query.NewSimpleError(query.ErrorTypeUnknown, err), false
			}
			converterDefault := func(err error) (error, bool) {
				return query.NewSimpleError(query.ErrorTypeRelationshipViolation, err), true
			}
			handler := query.NewErrorConvertChain(converterDefault, converterNonDefault)
			It("Should return a relationship violation error", func() {
				err := handler.Exec(fmt.Errorf("random error"))
				Expect(err.(query.Error).Type).To(Equal(query.ErrorTypeRelationshipViolation))
			})
		})
		Context("Neither handler handles the error", func() {
			converterNonDefault := func(err error) (error, bool) {
				return query.NewSimpleError(query.ErrorTypeItemNotFound, err), false
			}
			converterDefault := func(err error) (error, bool) {
				return query.NewSimpleError(query.ErrorTypeRelationshipViolation, err), false
			}
			handler := query.NewErrorConvertChain(converterDefault, converterNonDefault)
			It("Should return an unknown error", func() {
				err := handler.Exec(fmt.Errorf("random error"))
				sErr := err.(query.Error)
				Expect(sErr.Type).To(Equal(query.ErrorTypeUnknown))
				Expect(sErr.Message).To(Equal("query -> unknown error"))
			})
		})
	})
})
