package model_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Fields", func() {
	Describe("Creating a new reflect from fields", func() {
		var m *mock.ModelA
		BeforeEach(func() {
			m = &mock.ModelA{InnerModel: &mock.ModelB{ID: 96}}
		})
		It("Should create a new reflect ", func() {
			baseRfl := model.NewReflect(m)
			Expect(func() {
				baseRfl.FieldsByName("InnerModel").ToReflect()
			}).ToNot(Panic())
		})
		It("Should create a reflect with the correct items", func() {
			baseRfl := model.NewReflect(m)
			newRfl := baseRfl.FieldsByName("InnerModel").ToReflect()
			Expect(newRfl.IsChain()).To(BeTrue())
			Expect(newRfl.ChainValue().Len()).To(Equal(1))
			Expect(newRfl.ChainValueByIndex(0).PK().Raw()).To(Equal(96))
		})
	})

})
