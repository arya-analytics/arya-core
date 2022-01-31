package model_test

import (
	"github.com/arya-analytics/aryacore/pkg/storage/mock"
	"github.com/arya-analytics/aryacore/pkg/util/model"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"reflect"
)

var _ = Describe("Reflect", func() {
	Describe("Pointer Checks", func() {
		It("Should return true when the model is a pointer", func() {
			Expect(model.NewReflect(&mock.ModelA{}).IsPointer()).To(BeTrue())
		})
		It("Should return false when the model is a pointer", func() {
			Expect(model.NewReflect(mock.ModelA{}).IsPointer()).To(BeFalse())
		})
	})
	Describe("Pointer Creation", func() {
		It("Should create a new pointer for a non-pointer model", func() {
			Expect(model.NewReflect(mock.ModelA{}).ToNewPointer().IsPointer()).To(BeTrue())
		})
		It("Should create the pointer to the correct underlying value", func() {
			var baseVal []*mock.ModelA
			baseRfl := model.NewReflect(baseVal)
			Expect(baseRfl.ToNewPointer().RawValue().Interface()).To(Equal(baseVal))
		})
	})
	Context("Single Model", func() {
		var m = &mock.ModelA{
			ID: 22,
		}
		var mBaseType = reflect.TypeOf(mock.ModelA{})
		var mType = reflect.TypeOf(m)
		var refl = model.NewReflect(m)
		It("Should pass validation without panicking", func() {
			Expect(refl.Validate).ToNot(Panic())
		})
		It("Should return the correct pointer interface", func() {
			Expect(refl.Pointer()).To(Equal(m))
		})
		It("Should return the correct type", func() {
			Expect(refl.Type()).To(Equal(mBaseType))
		})
		It("Should return the correct value", func() {
			Expect(refl.StructValue().Type()).To(Equal(mBaseType))
		})
		It("Should return false for IsChain", func() {
			Expect(refl.IsChain()).To(BeFalse())
		})
		It("Should return true for IsStruct", func() {
			Expect(refl.IsStruct()).To(BeTrue())
		})
		It("Should return the correct struct field by name", func() {
			Expect(refl.StructValue().FieldByName("ID").Interface()).To(Equal(22))
		})
		It("Should return the correct struct field by role", func() {
			Expect(refl.StructFieldByRole("pk").Interface()).To(Equal(22))
		})
		It("Should panic if the role doesn't exist", func() {
			Expect(func() {
				refl.StructFieldByRole("nonexistentrole").Interface()
			}).To(Panic())
		})

		It("Should return the correct struct field by index", func() {
			Expect(refl.StructValue().Field(0).Interface()).To(Equal(22))
		})
		It("Should return the same item for the base value as for the value",
			func() {
				Expect(refl.RawValue()).To(Equal(refl.StructValue()))
			})
		It("Should return the same type for the base type as for the type", func() {
			Expect(refl.RawType()).To(Equal(refl.Type()))
		})
		It("Should return the correct pointer type", func() {
			Expect(refl.PointerType()).To(Equal(mType))
		})
		It("Should return the correct pointer value", func() {
			Expect(refl.PointerValue()).To(Equal(reflect.ValueOf(m)))
		})
		It("Should return the correct value for set", func() {
			Expect(refl.ValueForSet().Type()).To(Equal(mType))
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
				newM := refl.NewStruct()
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
		Describe("For Each", func() {
			It("Should provide the reflect itself", func() {
				refl.ForEach(func(rfl *model.Reflect, i int) {
					Expect(i).To(Equal(-1))
					Expect(rfl).To(Equal(refl))
				})
			})
		})
		Describe("PKS", func() {
			It("Should return the correct PK", func() {
				Expect(refl.PKChain()).To(HaveLen(1))
				Expect(refl.PKChain()[0].Interface()).To(Equal(m.ID))
			})
		})
		Describe("Accessing ChainValue", func() {
			It("Should panic", func() {
				Expect(func() {
					refl.ChainValue()
				}).To(Panic())
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
		It("Should pass validation without panicking", func() {
			Expect(refl.Validate).ToNot(Panic())
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
		It("Should return a slice for the base value", func() {
			Expect(refl.RawValue().Interface()).To(Equal(m))
			Expect(refl.RawType()).To(Equal(mBaseType))
		})
		It("Should return the correct pointer type", func() {
			Expect(refl.PointerType()).To(Equal(mType))
		})
		It("Should return the correct pointer value", func() {
			Expect(refl.PointerValue()).To(Equal(reflect.ValueOf(&m)))
		})
		It("Should return the correct value for set", func() {
			Expect(refl.ValueForSet().Type()).To(Equal(mBaseType))
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
				newM := refl.NewStruct()
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
		Describe("For Each", func() {
			It("Should iterate through the chain", func() {
				refl.ForEach(func(rfl *model.Reflect, i int) {
					Expect(rfl).To(Equal(refl.ChainValueByIndex(i)))
				})
			})
		})
		Describe("PKS", func() {
			It("Should get the correct value by PK", func() {
				val, ok := refl.ValueByPK(model.NewPK(m[0].ID))
				Expect(ok).To(BeTrue())
				Expect(val.Type()).To(Equal(mSingleBaseType))
			})
			It("Should return not ok when the value can't be found", func() {
				_, ok := refl.ValueByPK(model.NewPK(uuid.New()))
				Expect(ok).To(BeFalse())
			})
		})
		Describe("Accessing StructValue", func() {
			It("Should panic", func() {
				Expect(func() {
					refl.StructValue()
				}).To(Panic())
			})
		})
	})
	Describe("Errors + edge cases", func() {
		It("Should panic when a non pointer is provided", func() {
			Expect(model.NewReflect(mock.ModelA{ID: 22}).Validate).To(Panic())
		})
		It("Should panic when a non struct is provided", func() {
			i := 11
			Expect(model.NewReflect(&i).Validate).To(Panic())
		})
		Context("nil pointer", func() {
			It("Should not panic when initializing with a nil struct", func() {
				refl := model.NewReflect((*mock.ModelA)(nil))
				Expect(refl.Validate).ToNot(Panic())
			})
			It("Shouldn't panic when initializing with a nil chain", func() {
				var m []*mock.ModelA
				refl := model.NewReflect(&m)
				Expect(refl.Validate).ToNot(Panic())
			})
		})
	})
})
