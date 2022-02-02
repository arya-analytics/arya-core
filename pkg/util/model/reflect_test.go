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
			Expect(model.UnsafeNewReflect(mock.ModelA{}).IsPointer()).To(BeFalse())
		})
	})
	Describe("Pointer Creation", func() {
		It("Should create a new pointer for a non-pointer model", func() {
			Expect(model.UnsafeNewReflect(mock.ModelA{}).ToNewPointer().IsPointer()).To(BeTrue())
		})
		It("Should create the pointer to the correct underlying value", func() {
			var baseVal []*mock.ModelA
			baseRfl := model.UnsafeNewReflect(baseVal)
			Expect(baseRfl.ToNewPointer().RawValue().Interface()).To(Equal(baseVal))
		})
	})
	Context("Single Model", func() {
		var m = &mock.ModelA{
			ID: 22,
		}
		var mBaseType = reflect.TypeOf(mock.ModelA{})
		var mType = reflect.TypeOf(m)
		var rfl = model.NewReflect(m)
		It("Should pass validation without panicking", func() {
			Expect(rfl.Validate).ToNot(Panic())
		})
		It("Should return the correct pointer interface", func() {
			Expect(rfl.Pointer()).To(Equal(m))
		})
		It("Should return the correct type", func() {
			Expect(rfl.Type()).To(Equal(mBaseType))
		})
		It("Should return the correct value", func() {
			Expect(rfl.StructValue().Type()).To(Equal(mBaseType))
		})
		It("Should return false for IsChain", func() {
			Expect(rfl.IsChain()).To(BeFalse())
		})
		It("Should return true for IsStruct", func() {
			Expect(rfl.IsStruct()).To(BeTrue())
		})
		It("Should return the correct struct field by name", func() {
			Expect(rfl.StructValue().FieldByName("ID").Interface()).To(Equal(22))
		})
		It("Should return the correct struct field by role", func() {
			Expect(rfl.StructFieldByRole("pk").Interface()).To(Equal(22))
		})
		It("Should panic if the role doesn't exist", func() {
			Expect(func() {
				rfl.StructFieldByRole("nonexistentrole").Interface()
			}).To(Panic())
		})

		It("Should return the correct struct field by index", func() {
			Expect(rfl.StructValue().Field(0).Interface()).To(Equal(22))
		})
		It("Should return the same item for the raw value as for the value",
			func() {
				Expect(rfl.RawValue()).To(Equal(rfl.StructValue()))
			})
		It("Should return the same type for the raw type as for the type", func() {
			Expect(rfl.RawType()).To(Equal(rfl.Type()))
		})
		It("Should return the correct pointer type", func() {
			Expect(rfl.PointerType()).To(Equal(mType))
		})
		It("Should return the correct pointer value", func() {
			Expect(rfl.PointerValue()).To(Equal(reflect.ValueOf(m)))
		})
		Describe("New Chain", func() {
			It("Should return the correct type", func() {
				newC := rfl.NewChain()
				Expect(newC.RawType()).To(Equal(reflect.TypeOf([]*mock.ModelA{})))
				Expect(newC.Type()).To(Equal(mBaseType))
			})
		})
		Describe("New Model", func() {
			It("Should return the correct type", func() {
				newM := rfl.NewStruct()
				Expect(newM.PointerType()).To(Equal(mType))
				Expect(newM.Type()).To(Equal(mBaseType))
			})
		})
		Describe("New Raw", func() {
			It("Should return a single model", func() {
				newM := rfl.NewRaw()
				Expect(newM.PointerType()).To(Equal(mType))
				Expect(newM.Type()).To(Equal(mBaseType))
			})
		})
		Describe("For Each", func() {
			It("Should provide the reflect itself", func() {
				rfl.ForEach(func(rfl *model.Reflect, i int) {
					Expect(i).To(Equal(-1))
					Expect(rfl).To(Equal(rfl))
				})
			})
		})
		Describe("PKS", func() {
			It("Should return the correct PK", func() {
				Expect(rfl.PKChain()).To(HaveLen(1))
				Expect(rfl.PKChain()[0].Raw()).To(Equal(m.ID))
			})
		})
		Describe("Accessing ChainValue", func() {
			It("Should panic", func() {
				Expect(func() {
					rfl.ChainValue()
				}).To(Panic())
			})
		})
	})
	Context("Multiple Models", func() {
		var (
			m               []*mock.ModelA
			mBaseType       reflect.Type
			mType           reflect.Type
			mSingleBaseType reflect.Type
			mSingleType     reflect.Type
			rfl             *model.Reflect
		)
		BeforeEach(func() {
			m = []*mock.ModelA{
				{
					ID: 22,
				},
				{
					ID: 43,
				},
			}
			mBaseType = reflect.TypeOf(m)
			mType = reflect.TypeOf(&m)
			mSingleBaseType = reflect.TypeOf(mock.ModelA{})
			mSingleType = reflect.TypeOf(&mock.ModelA{})
			rfl = model.NewReflect(&m)
		})
		It("Should pass validation without panicking", func() {
			Expect(rfl.Validate).ToNot(Panic())
		})
		It("Should return the correct pointer interface", func() {
			Expect(rfl.Pointer()).To(Equal(&m))
		})
		It("Should return the correct type", func() {
			Expect(rfl.Type()).To(Equal(mSingleBaseType))
		})
		It("Should return true for IsChain", func() {
			Expect(rfl.IsChain()).To(BeTrue())
		})
		It("Should return false for IsStruct", func() {
			Expect(rfl.IsStruct()).To(BeFalse())
		})
		It("Should return the correct model value by index", func() {
			Expect(rfl.ChainValueByIndex(0).PointerType()).To(Equal(mSingleType))
			Expect(rfl.ChainValueByIndex(0).Type()).To(Equal(mSingleBaseType))
			Expect(rfl.ChainValueByIndexOrNew(0).Pointer()).To(Equal(m[0]))
		})
		It("Should create a new reflect if the index exceeds the chain value", func() {
			Expect(rfl.ChainValueByIndexOrNew(rfl.ChainValue().Len()).Type()).To(
				Equal(mSingleBaseType))
		})
		It("Should return a slice for the raw value", func() {
			Expect(rfl.RawValue().Interface()).To(Equal(m))
			Expect(rfl.RawType()).To(Equal(mBaseType))
		})
		It("Should return the correct pointer type", func() {
			Expect(rfl.PointerType()).To(Equal(mType))
		})
		It("Should return the correct pointer value", func() {
			Expect(rfl.PointerValue()).To(Equal(reflect.ValueOf(&m)))
		})
		Describe("Appending to a chain", func() {
			It("Should append to the chain correctly", func() {
				mTwo := &mock.ModelA{ID: 1}
				rfl.ChainAppend(model.NewReflect(mTwo))
				Expect(rfl.ChainValueByIndex(2).Pointer()).To(Equal(mTwo))
			})
		})
		Describe("New Chain", func() {
			It("Should return the correct type", func() {
				newC := rfl.NewChain()
				Expect(newC.RawType()).To(Equal(mBaseType))
				Expect(newC.Type()).To(Equal(mSingleBaseType))
			})
		})
		Describe("New Model", func() {
			It("Should return the correct type", func() {
				newM := rfl.NewStruct()
				Expect(newM.PointerType()).To(Equal(mSingleType))
				Expect(newM.Type()).To(Equal(mSingleBaseType))
			})
		})
		Describe("New Raw", func() {
			It("Should return a model chain", func() {
				newM := rfl.NewRaw()
				Expect(newM.PointerType()).To(Equal(mType))
				Expect(newM.RawType()).To(Equal(mBaseType))
				Expect(newM.Type()).To(Equal(mSingleBaseType))
			})
		})
		Describe("For Each", func() {
			It("Should iterate through the chain", func() {
				rfl.ForEach(func(rflIter *model.Reflect, i int) {
					Expect(rflIter).To(Equal(rfl.ChainValueByIndex(i)))
				})
			})
		})
		Describe("PKS", func() {
			It("Should get the correct value by PK", func() {
				val, ok := rfl.ValueByPK(model.NewPK(m[0].ID))
				Expect(ok).To(BeTrue())
				Expect(val.Type()).To(Equal(mSingleBaseType))
			})
			It("Should return not ok when the value can't be found", func() {
				_, ok := rfl.ValueByPK(model.NewPK(uuid.New()))
				Expect(ok).To(BeFalse())
			})
		})
		Describe("Accessing StructValue", func() {
			It("Should panic", func() {
				Expect(func() {
					rfl.StructValue()
				}).To(Panic())
			})
		})
	})
	Describe("Errors + edge cases", func() {
		It("Should panic when a non pointer is provided", func() {
			Expect(func() {
				model.NewReflect(mock.ModelA{ID: 22})
			}).To(Panic())
		})
		It("Should panic when a non struct is provided", func() {
			i := 11
			Expect(func() {
				model.NewReflect(&i)
			}).To(Panic())
		})
		Context("nil pointer", func() {
			It("Should panic when initializing with a nil struct", func() {
				Expect(func() {
					model.NewReflect((*mock.ModelA)(nil))
				}).To(Panic())
			})
			It("Shouldn't panic when initializing with a nil chain", func() {
				var m []*mock.ModelA
				Expect(func() {
					model.NewReflect(&m)
				}).ToNot(Panic())
			})
		})
	})
})
