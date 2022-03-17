package minio_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/internal"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Minio Engine", func() {
	Describe("adapter", func() {
		var a internal.Adapter
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
	})
})
