package model_test

import (
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/arya-analytics/aryacore/pkg/util/model/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Fields", func() {
	Describe("Creating a new reflect from fields", func() {
		Context("Inner model is not nil", func() {
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
			Specify("Changes to the field should reflect in teh original model", func() {
				baseRfl := model.NewReflect(m)
				newRfl := baseRfl.FieldsByName("InnerModel").ToReflect()
				newRfl.ChainValueByIndex(0).StructFieldByName("ID").Set(reflect.ValueOf(98))
				Expect(m.InnerModel.ID).To(Equal(98))
			})
		})
		Context("Inner model is nil", func() {
			var m *mock.ModelA
			BeforeEach(func() {
				m = &mock.ModelA{}
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
				Expect(newRfl.ChainValueByIndex(0).PK().Raw()).To(Equal(0))
			})
			Specify("Changes to the field should reflect in teh original model", func() {
				baseRfl := model.NewReflect(m)
				newRfl := baseRfl.FieldsByName("InnerModel").ToReflect()
				newRfl.ChainValueByIndex(0).StructFieldByName("ID").Set(reflect.ValueOf(98))
				Expect(m.InnerModel.ID).To(Equal(98))
			})
		})
	})
})
