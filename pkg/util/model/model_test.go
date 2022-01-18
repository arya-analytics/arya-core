package model_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"reflect"
)

var _ = Describe("Model", func() {
	Describe("Reflect", func() {
		Context("Single Model", func() {
			var m = &mock.ModelA{
				ID: 22,
			}
			var mBaseType = reflect.TypeOf(mock.ModelA{})
			var mType = reflect.TypeOf(m)
			var refl = model.NewReflect(m)
			It("Should pass validation without err", func() {
				Expect(refl.Validate()).To(BeNil())
			})
			It("Should return the correct pointer interface", func() {
				Expect(refl.Pointer()).To(Equal(m))
			})
			It("Should return the correct type", func() {
				Expect(refl.Type()).To(Equal(mBaseType))
			})
			It("Should return the correct value", func() {
				Expect(refl.Value().Type()).To(Equal(mBaseType))
			})
			It("Should return false for IsChain", func() {
				Expect(refl.IsChain()).To(BeFalse())
			})
			It("Should return true for IsStruct", func() {
				Expect(refl.IsStruct()).To(BeTrue())
			})
			It("Should return the correct struct field by name", func() {
				Expect(refl.Value().FieldByName("ID").Interface()).To(Equal(22))
			})
			It("Should return the correct struct field by index", func() {
				Expect(refl.Value().Field(0).Interface()).To(Equal(22))
			})
			It("Should return the same item for the raw value as for the value",
				func() {
					Expect(refl.RawValue()).To(Equal(refl.Value()))
				})
			It("Should return the same type for the raw type as for the type", func() {
				Expect(refl.RawType()).To(Equal(refl.Type()))
			})
			It("Should return the correct pointer type", func() {
				Expect(refl.PointerType()).To(Equal(mType))
			})
			It("Should return the correct pointer value", func() {
				Expect(refl.PointerValue()).To(Equal(reflect.ValueOf(m)))
			})
			Describe("New Chain", func() {
				It("Should return the correct type", func() {
					newC := refl.NewChain()
					Expect(newC.RawType()).To(Equal(reflect.TypeOf([]*mock.ModelA{})))
					Expect(newC.Type()).To(Equal(mBaseType))
				})
			})
			Describe("New Model", func() {
				It("Should return the correct type", func() {
					newM := refl.NewModel()
					Expect(newM.PointerType()).To(Equal(mType))
					Expect(newM.Type()).To(Equal(mBaseType))
				})
			})
			Describe("New Raw", func() {
				It("Should return a single model", func() {
					newM := refl.NewRaw()
					Expect(newM.PointerType()).To(Equal(mType))
					Expect(newM.Type()).To(Equal(mBaseType))
				})
			})
		})
		Context("Multiple Models", func() {
			var m = []*mock.ModelA{
				&mock.ModelA{
					ID: 22,
				},
				&mock.ModelA{
					ID: 43,
				},
			}
			var mBaseType = reflect.TypeOf(m)
			var mType = reflect.TypeOf(&m)
			var mSingleBaseType = reflect.TypeOf(mock.ModelA{})
			var mSingleType = reflect.TypeOf(&mock.ModelA{})
			var refl = model.NewReflect(&m)
			It("Should pass validation without error", func() {
				Expect(refl.Validate()).To(BeNil())
			})
			It("Should return the correct pointer interface", func() {
				Expect(refl.Pointer()).To(Equal(&m))
			})
			It("Should return the correct type", func() {
				Expect(refl.Type()).To(Equal(mSingleBaseType))
			})
			It("Should return true for IsChain", func() {
				Expect(refl.IsChain()).To(BeTrue())
			})
			It("Should return false for IsStruct", func() {
				Expect(refl.IsStruct()).To(BeFalse())
			})
			It("Should return the correct model value by index", func() {
				Expect(refl.ChainValueByIndex(0).PointerType()).To(Equal(mSingleType))
				Expect(refl.ChainValueByIndex(0).Type()).To(Equal(mSingleBaseType))
				Expect(refl.ChainValueByIndex(0).Pointer()).To(Equal(m[0]))
			})
			It("Should return a slice for the raw value", func() {
				Expect(refl.RawValue().Interface()).To(Equal(m))
				Expect(refl.RawType()).To(Equal(mBaseType))
			})
			It("Should return the correct pointer type", func() {
				Expect(refl.PointerType()).To(Equal(mType))
			})
			It("Should return the correct pointer value", func() {
				Expect(refl.PointerValue()).To(Equal(reflect.ValueOf(&m)))
			})
			Describe("Appending to a chain", func() {
				It("Should append to the chain correctly", func() {
					mTwo := &mock.ModelA{ID: 1}
					refl.ChainAppend(model.NewReflect(mTwo))
					Expect(refl.ChainValueByIndex(2).Pointer()).To(Equal(mTwo))
				})
			})
			Describe("New Chain", func() {
				It("Should return the correct type", func() {
					newC := refl.NewChain()
					Expect(newC.RawType()).To(Equal(mBaseType))
					Expect(newC.Type()).To(Equal(mSingleBaseType))
				})
			})
			Describe("New Model", func() {
				It("Should return the correct type", func() {
					newM := refl.NewModel()
					Expect(newM.PointerType()).To(Equal(mSingleType))
					Expect(newM.Type()).To(Equal(mSingleBaseType))
				})
			})
			Describe("New Raw", func() {
				It("Should return a model chain", func() {
					newM := refl.NewRaw()
					Expect(newM.PointerType()).To(Equal(mType))
					Expect(newM.RawType()).To(Equal(mBaseType))
					Expect(newM.Type()).To(Equal(mSingleBaseType))
				})
			})
		})
		Describe("Errors + edge cases", func() {
			It("Should return an error when a non pointer is provided", func() {
				Expect(model.NewReflect(mock.ModelA{ID: 22}).Validate()).ToNot(BeNil())
			})
			It("Should return an error when a non struct is provided", func() {
				i := 11
				Expect(model.NewReflect(&i).Validate()).ToNot(BeNil())
			})
			It("Should return an error when initializing with a nil pointer", func() {
				refl := model.NewReflect((*mock.ModelA)(nil))
				Expect(refl.Validate()).To(BeNil())
				log.Info(refl.NewModel().Pointer())
				Expect(refl.NewModel())
			})
			It("Should return an error when initializing with a nil pointer", func() {
				var m []*mock.ModelA
				refl := model.NewReflect(&m)
				Expect(refl.Validate()).To(BeNil())
				log.Info(refl.NewModel().Pointer())
				Expect(refl.NewModel())
				log.Info(refl.IsChain())
			})
		})
	})

})
