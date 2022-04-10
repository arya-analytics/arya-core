package model_test

import (
	"fmt"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"strconv"

	"github.com/arya-analytics/aryacore/pkg/util/model"
)

var _ = Describe("PKC", func() {
	Describe("Single PKC", func() {
		var uuidPk = uuid.New()
		Describe("Stringifying", func() {
			DescribeTable("Standard Usage", func(pk interface{}, expected interface{}) {
				Expect(model.NewPK(pk).String()).To(Equal(expected))
			},
				Entry("UUID", uuidPk, uuidPk.String()),
				Entry("int", 1, strconv.Itoa(1)),
				Entry("int32", int32(2), strconv.Itoa(2)),
				Entry("int64", int64(2), strconv.Itoa(2)),
				Entry("string", "hello", "hello"),
			)
			It("Should panic with an unknown pk type", func() {
				Expect(func() {
					_ = model.NewPK(123.2).String()
				}).To(Panic())
			})
		})
		Describe("Setting from a string", func() {
			DescribeTable("Normal usage",
				func(rawSource interface{}, rawDest interface{}) {
					sourcePK, destPK := model.NewPK(rawSource), model.NewPK(rawDest)
					newPK, err := destPK.NewFromString(sourcePK.String())
					Expect(err).To(BeNil())
					Expect(newPK.Type()).To(Equal(destPK.Type()))
					Expect(newPK.String()).To(Equal(sourcePK.String()))
				},
				Entry("Set UUID from string", uuid.New().String(), uuid.UUID{}),
				Entry("Set string from UUID", uuid.New(), ""),
				Entry("Set int from string", "123", 0),
				Entry("Set string from int", 123, ""),
				Entry("Set int32 from string", "123", int32(0)),
				Entry("Set string from int32", int32(123), ""),
				Entry("Set int64 from string", "123", int64(0)),
				Entry("Set string from int64", int64(123), ""),
			)
			It("Should panic when an unknown type is provided", func() {
				pk := model.NewPK(123.2)
				Expect(func() {
					newPK, _ := pk.NewFromString("123")
					// Just here so we don't get a compiler error
					fmt.Println(newPK)
				}).To(Panic())
			})
			Describe("Creating a PK chain from a string", func() {
				It("Should create a chain of PKs from a string", func() {
					pkcStr := []string{uuid.New().String(), uuid.New().String()}
					newPKC, err := model.NewPK(uuid.UUID{}).NewChainFromStrings(pkcStr...)
					Expect(err).To(BeNil())
					Expect(newPKC).To(HaveLen(len(pkcStr)))
				})
			})
		})
		Describe("Equality Check", func() {
			It("Should return true when two UUIDs are equal", func() {
				id := uuid.New()
				Expect(model.NewPK(id).Equals(model.NewPK(id)))
			})
		})
		Describe("Reflect StructValue", func() {
			It("Should return the correct reflect Val", func() {
				id := uuid.New()
				Expect(model.NewPK(id).Value().Interface()).To(Equal(id))
			})
		})
		Describe("Is Zero", func() {
			It("Should return true when the id is a zero Val", func() {
				var id int
				Expect(model.NewPK(id).IsZero()).To(BeTrue())
			})
		})
	})
	Describe("PK Chain", func() {
		Describe("Standard usage", func() {
			rawPks := []uuid.UUID{uuid.New(), uuid.New()}
			pks := model.NewPKChain(rawPks)
			It("Should have the correct length", func() {
				Expect(pks).To(HaveLen(2))
			})
			It("Should return the correct interface Val", func() {
				Expect(pks.Raw()[0]).To(Equal(rawPks[0]))
				Expect(pks.Raw()[1]).To(Equal(rawPks[1]))
			})
			It("Should return the correct PKS as strings", func() {
				Expect(pks.Strings()[0]).To(Equal(rawPks[0].String()))
				Expect(pks.Strings()[1]).To(Equal(rawPks[1].String()))
			})
			Context("All Zero", func() {
				It("Should return false when a pk is not zero", func() {
					Expect(pks.AllZero()).To(BeFalse())
				})
				It("Should return true when all pks are zero", func() {
					newPKC := model.NewPKChain([]uuid.UUID{{}, {}})
					Expect(newPKC.AllZero()).To(BeTrue())
				})
				It("Should return true when the chain is empty", func() {
					newPKC := model.NewPKChain([]uuid.UUID{})
					Expect(newPKC.AllZero()).To(BeTrue())
				})
			})
			Context("AllNonZero", func() {
				It("Should return true when all pks are zero", func() {
					Expect(pks.AllNonZero()).To(BeTrue())
				})
				It("Should return false when one of the pks is non zero", func() {
					newPKC := model.NewPKChain([]uuid.UUID{uuid.New(), {}})
					Expect(newPKC.AllNonZero()).To(BeFalse())
				})
			})
			Context("Contains", func() {
				It("Should return true when the chain contains the PK", func() {
					id1 := uuid.New()
					id2 := uuid.New()
					newPKC := model.NewPKChain([]uuid.UUID{id1, id2})
					Expect(newPKC.Contains(model.NewPK(id1))).To(BeTrue())
				})
				It("Should return false when the chain doesn't contain the PK", func() {
					id1 := uuid.New()
					id2 := uuid.New()
					newPKC := model.NewPKChain([]uuid.UUID{id1, id2})
					Expect(newPKC.Contains(model.NewPK(uuid.New()))).To(BeFalse())
				})
			})
			Context("Unique", func() {
				It("Should filter out duplicates", func() {
					id1 := uuid.New()
					id2 := uuid.New()
					newPKC := model.NewPKChain([]uuid.UUID{id1, id2, id2})
					Expect(newPKC.Unique()).To(HaveLen(2))
				})

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
