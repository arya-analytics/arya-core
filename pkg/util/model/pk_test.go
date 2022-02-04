package model_test

import (
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strconv"

	"github.com/arya-analytics/aryacore/pkg/util/model"
)

var _ = Describe("PK", func() {
	Describe("Single PK", func() {
		Describe("Stringifying", func() {
			It("Should return a UUID as a string", func() {
				id := uuid.New()
				Expect(model.NewPK(id).String()).To(Equal(id.String()))
			})
			It("Should return an int as a string", func() {
				i := 1
				Expect(model.NewPK(i).String()).To(Equal(strconv.Itoa(int(i))))
			})
			It("Should return an int32 as a string", func() {
				var id32 int32 = 1
				Expect(model.NewPK(id32).String()).To(Equal(strconv.Itoa(int(id32))))
			})
			It("Should return an int64 as a string", func() {
				var id64 int64 = 1
				Expect(model.NewPK(id64).String()).To(Equal(strconv.Itoa(int(id64))))
			})
			It("Should return a string as a string", func() {
				s := "Hello"
				Expect(model.NewPK(s).String()).To(Equal(s))
			})
			It("Should panic with an unknown pk type", func() {
				Expect(func() {
					_ = model.NewPK(123.2).String()
				}).To(Panic())
			})
		})
		Describe("Equality Check", func() {
			It("Should return true when two UUIDs are equal", func() {
				id := uuid.New()
				Expect(model.NewPK(id).Equals(model.NewPK(id)))
			})
		})
		Describe("Reflect StructValue", func() {
			It("Should return the correct reflect value", func() {
				id := uuid.New()
				Expect(model.NewPK(id).Value().Interface()).To(Equal(id))
			})
		})
		Describe("Is Zero", func() {
			It("Should return true when the id is a zero value", func() {
				var id int
				Expect(model.NewPK(id).IsZero()).To(BeTrue())
			})
		})
	})
	Describe("Multiple PKS", func() {
		Describe("Standard usage", func() {
			rawPks := []uuid.UUID{uuid.New(), uuid.New()}
			pks := model.NewPKChain(rawPks)
			It("Should have the correct length", func() {
				Expect(pks).To(HaveLen(2))
			})
			It("Should return the correct interface value", func() {
				Expect(pks.Raw()[0]).To(Equal(rawPks[0]))
				Expect(pks.Raw()[1]).To(Equal(rawPks[1]))
			})
		})
		Describe("Edge cases + errors", func() {
			It("Should panic when a non-slice is provided", func() {
				Expect(func() {
					model.NewPKChain(uuid.New())
				}).To(Panic())
			})
		})

	})
})
